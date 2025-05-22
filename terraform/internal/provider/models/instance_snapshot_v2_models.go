package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
)

var (
	InstanceSnapshotDetailDataType = map[string]attr.Type{
		"id":         types.StringType,
		"name":       types.StringType,
		"status":     types.StringType,
		"size":       types.Int64Type,
		"created_at": types.StringType,
		"labels": types.ListType{
			ElemType: types.StringType,
		},
	}

	InstanceSnapshotListItemType = map[string]attr.Type{
		"instance_id":   types.StringType,
		"instance_name": types.StringType,
		"snapshots": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: InstanceSnapshotDetailDataType,
			},
		},
	}
)

type TFInstanceSnapshot struct {
	Region     types.String `tfsdk:"region"`
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	InstanceID types.String `tfsdk:"instance_id"`
	Labels     types.List   `tfsdk:"labels"`
}

func (i *TFInstanceSnapshot) SetLabels(ctx context.Context, tags []string) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, types.StringType, tags)
	if d.HasError() {
		return d
	}
	i.Labels = l
	return d
}

func (i *TFInstanceSnapshot) GetLabels(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := i.Labels.ElementsAs(ctx, &ret, true)
	return ret, d
}

type TFInstanceSnapshotDSModel struct {
	Region  types.String `tfsdk:"region"`
	Details types.List   `tfsdk:"details"`
}

func (i *TFInstanceSnapshotDSModel) SetDetails(details []TFInstanceSnapshotListItem) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range details {
		vals = append(vals, types.ObjectValueMust(InstanceSnapshotListItemType, map[string]attr.Value{
			"instance_id":   x.InstanceID,
			"instance_name": x.InstanceName,
			"snapshots":     x.Snapshots,
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: InstanceSnapshotListItemType}, vals)
	if d.HasError() {
		return d
	}
	i.Details = l
	return d
}

type TFInstanceSnapshotListItem struct {
	InstanceID   types.String `tfsdk:"instance_id"`
	InstanceName types.String `tfsdk:"instance_name"`
	Snapshots    types.List   `tfsdk:"snapshots"`
}

func (i *TFInstanceSnapshotListItem) SetSnapshots(ctx context.Context, snapshots []api.SnapshotDetailsData) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range snapshots {
		labels, d := types.ListValueFrom(ctx, types.StringType, x.Labels)
		if d.HasError() {
			return d
		}
		vals = append(vals, types.ObjectValueMust(InstanceSnapshotDetailDataType, map[string]attr.Value{
			"id":         types.StringValue(x.ID),
			"name":       types.StringValue(x.Name),
			"status":     types.StringValue(x.Status),
			"size":       types.Int64Value(x.Size),
			"created_at": types.StringValue(x.GetFormattedTime()),
			"labels":     labels,
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: InstanceSnapshotDetailDataType}, vals)
	if d.HasError() {
		return d
	}
	i.Snapshots = l
	return d
}
