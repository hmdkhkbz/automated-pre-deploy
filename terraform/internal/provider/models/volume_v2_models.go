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
	VolumeV2ObjectType = map[string]attr.Type{
		"region":        types.StringType,
		"id":            types.StringType,
		"name":          types.StringType,
		"size":          types.Int64Type,
		"type":          types.StringType,
		"instance_name": types.StringType,
		"status":        types.StringType,
		"tags": types.ListType{
			ElemType: types.StringType,
		},
	}
)

type TFVolumeV2Model struct {
	Region       types.String `tfsdk:"region"`
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Size         types.Int64  `tfsdk:"size"`
	Type         types.String `tfsdk:"type"`
	InstanceName types.String `tfsdk:"instance_name"`
	Status       types.String `tfsdk:"status"`
	Tags         types.Set    `tfsdk:"tags"`
	SnapshotID   types.String `tfsdk:"snapshot_id"`
}

func (v *TFVolumeV2Model) GetTags(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := v.Tags.ElementsAs(ctx, &ret, true)
	return ret, d
}

func (v *TFVolumeV2Model) SetTags(ctx context.Context, tags []string) diag.Diagnostics {
	l, d := types.SetValueFrom(ctx, types.StringType, tags)
	if d.HasError() {
		return d
	}
	v.Tags = l
	return d
}

type TFVolumeV2DateSourceModel struct {
	Region  types.String `tfsdk:"region"`
	Volumes types.List   `tfsdk:"volumes"`
}

func (m *TFVolumeV2DateSourceModel) SetVolumes(ctx context.Context, region string, volumes []api.ListData) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range volumes {
		tags, d := types.ListValueFrom(ctx, types.StringType, x.Labels)
		if d.HasError() {
			return d
		}
		vals = append(vals, types.ObjectValueMust(VolumeV2ObjectType, map[string]attr.Value{
			"region":        types.StringValue(region),
			"id":            types.StringValue(x.ID),
			"name":          types.StringValue(x.Name),
			"size":          types.Int64Value(int64(x.Size)),
			"type":          types.StringValue(x.VolumeType),
			"instance_name": types.StringValue(x.InstanceName),
			"status":        types.StringValue(x.Status),
			"tags":          tags,
		}))
	}
	l, d := types.ListValue(basetypes.ObjectType{AttrTypes: VolumeV2ObjectType}, vals)
	if d.HasError() {
		return d
	}
	m.Volumes = l
	return d
}
