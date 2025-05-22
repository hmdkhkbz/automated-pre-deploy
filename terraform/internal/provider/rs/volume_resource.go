package rs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
)

type VolumeResource struct {
	client *api.Client
}

func (v *VolumeResource) SetAPIClient(c *api.Client) {
	v.client = c
}

func (v *VolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (v *VolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, v)
}

func (v *VolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"size": schema.Int64Attribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("hdd-g1"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ssd-g1", "hdd-g1"),
				},
			},
			"attachments": schema.ListAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.ObjectType{
					AttrTypes: models.AttachmentObject,
				},
			},
		},
	}
}

func (v *VolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TFVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := api.ServerVolume{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Size:        int32(data.Size.ValueInt64()),
		Type:        data.Type.ValueString(),
	}
	createResp, err := v.client.Volume.CreateVolume(ctx, data.Region.ValueString(), &createReq)
	if err != nil {
		resp.Diagnostics.AddError("error creating volume", err.Error())
		return
	}

	data.ID = types.StringValue(createResp.ID)
	tflog.Info(ctx, "create_volume_response", map[string]interface{}{"volume": createResp})
	var vals []models.TFVolumeAttachment
	for _, at := range createResp.Attachments {

		vals = append(vals, models.TFVolumeAttachment{
			ID:           types.StringValue(at.ID),
			Device:       types.StringValue(at.Device),
			VolumeID:     types.StringValue(at.VolumeID),
			AttachmentID: types.StringValue(at.AttachmentID),
			ServerID:     types.StringValue(at.ServerID),
			ServerName:   types.StringValue(at.ServerName),
			HostName:     types.StringValue(at.HostName),
		})

	}

	resp.Diagnostics.Append(data.SetAttachments(ctx, vals)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (v *VolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TFVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	respData, err := v.client.Volume.GetVolume(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok {
			if respErr.Code == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.AddError("error reading volume", err.Error())
		return
	}

	data.Name = types.StringValue(respData.Name)
	data.Description = types.StringValue(respData.Description)
	data.Size = types.Int64Value(int64(respData.Size))

	if len(respData.Attachments) > 0 {
		var vals []models.TFVolumeAttachment
		for _, at := range respData.Attachments {
			vals = append(vals, models.TFVolumeAttachment{
				ID:           types.StringValue(at.ID),
				Device:       types.StringValue(at.Device),
				VolumeID:     types.StringValue(at.VolumeID),
				AttachmentID: types.StringValue(at.AttachmentID),
				ServerID:     types.StringValue(at.ServerID),
				ServerName:   types.StringValue(at.ServerName),
				HostName:     types.StringValue(at.HostName),
			})
		}
		resp.Diagnostics.Append(data.SetAttachments(ctx, vals)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (v *VolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TFVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	err := v.client.WaitForCondition(ctx, time.Minute * 2, func() (bool, error) {
		vol, err := v.client.Volume.GetVolume(ctx, data.Region.ValueString(), data.ID.ValueString())
		if err != nil {
			return false, nil
		}
		return vol.Status == "available", nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error deleting volume", "volume not available")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	err = v.client.Volume.DeleteVolume(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting volume", err.Error())
	}

}

func (v *VolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.TFVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.TFVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.Size.Equal(stateData.Size) {
		if planData.Size.ValueInt64() < stateData.Size.ValueInt64() {
			resp.Diagnostics.AddError("volumes cannot be shrinked to smaller size", "volume size can only be extended")
			return
		}

		_, err := v.client.VolumeV2.Resize(ctx, planData.Region.ValueString(), planData.ID.ValueString(), &api.ResizeRequest{
			Size: int(planData.Size.ValueInt64()),
		})
		if err != nil {
			resp.Diagnostics.AddError("error resizing volume", err.Error())
			return
		}
	}

	if !planData.Name.Equal(stateData.Name) || !planData.Description.Equal(stateData.Description) {
		err := v.client.Volume.UpdateVolume(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.ServerVolume{
			Name:        planData.Name.ValueString(),
			Description: planData.Description.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("error updating volume", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func NewVolumeResource() resource.Resource {
	return &VolumeResource{}
}
