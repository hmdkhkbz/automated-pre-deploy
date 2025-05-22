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

type InstanceSnapshotResource struct {
	client *api.Client
}

func (i *InstanceSnapshotResource) SetAPIClient(c *api.Client) {
	i.client = c
}

func (i *InstanceSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_snapshot_v2"
}

func (i *InstanceSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, i)
}

func (i *InstanceSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"instance_id": schema.StringAttribute{
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
		},
	}
}

func (i *InstanceSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TFInstanceSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := i.client.BackupV2.CreateInstanceSnapshot(ctx, data.Region.ValueString(), &api.CreateInstanceSnapshotRequest{
		Name:       data.Name.ValueString(),
		InstanceID: data.InstanceID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error creating instance snapshot", err.Error())
		return
	}
	data.ID = types.StringValue(apiResp.SnapshotID)
	labels, d := data.GetLabels(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err = i.client.BackupV2.EditSnapshotLabels(ctx, data.Region.ValueString(), apiResp.SnapshotID, &api.EditSnapshotLabels{
		Labels: labels,
	})
	if err != nil {
		resp.Diagnostics.AddError("error editing instance snapshot labels", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (i *InstanceSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TFInstanceSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := i.client.BackupV2.GetSnapshotDetails(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok {
			if respErr.Code == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.AddError("error reading instance snapshot", err.Error())
		return
	}
	data.Name = types.StringValue(apiResp.Data.Name)
	resp.Diagnostics.Append(data.SetLabels(ctx, apiResp.Data.Labels)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *InstanceSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TFInstanceSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := i.client.BackupV2.DeleteInstanceSnapshot(ctx, data.Region.ValueString(), &api.DeleteInstanceSnapshotsRequest{
		InstanceIDs: []string{data.ID.ValueString()},
	})
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting instance snapshot", err.Error())
	}
}

func (i *InstanceSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.TFInstanceSnapshot
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.TFInstanceSnapshot
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.Name.Equal(stateData.Name) {
		_, err := i.client.BackupV2.EditSnapshotName(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.EditSnapshotName{
			Name: planData.Name.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("error editing instance snapshot name", err.Error())
			return
		}
	}

	if !planData.Labels.Equal(stateData.Labels) {
		labels, d := planData.GetLabels(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		_, err := i.client.BackupV2.EditSnapshotLabels(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), &api.EditSnapshotLabels{
			Labels: labels,
		})
		if err != nil {
			resp.Diagnostics.AddError("error editing instance snapshot labels", err.Error())
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func NewInstanceSnapshotResource() resource.Resource {
	return &InstanceSnapshotResource{}
}
