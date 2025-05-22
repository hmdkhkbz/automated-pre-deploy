package ds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
)

type VolumeV2Datasource struct {
	client *api.Client
}

func (v *VolumeV2Datasource) SetAPIClient(c *api.Client) {
	v.client = c
}

func (v *VolumeV2Datasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volumes_v2"
}

func (v *VolumeV2Datasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, v)
}

/*
ID           string   `json:"id"`
	Name         string   `json:"name"`
	Size         int      `json:"size"`
	InstanceName string   `json:"instance_name"`
	VolumeType   string   `json:"type"`
	Status       string   `json:"status"`
	Labels       []string `json:"labels"`
*/

func (v *VolumeV2Datasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"volumes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"size": schema.Int64Attribute{
							Computed: true,
						},
						"instance_name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (v *VolumeV2Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFVolumeV2DateSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := v.client.VolumeV2.List(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching volumes", err.Error())
		return
	}

	resp.Diagnostics.Append(data.SetVolumes(ctx, data.Region.ValueString(), apiResp.Data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewVolumeV2Datasource() datasource.DataSource {
	return &VolumeV2Datasource{}
}
