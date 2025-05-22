package models

import (
	"context"
	"terraform-provider-hashicups-pf/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	serverGroupObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"policies": types.ListType{
				ElemType: types.StringType,
			},
			"members": types.ListType{
				ElemType: types.StringType,
			},
		},
	}
)

type TFServerGroup struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Policies types.List   `tfsdk:"policies"`
	Members  types.List   `tfsdk:"members"`
}

type TFServerGroupDatasourceModel struct {
	Region       types.String        `tfsdk:"region"`
	ServerGroups basetypes.ListValue `tfsdk:"server_groups"`
}

func (s *TFServerGroupDatasourceModel) SetServerGroups(ctx context.Context, serverGroups []api.ServerGroupDetail) diag.Diagnostics {
	var k []TFServerGroup
	for _, x := range serverGroups {
		policies, d := types.ListValueFrom(ctx, types.StringType, x.Policies)
		if d.HasError() {
			return d
		}

		members, d := types.ListValueFrom(ctx, types.StringType, x.Members)
		if d.HasError() {
			return d
		}

		k = append(k, TFServerGroup{
			ID:       types.StringValue(x.ID),
			Name:     types.StringValue(x.Name),
			Policies: policies,
			Members:  members,
		})
	}
	l, d := types.ListValueFrom(ctx, serverGroupObjType, k)
	if d.HasError() {
		return d
	}
	s.ServerGroups = l
	return d
}
