package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-hashicups-pf/internal/api"
)

var (
	subnetServerObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":      types.StringType,
			"ip":      types.StringType,
			"name":    types.StringType,
			"port_id": types.StringType,
		},
	}

	ipRangeObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"start": types.StringType,
			"end":   types.StringType,
		},
	}

	subnetObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"subnet_id":   types.StringType,
			"network_id":  types.StringType,
			"name":        types.StringType,
			"description": types.StringType,
			"dhcp_range":  ipRangeObjType,
			"dns_servers": basetypes.SetType{
				ElemType: types.StringType,
			},
			"enable_dhcp":    types.BoolType,
			"enable_gateway": types.BoolType,
			"gateway_ip":     types.StringType,
			"cidr":           types.StringType,
			"servers": basetypes.SetType{
				ElemType: subnetServerObjType,
			},
			"shared": types.BoolType,
			"mtu":    types.Int64Type,
		},
	}
)

type TFIPRange struct {
	Start types.String `tfsdk:"start"`
	End   types.String `tfsdk:"end"`
}

type TFSubnetModel struct {
	ID            types.String `tfsdk:"id"`
	Region        types.String `tfsdk:"region"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	DHCPRange     types.Object `tfsdk:"dhcp_range"`
	DNSServer     types.List   `tfsdk:"dns_servers"`
	EnableDHCP    types.Bool   `tfsdk:"enable_dhcp"`
	EnableGateway types.Bool   `tfsdk:"enable_gateway"`
	GatewayIP     types.String `tfsdk:"gateway_ip"`
	CIDR          types.String `tfsdk:"cidr"`
	NetworkID     types.String `tfsdk:"network_id"`
}

func (s *TFSubnetModel) GetDHCPRange(ctx context.Context) (TFIPRange, diag.Diagnostics) {
	var ret TFIPRange
	d := s.DHCPRange.As(ctx, &ret, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	return ret, d
}

func (s *TFSubnetModel) SetDHCPRange(ctx context.Context, r TFIPRange) diag.Diagnostics {
	obj, d := types.ObjectValueFrom(ctx, ipRangeObjType.AttrTypes, &r)
	if d.HasError() {
		return d
	}
	s.DHCPRange = obj
	return d
}

func (s *TFSubnetModel) GetDNSServers(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := s.DNSServer.ElementsAs(ctx, &ret, true)
	return ret, d
}

func (s *TFSubnetModel) SetDNSServers(ctx context.Context, ns []string) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, types.StringType, &ns)
	if d.HasError() {
		return d
	}
	s.DNSServer = l
	return d
}

type TFSubnetServerDatasourceModel struct {
	ID     types.String `tfsdk:"id"`
	IP     types.String `tfsdk:"ip"`
	Name   types.String `tfsdk:"name"`
	PortID types.String `tfsdk:"port_id"`
}

type TFIPRangeDatasourceModel struct {
	Start types.String `tfsdk:"start"`
	End   types.String `tfsdk:"end"`
}

type TFSubnetDatasourceModel struct {
	Shared        types.Bool            `tfsdk:"shared"`
	MTU           types.Int64           `tfsdk:"mtu"`
	SubnetID      types.String          `tfsdk:"subnet_id"`
	NetworkID     types.String          `tfsdk:"network_id"`
	Name          types.String          `tfsdk:"name"`
	Description   types.String          `tfsdk:"description"`
	DHCPRange     basetypes.ObjectValue `tfsdk:"dhcp_range"`
	DNSServer     basetypes.SetValue    `tfsdk:"dns_servers"`
	EnableDHCP    types.Bool            `tfsdk:"enable_dhcp"`
	EnableGateway types.Bool            `tfsdk:"enable_gateway"`
	GatewayIP     types.String          `tfsdk:"gateway_ip"`
	CIDR          types.String          `tfsdk:"cidr"`
	Servers       basetypes.SetValue    `tfsdk:"servers"`
}

func TFSubnetDatasourceModelFromAPINetwork(ctx context.Context, n *api.Network) (*TFSubnetDatasourceModel, diag.Diagnostics) {
	ret := TFSubnetDatasourceModel{
		NetworkID:   types.StringValue(n.ID),
		Name:        types.StringValue(n.Name),
		Description: types.StringValue(n.Description),
		DHCPRange:   types.ObjectNull(ipRangeObjType.AttrTypes),
		DNSServer:   types.SetNull(types.StringType),
		Servers:     types.SetNull(subnetServerObjType),
		Shared:      types.BoolValue(n.Shared),
		MTU:         types.Int64Value(int64(n.MTU)),
	}
	if len(n.Subnets) > 0 {
		s := n.Subnets[0]
		ret.SubnetID = types.StringValue(s.ID)
		if len(s.AllocationPools) > 0 {
			dRange := TFIPRangeDatasourceModel{
				Start: types.StringValue(s.AllocationPools[0].Start),
				End:   types.StringValue(s.AllocationPools[0].End),
			}
			obj, d := types.ObjectValueFrom(ctx, ipRangeObjType.AttrTypes, dRange)
			if d.HasError() {
				return nil, d
			}
			ret.DHCPRange = obj
		}
		if len(s.DNSNameservers) > 0 {
			var ns []attr.Value
			for _, x := range s.DNSNameservers {
				ns = append(ns, types.StringValue(x))
			}
			ret.DNSServer = types.SetValueMust(types.StringType, ns)
		}

		ret.EnableDHCP = types.BoolValue(s.EnableDHCP)
		ret.EnableGateway = types.BoolValue(s.GatewayIP != nil)
		if s.GatewayIP != nil {
			ret.GatewayIP = types.StringValue(*s.GatewayIP)
		}

		ret.CIDR = types.StringValue(s.CIDR)

		if len(s.Servers) > 0 {
			var tfServers []TFSubnetServerDatasourceModel
			for _, i := range s.Servers {
				tfs := TFSubnetServerDatasourceModel{
					ID:   types.StringValue(i.ID),
					Name: types.StringValue(i.Name),
				}
				for _, ip := range i.IPs {
					if ip.SubnetID == s.ID {
						tfs.IP = types.StringValue(ip.IP)
						tfs.PortID = types.StringValue(ip.PortID)
						break
					}
				}
				tfServers = append(tfServers, tfs)
			}
			servers, d := types.SetValueFrom(ctx, subnetServerObjType, tfServers)
			if d.HasError() {
				return nil, d
			}
			ret.Servers = servers
		}
	}
	return &ret, diag.Diagnostics{}
}

type TFNetworkFilter struct {
	Public    types.Bool   `tfsdk:"public"`
	Name      types.String `tfsdk:"name"`
	NetworkID types.String `tfsdk:"network_id"`
}

func (f *TFNetworkFilter) Filter(n *TFSubnetDatasourceModel) bool {
	var res = true
	if !f.Public.IsNull() {
		res = n.Shared.Equal(f.Public)
	}
	if !f.NetworkID.IsNull() {
		res = res && n.NetworkID.Equal(f.NetworkID)
	}
	if !f.Name.IsNull() {
		res = res && n.Name.Equal(f.Name)
	}
	return res
}

type TFNetworkDatasourceModel struct {
	Filters  basetypes.ObjectValue `tfsdk:"filters"`
	Region   types.String          `tfsdk:"region"`
	Networks basetypes.SetValue    `tfsdk:"networks"`
}

func (n *TFNetworkDatasourceModel) SetNetworks(ctx context.Context, nets []TFSubnetDatasourceModel) diag.Diagnostics {
	var filteredNets = nets
	if !n.Filters.IsUnknown() && !n.Filters.IsNull() {
		f, d := n.GetFilters(ctx)
		if d.HasError() {
			return d
		}
		filteredNets = nil
		for _, x := range nets {
			if f.Filter(&x) {
				filteredNets = append(filteredNets, x)
			}
		}
	}

	s, d := types.SetValueFrom(ctx, subnetObjType, filteredNets)
	if d.HasError() {
		return d
	}

	n.Networks = s
	return d
}

func (n *TFNetworkDatasourceModel) GetFilters(ctx context.Context) (*TFNetworkFilter, diag.Diagnostics) {
	var ret TFNetworkFilter
	d := n.Filters.As(ctx, &ret, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if d.HasError() {
		return nil, d
	}
	return &ret, d
}
