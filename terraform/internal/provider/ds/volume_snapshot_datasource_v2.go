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

type VolumeSnapshotV2Datasource struct {
	client *api.Client
}

func (v *VolumeSnapshotV2Datasource) SetAPIClient(c *api.Client) {
	v.client = c
}

func (v *VolumeSnapshotV2Datasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume_snapshots_v2"
}

func (v *VolumeSnapshotV2Datasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, v)
}

func (v *VolumeSnapshotV2Datasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"details": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"volume_id": schema.StringAttribute{
							Computed: true,
						},
						"volume_name": schema.StringAttribute{
							Computed: true,
						},
						"snapshots": schema.ListNestedAttribute{
							Optional: true,
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"status": schema.StringAttribute{
										Computed: true,
									},
									"size": schema.Int64Attribute{
										Computed: true,
									},
									"created_at": schema.StringAttribute{
										Computed: true,
									},
									"labels": schema.ListAttribute{
										Optional:    true,
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (v *VolumeSnapshotV2Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFVolumeSnapshotV2DSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := v.client.BackupV2.ListVolumeSnapshots(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching volume snapshots", err.Error())
		return
	}
	var tfList []models.TFVolumeSnapshotV2DSListItem
	for _, x := range apiResp.Data {
		item := models.TFVolumeSnapshotV2DSListItem{
			VolumeID:   types.StringValue(x.VolumeID),
			VolumeName: types.StringValue(x.VolumeName),
		}
		apiSnapshots, err := v.client.BackupV2.VolumeSnapshotDetails(ctx, data.Region.ValueString(), x.VolumeID)
		if err != nil {
			resp.Diagnostics.AddError("error fetching snapshot details", err.Error())
			return
		}
		resp.Diagnostics.Append(item.SetSnapshots(ctx, apiSnapshots.Snapshots)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tfList = append(tfList, item)
	}
	resp.Diagnostics.Append(data.SetDetails(tfList)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewVolumeSnapshotV2Datasource() datasource.DataSource {
	return &VolumeSnapshotV2Datasource{}
}
