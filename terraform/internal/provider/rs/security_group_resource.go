package rs

import (
	"context"
	"net"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/utl"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecurityGroupResource struct {
	client *api.Client
}

func (s *SecurityGroupResource) SetAPIClient(c *api.Client) {
	s.client = c
}

func (s *SecurityGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_group"
}

func (s *SecurityGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, s)
}

func (s *SecurityGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"readonly": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"rules": schema.SetNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"description": schema.StringAttribute{
							Optional: true,
						},
						"direction": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ingress", "egress"),
							},
						},
						"ether_type": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"group_id": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"ip": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								utl.CIDRValidator(),
							},
						},
						"port_from": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								utl.PortValidator(),
							},
						},
						"port_to": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								utl.PortValidator(),
							},
						},
						"protocol": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (s *SecurityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData models.TFSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiSg, err := s.client.Firewall.CreateSecurityGroup(ctx, planData.Region.ValueString(), planData.Name.ValueString(), planData.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error creating security group", err.Error())
		return
	}

	//this is because the api creates two default rules which will conflict with plan data later
	//whe the resource will be refreshed or read
	for _, x := range apiSg.Rules {
		err = s.client.Firewall.DeleteRule(ctx, planData.Region.ValueString(), x.ID)
		if err != nil {
			resp.Diagnostics.AddError("error deleting security group", err.Error())
			return
		}
	}

	newState := &models.TFSecurityGroupModel{
		Region:      planData.Region,
		ID:          types.StringValue(apiSg.ID),
		Name:        planData.Name,
		Description: planData.Description,
		Default:     types.BoolValue(apiSg.Default),
		ReadOnly:    types.BoolValue(apiSg.ReadOnly),	
	}
	ruleIdMap := make(map[string]bool)
	planRules, d := planData.GetRules(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	for idx, x := range planRules {
		// validate rule IP address to be in CIDR format
		if !x.IP.IsNull() {
			_, _, err := net.ParseCIDR(x.IP.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("IP address must be in CIDR format(Ex. 10.0.0.1/32)", x.IP.ValueString())
			}
		}

		ruleApiReq := x.ToCreateRuleAPIReq(ctx)

		err := s.client.Firewall.CreateRule(ctx, planData.Region.ValueString(), apiSg.ID, ruleApiReq)
		if err != nil {
			resp.Diagnostics.AddError("error creating rule", err.Error())
			return
		}

		groupResp, err := s.client.Firewall.GetSecurityGroupByID(ctx, planData.Region.ValueString(), apiSg.ID)
		if err != nil {
			resp.Diagnostics.AddError("error fetching security group", err.Error())
			return
		}

		for _, rule := range groupResp.Rules {
			if ruleIdMap[rule.ID] {
				continue
			}
			resp.Diagnostics.Append(planRules[idx].PopulateFromAPIResponse(ctx, rule)...)
			ruleIdMap[rule.ID] = true
			break
		}
	}
	resp.Diagnostics.Append(newState.SetRules(ctx, planRules)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (s *SecurityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateData models.TFSecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := s.client.Firewall.GetSecurityGroupByID(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
	if err != nil {

		if misc.RemoveResourceIfNotFound(ctx, resp, err) {
			return
		}
		resp.Diagnostics.AddError("error fetching security group", err.Error())
		return
	}
	utl.AssignStringIfChanged(&stateData.Name, apiResp.Name)
	utl.AssignStringIfChanged(&stateData.Description, apiResp.Description)
	if !stateData.Default.Equal(types.BoolValue(apiResp.Default)) {
		stateData.Default = types.BoolValue(apiResp.Default)
	}
	var newTFRules []models.TFSecGroupRuleModel
	for _, r := range apiResp.Rules {
		var tfRule models.TFSecGroupRuleModel
		resp.Diagnostics.Append(tfRule.PopulateFromAPIResponse(ctx, r)...)
		if resp.Diagnostics.HasError() {
			return
		}
		newTFRules = append(newTFRules, tfRule)
	}
	resp.Diagnostics.Append(stateData.SetRules(ctx, newTFRules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (s *SecurityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateData models.TFSecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallInstances, err := s.client.FirewallV2.GetFirewallConnectedInstances(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error getting security group info", err.Error())
		return
	}

	// detach all instances from firewall before deleting
	instanceIDs := make([]string, 0) 
	for _, i := range firewallInstances {
		instanceIDs = append(instanceIDs, i.InstanceID)
	}

	err = s.client.FirewallV2.DetachInstancesFromFirewall(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), instanceIDs)
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Message != utl.ErrFirewalNotAttached {
			resp.Diagnostics.AddError("error detaching instances from security group", err.Error())
			return
		}
	}

	// finally delete the security group
	err = s.client.Firewall.DeleteSecurityGroup(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting security group", err.Error())
		return
	}
}

func (s *SecurityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateData models.TFSecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planData models.TFSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.Rules.Equal(stateData.Rules) {
		var stateRuleIDs []types.String
		stateRules, d := stateData.GetRules(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, x := range stateRules {
			stateRuleIDs = append(stateRuleIDs, x.ID)
		}

		var planRuleIDs []types.String
		planRules, d := planData.GetRules(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, x := range planRules {
			planRuleIDs = append(planRuleIDs, x.ID)
		}
		planRuleSet := utl.ListToSet(planRuleIDs)
		if planRuleSet[""] {
			delete(planRuleSet, "")
		}

		stateRuleSet := utl.ListToSet(stateRuleIDs)

		for i := 0; i < len(planRules); i++ {
			if planRules[i].ID.IsNull() || planRules[i].ID.IsUnknown() {
				err := s.client.Firewall.CreateRule(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), planRules[i].ToCreateRuleAPIReq(ctx))
				if err != nil {
					resp.Diagnostics.AddError("error creating rule", err.Error())
					return
				}
				resp.Diagnostics.Append(s.populateCreatedRule(ctx, &stateData, &planRules[i], stateRuleSet)...)
				if resp.Diagnostics.HasError() {
					return
				}
				planRuleSet[planRules[i].ID.ValueString()] = true
			}
		}

		for k, _ := range stateRuleSet {
			if planRuleSet[k] {
				continue
			}
			err := s.client.Firewall.DeleteRule(ctx, stateData.Region.ValueString(), k)
			if err != nil {
				resp.Diagnostics.AddError("error deleting rule", err.Error())
				return
			}
		}
		resp.Diagnostics.Append(planData.SetRules(ctx, planRules)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (s *SecurityGroupResource) populateCreatedRule(ctx context.Context, state *models.TFSecurityGroupModel, rule *models.TFSecGroupRuleModel, existingRuleSet map[string]bool) diag.Diagnostics {
	var d diag.Diagnostics
	sg, err := s.client.Firewall.GetSecurityGroupByID(ctx, state.Region.ValueString(), state.ID.ValueString())
	if err != nil {
		d.AddError("error fetching security groups", err.Error())
		return d
	}
	for _, r := range sg.Rules {
		if existingRuleSet[r.ID] {
			continue
		}
		d := rule.PopulateFromAPIResponse(ctx, r)
		if d.HasError() {
			return d
		}
		existingRuleSet[r.ID] = true
		break
	}
	return d
}

func NewSecurityGroupResource() resource.Resource {
	return &SecurityGroupResource{}
}
