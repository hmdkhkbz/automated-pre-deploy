package models

import (
	"context"
	"fmt"
	"net"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/utl"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	ruleType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"description": types.StringType,
			"direction":   types.StringType,
			"ether_type":  types.StringType,
			"group_id":    types.StringType,
			"ip":          types.StringType,
			"port_from":   types.StringType,
			"port_to":     types.StringType,
			"protocol":    types.StringType,
		},
	}

	groupType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":      types.StringType,
			"id":          types.StringType,
			"name":        types.StringType,
			"description": types.StringType,
			"readonly":    types.BoolType,
			"default":     types.BoolType,
			"rules": basetypes.SetType{
				ElemType: ruleType,
			},
		},
	}
)

type TFSecurityGroup struct {
	ID          types.String     `tfsdk:"id"`
	Name        types.String     `tfsdk:"name"`
	Description types.String     `tfsdk:"description"`
	Default     types.Bool       `tfsdk:"default"`
	ReadOnly    types.Bool       `tfsdk:"readonly"`
	IPAddresses []types.String   `tfsdk:"ip_addresses"`
	Rules       []TFSecGroupRule `tfsdk:"rules"`
}

type TFSecurityGroupModel struct {
	Region      types.String `tfsdk:"region"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
	ReadOnly    types.Bool   `tfsdk:"readonly"`
	Rules       types.Set    `tfsdk:"rules"`
}

func (s *TFSecurityGroupModel) GetRules(ctx context.Context) ([]TFSecGroupRuleModel, diag.Diagnostics) {
	var d diag.Diagnostics
	var ret []TFSecGroupRuleModel
	d = s.Rules.ElementsAs(ctx, &ret, true)

	if d.HasError() {
		return nil, d
	}
	return ret, d
}

func (s *TFSecurityGroupModel) SetRules(ctx context.Context, rules []TFSecGroupRuleModel) diag.Diagnostics {
	r, d := types.SetValueFrom(ctx, ruleType, &rules)
	if d.HasError() {
		return d
	}
	s.Rules = r
	return d
}

type TFSecGroupRuleModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Direction   types.String `tfsdk:"direction"`
	EtherType   types.String `tfsdk:"ether_type"`
	GroupID     types.String `tfsdk:"group_id"`
	IP          types.String `tfsdk:"ip"`
	PortFrom    types.String `tfsdk:"port_from"`
	PortTo      types.String `tfsdk:"port_to"`
	Protocol    types.String `tfsdk:"protocol"`
}

func (r *TFSecGroupRuleModel) ToCreateRuleAPIReq(ctx context.Context) *api.RuleRequest {
	ret := &api.RuleRequest{
		Description: r.Description.ValueString(),
		Direction:   r.Direction.ValueString(),
		PortStart:   r.PortFrom.ValueString(),
		PortEnd:     r.PortTo.ValueString(),
		Protocol:    r.Protocol.ValueString(),
	}
	if !r.IP.IsNull() {
		ret.IP = []string{r.IP.ValueString()}
	}
	return ret
}

func (r *TFSecGroupRuleModel) PopulateFromAPIResponse(ctx context.Context, apiResp *api.Rule) diag.Diagnostics {
	var d diag.Diagnostics
	utl.AssignStringIfChanged(&r.ID, apiResp.ID)
	utl.AssignStringIfChanged(&r.GroupID, apiResp.GroupID)
	utl.AssignStringIfChanged(&r.Description, apiResp.Description)
	utl.AssignStringIfChanged(&r.Direction, apiResp.Direction)
	utl.AssignStringIfChanged(&r.EtherType, apiResp.EtherType)
	utl.AssignStringIfChanged(&r.Protocol, apiResp.Protocol)
	if apiResp.PortStart != 0 {
		utl.AssignStringIfChanged(&r.PortFrom, fmt.Sprintf("%d", apiResp.PortStart))
	}
	if apiResp.PortEnd != 0 {
		utl.AssignStringIfChanged(&r.PortTo, fmt.Sprintf("%d", apiResp.PortEnd))
	}

	if apiResp.IP != "" {
		_, _, err := net.ParseCIDR(apiResp.IP)
		if err != nil {
			d.AddError("error parsing cidr", err.Error())
			return d
		}
		utl.AssignStringIfChanged(&r.IP, apiResp.IP)
	}
	return d

}

type TFSecGroupRule struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Direction   types.String `tfsdk:"direction"`
	EtherType   types.String `tfsdk:"ether_type"`
	GroupID     types.String `tfsdk:"group_id"`
	IP          types.String `tfsdk:"ip"`
	PortStart   types.Int64  `tfsdk:"port_start"`
	PortEnd     types.Int64  `tfsdk:"port_end"`
	Protocol    types.String `tfsdk:"protocol"`
}

type TFSecurityGroupDataSourceModel struct {
	Region types.String        `tfsdk:"region"`
	Groups basetypes.ListValue `tfsdk:"groups"`
}

func (s *TFSecurityGroupDataSourceModel) SetGroups(ctx context.Context, groups []TFSecurityGroupModel) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, groupType, groups)
	if d.HasError() {
		return d
	}
	s.Groups = l
	return d
}
