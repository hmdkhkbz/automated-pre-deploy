package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	volumeSnapshotObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":      types.StringType,
			"id":          types.StringType,
			"volume_id":   types.StringType,
			"name":        types.StringType,
			"description": types.StringType,
			"volume_name": types.StringType,
			"status":      types.StringType,
			"size":        types.Int64Type,
		},
	}

	serverSnapshotObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"region":      types.StringType,
			"id":          types.StringType,
			"size":        types.Int64Type,
			"name":        types.StringType,
			"description": types.StringType,
			"volume_id":   types.StringType,
			"image_id":    types.StringType,
			"status":      types.StringType,
			"server_id":   types.StringType,
			"server_name": types.StringType,
		},
	}
)

type VolumeSnapshot struct {
	Region      types.String `tfsdk:"region"`
	ID          types.String `tfsdk:"id"`
	VolumeID    types.String `tfsdk:"volume_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VolumeName  types.String `tfsdk:"volume_name"`
	Status      types.String `tfsdk:"status"`
	Size        types.Int64  `tfsdk:"size"`
}

type SnapshotDatasourceModel struct {
	Region    types.String        `tfsdk:"region"`
	Snapshots basetypes.ListValue `tfsdk:"snapshots"`
}

func (s *SnapshotDatasourceModel) SetVolumeSnapshots(ctx context.Context, data []VolumeSnapshot) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, volumeSnapshotObjType, data)
	if d.HasError() {
		return d
	}
	s.Snapshots = l
	return d
}

type ServerSnapshot struct {
	Region      types.String `tfsdk:"region"`
	ID          types.String `tfsdk:"id"`
	Size        types.Int64  `tfsdk:"size"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VolumeID    types.String `tfsdk:"volume_id"`
	ImageID     types.String `tfsdk:"image_id"`
	Status      types.String `tfsdk:"status"`
	ServerID    types.String `tfsdk:"server_id"`
	ServerName  types.String `tfsdk:"server_name"`
}

type TFServerSnapshotDatasourceModel struct {
	Region    types.String        `tfsdk:"region"`
	Snapshots basetypes.ListValue `tfsdk:"snapshots"`
}

func (s *TFServerSnapshotDatasourceModel) SetSnapshots(ctx context.Context, data []ServerSnapshot) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, serverSnapshotObjType, data)
	if d.HasError() {
		return d
	}
	s.Snapshots = l
	return d
}

type SnapshotPersonalImageResourceModel struct {
	Region     types.String `tfsdk:"region"`
	Name       types.String `tfsdk:"name"`
	SnapshotID types.String `tfsdk:"snapshot_id"`
	ID         types.String `tfsdk:"id"`
	Size       types.Int64  `tfsdk:"size"`
}
