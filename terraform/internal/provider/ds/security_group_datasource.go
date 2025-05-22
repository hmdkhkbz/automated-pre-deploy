package ds

import (
	"context"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecurityGroupDatasource struct {
	client *api.Client
}

func (s *SecurityGroupDatasource) SetAPIClient(client *api.Client) {
	s.client = client
}

func (s *SecurityGroupDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_groups"
}

func (s *SecurityGroupDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, s)
}

func (s *SecurityGroupDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"default": schema.BoolAttribute{
							Computed: true,
						},
						"readonly": schema.BoolAttribute{
							Computed: true,
						},
						"rules": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"description": schema.StringAttribute{
										Computed: true,
									},
									"direction": schema.StringAttribute{
										Computed: true,
									},
									"ether_type": schema.StringAttribute{
										Computed: true,
									},
									"group_id": schema.StringAttribute{
										Computed: true,
									},
									"ip": schema.StringAttribute{
										Computed: true,
									},
									"port_from": schema.StringAttribute{
										Computed: true,
									},
									"port_to": schema.StringAttribute{
										Computed: true,
									},
									"protocol": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *SecurityGroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFSecurityGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := s.client.Firewall.GetAllSecurityGroups(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching security groups", err.Error())
		return
	}

	var tfGroups []models.TFSecurityGroupModel
	for _, x := range apiResp {
		tfG := models.TFSecurityGroupModel{
			Region:      data.Region,
			ID:          types.StringValue(x.ID),
			Name:        types.StringValue(x.Name),
			Description: types.StringValue(x.Description),
			Default:     types.BoolValue(x.Default),
			ReadOnly:    types.BoolValue(x.ReadOnly),
		}
		var tfRules []models.TFSecGroupRuleModel
		for _, r := range x.Rules {
			tfR := models.TFSecGroupRuleModel{}
			resp.Diagnostics.Append(tfR.PopulateFromAPIResponse(ctx, r)...)
			if resp.Diagnostics.HasError() {
				return
			}
			tfRules = append(tfRules, tfR)
		}
		resp.Diagnostics.Append(tfG.SetRules(ctx, tfRules)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tfGroups = append(tfGroups, tfG)
	}

	resp.Diagnostics.Append(data.SetGroups(ctx, tfGroups)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewSecurityGroupDatasource() datasource.DataSource {
	return &SecurityGroupDatasource{}
}
