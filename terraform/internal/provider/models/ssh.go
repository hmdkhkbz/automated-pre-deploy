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
	sshKeyObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":       types.StringType,
			"public_key": types.StringType,
		},
	}
)

type TFSSHKey struct {
	Name      types.String `tfsdk:"name"`
	PublicKey types.String `tfsdk:"public_key"`
}

type TFSSHKeyDatasourceModel struct {
	Region types.String        `tfsdk:"region"`
	Keys   basetypes.ListValue `tfsdk:"keys"`
}

func (s *TFSSHKeyDatasourceModel) SetKeys(ctx context.Context, keys []*api.SSHKey) diag.Diagnostics {
	var k []TFSSHKey
	for _, x := range keys {
		k = append(k, TFSSHKey{
			Name:      types.StringValue(x.Name),
			PublicKey: types.StringValue(x.PublicKey),
		})
	}
	l, d := types.ListValueFrom(ctx, sshKeyObjType, k)
	if d.HasError() {
		return d
	}
	s.Keys = l
	return d
}
