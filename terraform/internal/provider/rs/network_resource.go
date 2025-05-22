package rs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
)

type NetworkResource struct {
	client *api.Client
}

func (n *NetworkResource) SetAPIClient(c *api.Client) {
	n.client = c
}

func (n *NetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (n *NetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, n)
}

func (n *NetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"dhcp_range": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"start": schema.StringAttribute{
						Required: true,
					},
					"end": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"dns_servers": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"enable_dhcp": schema.BoolAttribute{
				Required: true,
			},
			"enable_gateway": schema.BoolAttribute{
				Required: true,
			},
			"gateway_ip": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{
					//misc.NewRequiredWhenAttributeHasBoolValue(true, path.Root("enable_gateway")),
				},
			},
			"cidr": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (n *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData models.TFSubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planData.EnableDHCP.ValueBool() {
		if planData.DNSServer.IsNull() || planData.DNSServer.IsUnknown() {
			resp.Diagnostics.AddError("invalid operation", "dns_servers must be set when dhcp is enabled")
			return
		}

		if planData.DHCPRange.IsNull() || planData.DHCPRange.IsUnknown() {
			resp.Diagnostics.AddError("invalid operation", "dhcp_range must be set when dhcp is enabled")
			return
		}
	} else {
		if !planData.DNSServer.IsNull() {
			resp.Diagnostics.AddError("invalid operation", "dns_servers must not be set when dhcp is disabled")
			return
		}

		if !planData.DHCPRange.IsNull() {
			resp.Diagnostics.AddError("invalid operation", "dhcp_range must not be set when dhcp is disabled")
			return
		}
	}

	apiReq := api.Subnet{
		Name:          planData.Name.ValueString(),
		EnableDHCP:    planData.EnableDHCP.ValueBool(),
		EnableGateway: planData.EnableGateway.ValueBool(),
		SubnetGateway: planData.GatewayIP.ValueString(),
		CIDR:          planData.CIDR.ValueString(),
		Description:   planData.Description.ValueString(),
	}

	if !planData.DHCPRange.IsUnknown() && !planData.DHCPRange.IsNull() && !planData.DNSServer.IsNull() && !planData.DNSServer.IsUnknown() && planData.EnableDHCP.ValueBool() {
		dhcpRange, d := planData.GetDHCPRange(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.DHCPRange = strings.Join([]string{dhcpRange.Start.ValueString(), dhcpRange.End.ValueString()}, ",")
		slcDNS, diag := planData.GetDNSServers(ctx)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		dnsServers := ""
		for idx, x := range slcDNS {
			dnsServers += x
			if idx < len(slcDNS)-1 {
				dnsServers += "\n"
			}
		}
		apiReq.DNSServers = dnsServers
	}

	apiResp, err := n.client.Subnet.CreatePrivateNetwork(ctx, planData.Region.ValueString(), &apiReq)
	if err != nil {
		resp.Diagnostics.AddError("error creating network", err.Error())
		return
	}

	planData.NetworkID = types.StringValue(apiResp.NetworkID)
	planData.ID = types.StringValue(apiResp.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (n *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "PRIVATE_NET_READ_STARTED")
	var state models.TFSubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := n.client.Subnet.GetPrivateNetwork(ctx, state.Region.ValueString(), state.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok {
			if respErr.Code == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.AddError("error fetching network", err.Error())
		return
	}
	//state.CIDR.Replace(apiResp.CIDR)
	if state.Name.ValueString() != apiResp.Name {
		state.Name = types.StringValue(apiResp.Name)
	}

	if !state.Description.IsNull() && !state.Description.IsUnknown() && state.Description.ValueString() != apiResp.Description {
		state.Description = types.StringValue(apiResp.Description)
	}

	if apiResp.GatewayIP != nil && state.GatewayIP.ValueString() != *apiResp.GatewayIP {
		state.GatewayIP = types.StringValue(*apiResp.GatewayIP)
	}

	if len(apiResp.AllocationPools) == 1 && apiResp.EnableDHCP {
		resp.Diagnostics.Append(state.SetDHCPRange(ctx, models.TFIPRange{
			Start: types.StringValue(apiResp.AllocationPools[0].Start),
			End:   types.StringValue(apiResp.AllocationPools[0].End),
		})...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if apiResp.EnableDHCP {
		resp.Diagnostics.Append(state.SetDNSServers(ctx, apiResp.DNSNameservers)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !state.EnableDHCP.Equal(types.BoolValue(apiResp.EnableDHCP)) {
		state.EnableDHCP = types.BoolValue(apiResp.EnableDHCP)
	}
	enableGateway := apiResp.GatewayIP != nil && *apiResp.GatewayIP != ""
	if !state.EnableGateway.Equal(types.BoolValue(enableGateway)) {
		state.EnableGateway = types.BoolValue(enableGateway)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (n *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.TFSubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := n.client.Subnet.DeletePrivateNetwork(ctx, state.Region.ValueString(), state.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting network", err.Error())
		return
	}
}

func (n *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.TFSubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planData.EnableDHCP.ValueBool() {
		if planData.DNSServer.IsNull() || planData.DNSServer.IsUnknown() {
			resp.Diagnostics.AddError("invalid operation", "dns_servers must be set when dhcp is enabled")
			return
		}

		if planData.DHCPRange.IsNull() || planData.DHCPRange.IsUnknown() {
			resp.Diagnostics.AddError("invalid operation", "dhcp_range must be set when dhcp is enabled")
			return
		}
	} else {
		if !planData.DNSServer.IsNull() {
			resp.Diagnostics.AddError("invalid operation", "dns_servers must not be set when dhcp is disabled")
			return
		}

		if !planData.DHCPRange.IsNull() {
			resp.Diagnostics.AddError("invalid operation", "dhcp_range must not be set when dhcp is disabled")
			return
		}
	}

	var stateData models.TFSubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	n.handleSimpleUpdates(ctx, &stateData, &planData, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := &api.Subnet{
		Name:          planData.Name.ValueString(),
		EnableDHCP:    planData.EnableDHCP.ValueBool(),
		EnableGateway: planData.EnableGateway.ValueBool(),
		NetworkID:     stateData.NetworkID.ValueString(),
		SubnetGateway: planData.GatewayIP.ValueString(),
		SubnetID:      planData.ID.ValueString(),
		CIDR:          stateData.CIDR.ValueString(),
		Description:   planData.Description.ValueString(),
	}

	if !planData.DHCPRange.IsNull() && !planData.EnableDHCP.IsUnknown() && planData.EnableDHCP.ValueBool() {
		dhcpRange, d := planData.GetDHCPRange(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.DHCPRange = strings.Join([]string{dhcpRange.Start.ValueString(), dhcpRange.End.ValueString()}, ",")
	}

	if !planData.DNSServer.IsNull() && !planData.DNSServer.IsUnknown() && planData.EnableDHCP.ValueBool() {
		dnsList, diags := planData.GetDNSServers(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		dnsServers := ""
		for idx, x := range dnsList {
			dnsServers += x
			if idx < len(dnsList)-1 {
				dnsServers += "\n"
			}
		}
		apiReq.DNSServers = dnsServers
	}

	err := n.client.Subnet.UpdatePrivateNetwork(ctx, stateData.Region.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("error updating network", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (n *NetworkResource) handleSimpleUpdates(ctx context.Context, stateData, planData *models.TFSubnetModel, resp *resource.UpdateResponse) {
	if !planData.CIDR.Equal(stateData.CIDR) {
		tflog.Info(ctx, "CIDR_CHANGE", map[string]interface{}{"STATE": stateData.CIDR.ValueString(), "PLAN": planData.CIDR.ValueString()})
		resp.Diagnostics.AddError("invalid operation", "cidr can only be set at creation time")
		//planData.CIDR = stateData.CIDR
	}
	return
}

func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}
