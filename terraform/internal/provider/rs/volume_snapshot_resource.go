package rs

import (
	"context"
	"fmt"
	"time"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/utl"
)

type VolumeSnapshotResource struct {
	client *api.Client
}

func NewVolumeSnapshotResource() resource.Resource {
	return &VolumeSnapshotResource{}
}

func (v *VolumeSnapshotResource) SetAPIClient(c *api.Client) {
	v.client = c
}

func (v *VolumeSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_volume_snapshot"
}

func (v *VolumeSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, v)
}

func (v *VolumeSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"volume_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"volume_name": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"size": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (v *VolumeSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData models.VolumeSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := v.client.WaitForCondition(ctx, time.Minute * 2, func() (bool, error) {
		volResp, err := v.client.Volume.GetVolume(ctx, planData.Region.ValueString(), planData.VolumeID.ValueString())
		if err != nil {
			return false, err
		}
		cond := volResp.Status == "in-use" || volResp.Status == "available"
		return cond, nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error creating volume snapshot", fmt.Sprintf("volume %s is not ready for creating snapshot", planData.VolumeID.ValueString()))
		return
	}

	apiResp, err := v.client.SnapshotClient.CreateVolumeSnapshot(ctx, planData.Region.ValueString(), planData.VolumeID.ValueString(), &api.SnapshotRequest{
		Description: planData.Description.ValueString(),
		Name:        planData.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error creating snapshot", err.Error())
		return
	}
	planData.ID = types.StringValue(apiResp.ID)
	planData.VolumeName = types.StringValue(apiResp.VolumeName)
	planData.Status = types.StringValue(apiResp.Status)
	planData.Size = types.Int64Value(int64(apiResp.Size))

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (v *VolumeSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.VolumeSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := v.client.SnapshotClient.GetSnapshotByID(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if misc.RemoveResourceIfNotFound(ctx, resp, err) {
			return
		}
		resp.Diagnostics.AddError("error fetching volume snapshot", err.Error())
		return
	}
	utl.AssignStringIfChanged(&data.Status, apiResp.Status)
	utl.AssignStringIfChanged(&data.Description, apiResp.Description)
	utl.AssignStringIfChanged(&data.Name, apiResp.Name)
	utl.AssignStringIfChanged(&data.VolumeName, apiResp.VolumeName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (v *VolumeSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.VolumeSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := v.client.SnapshotClient.DeleteSnapshot(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting snapshot", err.Error())
		return
	}
}

func (v *VolumeSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.VolumeSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.VolumeSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.VolumeID.Equal(stateData.VolumeID) {
		resp.Diagnostics.AddError("invalid operation", "volume id can only be set at creation time")
		return
	}

	apiResp, err := v.client.SnapshotClient.UpdateVolumeSnapshot(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.UpdateVolumeSnapshot{
		Name:        planData.Name.ValueString(),
		Description: planData.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error updating volume snapshot", err.Error())
		return
	}
	planData.Status = types.StringValue(apiResp.Status)
	planData.VolumeName = types.StringValue(apiResp.VolumeName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
