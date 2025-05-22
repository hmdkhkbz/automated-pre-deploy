package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	FloatingIPObjectType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                     types.StringType,
			"status":                 types.StringType,
			"description":            types.StringType,
			"fixed_ip":               types.StringType,
			"floating_ip":            types.StringType,
			"port_id":                types.StringType,
			"attached_instance_id":   types.StringType,
			"attached_instance_name": types.StringType,
		},
	}
)

type TFFloatingIPModel struct {
	Region      types.String `tfsdk:"region"`
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
	Address     types.String `tfsdk:"address"`
}

type TFFloatingIPDataSourceItem struct {
	ID                   types.String `tfsdk:"id"`
	Status               types.String `tfsdk:"status"`
	Description          types.String `tfsdk:"description"`
	FixedIP              types.String `tfsdk:"fixed_ip"`
	FloatingIP           types.String `tfsdk:"floating_ip"`
	PortID               types.String `tfsdk:"port_id"`
	AttachedInstanceID   types.String `tfsdk:"attached_instance_id"`
	AttachedInstanceName types.String `tfsdk:"attached_instance_name"`
}

type TFFloatingIPDataSourceModel struct {
	Region      types.String        `tfsdk:"region"`
	FloatingIPs basetypes.ListValue `tfsdk:"floating_ips"`
}

func (f *TFFloatingIPDataSourceModel) SetFloatingIPs(ctx context.Context, ips []TFFloatingIPDataSourceItem) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, FloatingIPObjectType, ips)
	if d.HasError() {
		return d
	}
	f.FloatingIPs = l
	return d
}
