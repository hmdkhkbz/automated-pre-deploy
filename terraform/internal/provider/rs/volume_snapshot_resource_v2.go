package rs

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
)

type VolumeSnapshotV2Resource struct {
	client *api.Client
}

func (s *VolumeSnapshotV2Resource) SetAPIClient(c *api.Client) {
	s.client = c
}

func (s *VolumeSnapshotV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume_snapshot_v2"
}

func (s *VolumeSnapshotV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, s)
}

func (s *VolumeSnapshotV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("built_by_terraform")})),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"volume_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (s *VolumeSnapshotV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TFVolumeSnapshotV2
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := s.client.BackupV2.CreateVolumeSnapshot(ctx, data.Region.ValueString(), &api.CreateVolumeSnapshot{
		Name:     data.Name.ValueString(),
		VolumeID: data.VolumeID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error creating volume snapshot", err.Error())
		return
	}
	data.ID = types.StringValue(apiResp.SnapshotID)
	tags, d := data.GetLabels(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err = s.client.BackupV2.EditSnapshotLabels(ctx, data.Region.ValueString(), apiResp.SnapshotID, &api.EditSnapshotLabels{
		Labels: tags,
	})
	if err != nil {
		resp.Diagnostics.AddError("error adding tags to volume snapshot", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s *VolumeSnapshotV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TFVolumeSnapshotV2
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := s.client.BackupV2.GetVolumeSnapshot(ctx, data.Region.ValueString(), data.VolumeID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching volume snapshot", err.Error())
		return
	}
	data.Name = types.StringValue(apiResp.Name)
	resp.Diagnostics.Append(data.SetLabels(ctx, apiResp.Labels)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (s *VolumeSnapshotV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TFVolumeSnapshotV2
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := s.client.BackupV2.DeleteSnapshot(ctx, data.Region.ValueString(), &api.DeleteSnapshot{
		SnapshotIDs: []string{data.ID.ValueString()},
	})
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting snapshot", err.Error())
		return
	}
}

func (s *VolumeSnapshotV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.TFVolumeSnapshotV2
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var stateData models.TFVolumeSnapshotV2
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !planData.VolumeID.Equal(stateData.VolumeID) {
		resp.Diagnostics.AddError("invalid operation", "volume id can not be changed")
		return
	}
	if !planData.Name.Equal(stateData.Name) {
		_, err := s.client.BackupV2.EditSnapshotName(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.EditSnapshotName{
			Name: planData.Name.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("error updating snapshot name", err.Error())
			return
		}
	}

	if !planData.Labels.Equal(stateData.Labels) {
		tags, d := planData.GetLabels(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		_, err := s.client.BackupV2.EditSnapshotLabels(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.EditSnapshotLabels{
			Labels: tags,
		})
		if err != nil {
			resp.Diagnostics.AddError("error updating volume snapshot labels", err.Error())
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func NewVolumeSnapshotV2Resource() resource.Resource {
	return &VolumeSnapshotV2Resource{}
}
