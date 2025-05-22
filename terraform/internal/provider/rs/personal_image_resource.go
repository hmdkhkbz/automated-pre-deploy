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
)

type PersonalImageResource struct {
	client *api.Client
}

func (p *PersonalImageResource) SetAPIClient(c *api.Client) {
	p.client = c
}

func (p *PersonalImageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_image"
}

func (p *PersonalImageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, p)
}

func (p *PersonalImageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"snapshot_id": schema.StringAttribute{
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
		},
	}
}

func (p *PersonalImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.SnapshotPersonalImageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	img, err := p.client.Img.GetPrivateImageByID(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if misc.RemoveResourceIfNotFound(ctx, resp, err) {
			return
		}
		resp.Diagnostics.AddError("error fetching private image", err.Error())
		return
	}
	data.Name = types.StringValue(img.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PersonalImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.SnapshotPersonalImageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := p.client.SnapshotClient.CreatePersonalImage(ctx, data.Region.ValueString(), data.Name.ValueString(), data.SnapshotID.ValueString())
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError("error creating personal image from snapshot", err.Error())
			return
		}
	}
	data.ID = types.StringValue(apiResp.ID)
	data.Size = types.Int64Value(int64(apiResp.Size))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (p *PersonalImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.SnapshotPersonalImageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := p.client.Volume.DeleteVolume(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting personal image", err.Error())
		return
	}
}

func (p *PersonalImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func NewPersonalImageResource() resource.Resource {
	return &PersonalImageResource{}
}
