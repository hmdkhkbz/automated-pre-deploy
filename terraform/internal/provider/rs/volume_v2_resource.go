package rs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
)

type VolumeV2Resource struct {
	c *api.Client
}

func (v *VolumeV2Resource) SetAPIClient(c *api.Client) {
	v.c = c
}

func (v *VolumeV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + "volume_v2"
}

func (v *VolumeV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, v)
}

func (v *VolumeV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"snapshot_id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("built_by_terraform")})),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (v *VolumeV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TFVolumeV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	respData, err := v.c.VolumeV2.Inquire(ctx, data.Region.ValueString(), data.ID.ValueString())
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

	tags, err := v.c.VolumeV2.GetVolumeByID(ctx, data.Region.ValueString(), data.ID.ValueString())
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
	if tags.Labels != nil {
		resp.Diagnostics.Append(data.SetTags(ctx, tags.Labels)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	data.Name = types.StringValue(respData.Name)
	data.Size = types.Int64Value(int64(respData.Size))
	data.Status = types.StringValue(respData.Status)
	data.InstanceName = types.StringValue(respData.InstanceName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (v *VolumeV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TFVolumeV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, d := data.GetTags(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.SnapshotID.IsNull() && !data.SnapshotID.IsUnknown() {
		if !data.Size.IsUnknown() && !data.Size.IsNull() {
			resp.Diagnostics.AddError("invalid operation", "when creating volume from snapshot size must not be provided")
			return
		}

		if !data.Type.IsNull() && !data.Type.IsUnknown() {
			resp.Diagnostics.AddError("invalid operation", "when creating volume from snapshot volume type must not be provided")
			return
		}
		cResp, err := v.c.BackupV2.CreateVolumeFromSnapshot(ctx, data.Region.ValueString(), data.SnapshotID.ValueString(), &api.CreateVolumeFromSnapshot{
			Name: data.Name.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("error creating volume from snapshot", err.Error())
			return
		}
		data.ID = types.StringValue(cResp.VolumeID)
		data.Size = types.Int64Value(int64(cResp.VolumeSize))

		_, err = v.c.VolumeV2.EditLabels(ctx, data.Region.ValueString(), cResp.VolumeID, &api.EditLabelsRequest{
			Labels: tags,
		})

		if err != nil {
			resp.Diagnostics.AddError("error adding labels", err.Error())
			return
		}

		inqData, err := v.c.VolumeV2.Inquire(ctx, data.Region.ValueString(), data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error fetching volume details", err.Error())
			return
		}
		data.Status = types.StringValue(inqData.Status)
		data.InstanceName = types.StringValue(inqData.InstanceName)

		vDetails, err := v.c.VolumeV2.GetVolumeByID(ctx, data.Region.ValueString(), data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error fetching volume details", err.Error())
			return
		}
		data.Type = types.StringValue(vDetails.VolumeType)

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	apiReq := api.VolumeV2CreateRequest{
		Name: data.Name.ValueString(),
		Size: int(data.Size.ValueInt64()),
		Type: data.Type.ValueString(),
	}
	apiResp, err := v.c.VolumeV2.Create(ctx, data.Region.ValueString(), &apiReq)
	if err != nil {
		resp.Diagnostics.AddError("error creating volume", err.Error())
		return
	}

	_, err = v.c.VolumeV2.EditLabels(ctx, data.Region.ValueString(), apiResp.VolumeID, &api.EditLabelsRequest{
		Labels: tags,
	})

	if err != nil {
		resp.Diagnostics.AddError("error adding labels", err.Error())
		return
	}

	data.ID = types.StringValue(apiResp.VolumeID)
	details, err := v.c.VolumeV2.Inquire(ctx, data.Region.ValueString(), apiResp.VolumeID)
	if err != nil {
		resp.Diagnostics.AddError("error reading volume details", err.Error())
		return
	}
	data.Status = types.StringValue(details.Status)
	data.InstanceName = types.StringValue(details.InstanceName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (v *VolumeV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TFVolumeV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := v.c.WaitForCondition(ctx, time.Minute * 2, func() (bool, error) {
		vol, err := v.c.VolumeV2.Inquire(ctx, data.Region.ValueString(), data.ID.ValueString())
		if err != nil {
			return false, nil
		}
		return vol.Status == "available", nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error deleting volume", "volume is not available")
		return
	}

	_, err = v.c.VolumeV2.Delete(ctx, data.Region.ValueString(), &api.VolumeV2DeleteRequest{
		VolumeIDs: []string{data.ID.ValueString()},
	})

	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting volume", err.Error())
	}

}

func (v *VolumeV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.TFVolumeV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.TFVolumeV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.Type.Equal(stateData.Type) {
		resp.Diagnostics.AddError("invalid operation", "can not change volume type after creation")
		return
	}

	if !planData.SnapshotID.Equal(stateData.SnapshotID) {
		resp.Diagnostics.AddError("invalid operation", "can not change snapshot id after volume creation")
		return
	}

	if !planData.Size.Equal(stateData.Size) {
		_, err := v.c.VolumeV2.Resize(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.ResizeRequest{
			Size: int(planData.Size.ValueInt64()),
		})
		if err != nil {
			resp.Diagnostics.AddError("error resizing volume", err.Error())
			return
		}
	}

	if !planData.Name.Equal(stateData.Name) {
		_, err := v.c.VolumeV2.EditName(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.EditNameRequest{
			Name: planData.Name.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("error editing volume name", err.Error())
			return
		}
	}

	if !planData.Tags.Equal(stateData.Tags) {
		tags, d := planData.GetTags(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		_, err := v.c.VolumeV2.EditLabels(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.EditLabelsRequest{
			Labels: tags,
		})
		if err != nil {
			resp.Diagnostics.AddError("error updating labels", err.Error())
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func NewVolumeV2Resource() resource.Resource {
	return &VolumeV2Resource{}
}
