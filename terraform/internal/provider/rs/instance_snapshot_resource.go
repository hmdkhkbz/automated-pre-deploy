package rs

import (
	"context"
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

type ServerSnapshotResource struct {
	client *api.Client
}

func (s *ServerSnapshotResource) SetAPIClient(c *api.Client) {
	s.client = c
}

func (s *ServerSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_snapshot"
}

func (s *ServerSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, s)
}

func (s *ServerSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"size": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"volume_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"server_id": schema.StringAttribute{
				Required: true,
			},
			"server_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (s *ServerSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateData models.ServerSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := s.client.SnapshotClient.GetSnapshotByID(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
	if err != nil {
		if misc.RemoveResourceIfNotFound(ctx, resp, err) {
			return
		}
		resp.Diagnostics.AddError("error fetching server snapshot", err.Error())
		return
	}
	utl.AssignStringIfChanged(&stateData.Name, apiResp.Name)
	utl.AssignStringIfChanged(&stateData.Description, apiResp.Description)
	utl.AssignStringIfChanged(&stateData.Status, apiResp.Status)
	utl.AssignStringIfChanged(&stateData.ServerName, apiResp.ServerName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)

}

func (s *ServerSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData models.ServerSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := s.client.SnapshotClient.CreateServerSnapshot(ctx, planData.Region.ValueString(), planData.ServerID.ValueString(), &api.SnapshotRequest{
		Description: planData.Description.ValueString(),
		Name:        planData.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error creating server snapshot", err.Error())
		return
	}
	planData.ID = types.StringValue(apiResp.ID)
	planData.Status = types.StringValue(apiResp.Status)
	planData.ServerName = types.StringValue(apiResp.ServerName)
	planData.VolumeID = types.StringValue(apiResp.VolumeID)
	planData.ImageID = types.StringValue(apiResp.ImageID)
	planData.Size = types.Int64Value(int64(apiResp.Size))

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (s *ServerSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.ServerSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := s.client.SnapshotClient.DeleteSnapshot(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting snapshot", err.Error())
		return
	}
}

func (s *ServerSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.ServerSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.ServerSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !planData.ServerID.Equal(stateData.ServerID) {
		resp.Diagnostics.AddError("invalid operation", "server id can only be set at creation time")
		return
	}

	apiResp, err := s.client.SnapshotClient.UpdateVolumeSnapshot(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.UpdateVolumeSnapshot{
		Name:        planData.Name.ValueString(),
		Description: planData.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error updating server snapshot", err.Error())
		return
	}

	planData.Status = types.StringValue(apiResp.Status)
	planData.ServerName = types.StringValue(apiResp.ServerName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func NewServerSnapshotResource() resource.Resource {
	return &ServerSnapshotResource{}
}
