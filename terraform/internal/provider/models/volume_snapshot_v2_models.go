package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
)

var (
	VolumeSnapshotV2DetailType = map[string]attr.Type{
		"id":         types.StringType,
		"name":       types.StringType,
		"status":     types.StringType,
		"size":       types.Int64Type,
		"created_at": types.StringType,
		"tags": types.ListType{
			ElemType: types.StringType,
		},
	}

	VolumeSnapshotV2ListItemType = map[string]attr.Type{
		"volume_id":   types.StringType,
		"volume_name": types.StringType,
		"snapshots": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: VolumeSnapshotV2DetailType,
			},
		},
	}
)

type TFVolumeSnapshotV2 struct {
	Region   types.String `tfsdk:"region"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Labels   types.List   `tfsdk:"labels"`
	VolumeID types.String `tfsdk:"volume_id"`
}

func (v *TFVolumeSnapshotV2) SetLabels(ctx context.Context, tags []string) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, types.StringType, tags)
	if d.HasError() {
		return d
	}
	v.Labels = l
	return d
}

func (v *TFVolumeSnapshotV2) GetLabels(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := v.Labels.ElementsAs(ctx, &ret, true)
	return ret, d
}

type TFVolumeSnapshotV2DSModel struct {
	Region  types.String `tfsdk:"region"`
	Details types.List   `tfsdk:"details"`
}

func (m *TFVolumeSnapshotV2DSModel) SetDetails(details []TFVolumeSnapshotV2DSListItem) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range details {
		vals = append(vals, types.ObjectValueMust(VolumeSnapshotV2ListItemType, map[string]attr.Value{
			"volume_id":   x.VolumeID,
			"volume_name": x.VolumeName,
			"snapshots":   x.Snapshots,
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: VolumeSnapshotV2ListItemType}, vals)
	if d.HasError() {
		return d
	}
	m.Details = l
	return d
}

type TFVolumeSnapshotV2DSListItem struct {
	VolumeID   types.String `tfsdk:"volume_id"`
	VolumeName types.String `tfsdk:"volume_name"`
	Snapshots  types.List   `tfsdk:"snapshots"`
}

func (i *TFVolumeSnapshotV2DSListItem) SetSnapshots(ctx context.Context, snapshots []api.SnapshotDetailsData) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range snapshots {
		labels, d := types.ListValueFrom(ctx, types.StringType, x.Labels)
		if d.HasError() {
			return d
		}
		vals = append(vals, types.ObjectValueMust(VolumeSnapshotV2DetailType, map[string]attr.Value{
			"id":         types.StringValue(x.ID),
			"name":       types.StringValue(x.Name),
			"status":     types.StringValue(x.Status),
			"size":       types.Int64Value(x.Size),
			"created_at": types.StringValue(x.GetFormattedTime()),
			"labels":     labels,
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: VolumeSnapshotV2DetailType}, vals)
	if d.HasError() {
		return d
	}
	i.Snapshots = l
	return d
}

type TFVolumeSnapshotV2DSDetailItem struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Status    types.String `tfsdk:"status"`
	Size      types.Int64  `tfsdk:"size"`
	CreatedAt types.String `tfsdk:"created_at"`
	Labels    types.List   `tfsdk:"labels"`
}

func (i *TFVolumeSnapshotV2DSDetailItem) SetLabels(ctx context.Context, labels []string) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, types.StringType, labels)
	if d.HasError() {
		return d
	}
	i.Labels = l
	return d
}
