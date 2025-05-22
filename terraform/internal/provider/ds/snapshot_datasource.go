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

type SnapshotDatasource struct {
	client *api.Client
}

func (s *SnapshotDatasource) SetAPIClient(c *api.Client) {
	s.client = c
}

func (s *SnapshotDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshots"
}

func (s *SnapshotDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, s)
}

func (s *SnapshotDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"snapshots": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"volume_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"volume_name": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"size": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (s *SnapshotDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.SnapshotDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := s.client.SnapshotClient.ListVolumeSnapshots(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching volume snapshots", err.Error())
		return
	}

	var tfSnaps []models.VolumeSnapshot
	for _, x := range apiResp {
		snp := models.VolumeSnapshot{
			Region:      data.Region,
			ID:          types.StringValue(x.ID),
			VolumeID:    types.StringValue(x.VolumeID),
			Name:        types.StringValue(x.Name),
			Description: types.StringValue(x.Description),
			VolumeName:  types.StringValue(x.VolumeName),
			Status:      types.StringValue(x.Status),
			Size:        types.Int64Value(int64(x.Size)),
		}
		tfSnaps = append(tfSnaps, snp)
	}
	resp.Diagnostics.Append(data.SetVolumeSnapshots(ctx, tfSnaps)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewSnapshotDatasource() datasource.DataSource {
	return &SnapshotDatasource{}
}
