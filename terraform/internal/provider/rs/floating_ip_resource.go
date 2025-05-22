package rs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/utl"
)

type FloatingIPResource struct {
	client *api.Client
}

func (f *FloatingIPResource) SetAPIClient(client *api.Client) {
	f.client = client
}

func (f *FloatingIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_floating_ip"
}

func (f *FloatingIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, f)
}

func (f *FloatingIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"description": schema.StringAttribute{
				Required: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (f *FloatingIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planData models.TFFloatingIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := f.client.FIPClient.CreateFloatingIP(ctx, planData.Region.ValueString(), planData.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error deleting floating ip", err.Error())
		return
	}
	planData.ID = types.StringValue(apiResp.ID)
	planData.Address = types.StringValue(apiResp.FloatingIPAddress)
	planData.Status = types.StringValue(apiResp.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (f *FloatingIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TFFloatingIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := f.client.FIPClient.GetFloatingIP(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if misc.RemoveResourceIfNotFound(ctx, resp, err) {
			return
		}
		resp.Diagnostics.AddError("error fetching floating ip", err.Error())
		return
	}
	utl.AssignStringIfChanged(&data.Description, apiResp.Description)
	utl.AssignStringIfChanged(&data.Status, apiResp.Status)
	utl.AssignStringIfChanged(&data.Address, apiResp.FloatingIPAddress)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (f *FloatingIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TFFloatingIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := f.client.WaitForCondition(ctx, 2 * time.Minute, func() (bool, error) {
		apiResp, err := f.client.FIPClient.GetFloatingIP(ctx, data.Region.ValueString(), data.ID.ValueString())
		if err != nil {
			return false, err
		}
		return apiResp.Status == "DOWN", nil
	})

	if err != nil {
		resp.Diagnostics.AddError("error deleting floating ip", "floating ip is not ready to delete")
		return
	}

	err = f.client.FIPClient.DeleteFloatingIP(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting floating ip", err.Error())
		return
	}
}

func (f *FloatingIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData models.TFFloatingIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.TFFloatingIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planData.Status = stateData.Status
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)

}

func NewFloatingIPResource() resource.Resource {
	return &FloatingIPResource{}
}
