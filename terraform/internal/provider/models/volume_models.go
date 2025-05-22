package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	AttachmentObject = map[string]attr.Type{
		"id":            types.StringType,
		"device":        types.StringType,
		"volume_id":     types.StringType,
		"attachment_id": types.StringType,
		"server_id":     types.StringType,
		"server_name":   types.StringType,
		"host_name":     types.StringType,
	}

	VolumeObjectType = map[string]attr.Type{
		"region":      types.StringType,
		"id":          types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"type":        types.StringType,
		"size":        types.Int64Type,
		"attachments": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: AttachmentObject,
			},
		},
	}
)

type TFVolumeModel struct {
	Region      types.String        `tfsdk:"region"`
	ID          types.String        `tfsdk:"id"`
	Name        types.String        `tfsdk:"name"`
	Description types.String        `tfsdk:"description"`
	Size        types.Int64         `tfsdk:"size"`
	Type        types.String        `tfsdk:"type"`
	Attachments basetypes.ListValue `tfsdk:"attachments"`
}

func (v *TFVolumeModel) GetAttachments(ctx context.Context) ([]TFVolumeAttachment, diag.Diagnostics) {
	var d diag.Diagnostics
	var ret []TFVolumeAttachment
	d = v.Attachments.ElementsAs(ctx, &ret, true)
	return ret, d
}

func (v *TFVolumeModel) SetAttachments(ctx context.Context, attrs []TFVolumeAttachment) diag.Diagnostics {
	var vals []attr.Value
	for _, at := range attrs {

		vals = append(vals, types.ObjectValueMust(AttachmentObject, map[string]attr.Value{
			"id":            at.ID,
			"device":        at.Device,
			"volume_id":     at.VolumeID,
			"attachment_id": at.AttachmentID,
			"server_id":     at.ServerID,
			"server_name":   at.ServerName,
			"host_name":     at.HostName,
		}))

	}
	var d diag.Diagnostics
	var l basetypes.ListValue
	l, d = types.ListValue(basetypes.ObjectType{AttrTypes: AttachmentObject}, vals)
	if d.HasError() {
		return d
	}
	v.Attachments = l
	return d
}

type TFVolumeAttachment struct {
	ID           types.String `tfsdk:"id"`
	Device       types.String `tfsdk:"device"`
	VolumeID     types.String `tfsdk:"volume_id"`
	AttachmentID types.String `tfsdk:"attachment_id"`
	ServerID     types.String `tfsdk:"server_id"`
	ServerName   types.String `tfsdk:"server_name"`
	HostName     types.String `tfsdk:"host_name"`
}

type TFVolumeDataSourceModel struct {
	Region  types.String        `tfsdk:"region"`
	Volumes basetypes.ListValue `tfsdk:"volumes"`
}

func (v *TFVolumeDataSourceModel) SetVolumes(ctx context.Context, vols []TFVolumeModel) diag.Diagnostics {
	/*var vals []attr.Value
	for _, x := range vols {
		var attachmentAttrs []attr.Value
		tfAttrs, d := x.GetAttachments(ctx)
		if d.HasError() {
			return d
		}
		for _, at := range tfAttrs {

			attachmentAttrs = append(attachmentAttrs, types.ObjectValueMust(AttachmentObject, map[string]attr.Value{
				"id":            at.ID,
				"device":        at.Device,
				"volume_id":     at.VolumeID,
				"attachment_id": at.AttachmentID,
				"server_id":     at.ServerID,
				"server_name":   at.ServerName,
				"host_name":     at.HostName,
			}))

		}
		vals = append(vals, types.ObjectValueMust(VolumeObjectType, map[string]attr.Value{
			"region":      x.Region,
			"id":          x.ID,
			"attachments": types.ListValueMust(types.ObjectType{AttrTypes: AttachmentObject}, attachmentAttrs),
			"name":        x.Name,
			"description": x.Description,
			"size":        x.Size,
		}))
	}*/
	l, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: VolumeObjectType}, vols)
	if d.HasError() {
		return d
	}
	v.Volumes = l
	return d
}
