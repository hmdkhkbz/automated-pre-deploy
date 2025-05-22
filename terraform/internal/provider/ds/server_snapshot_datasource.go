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

type ServerSnapshotDatasource struct {
	client *api.Client
}

func (s *ServerSnapshotDatasource) SetAPIClient(c *api.Client) {
	s.client = c
}

func (s *ServerSnapshotDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_snapshot"
}

func (s *ServerSnapshotDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, s)
}

func (s *ServerSnapshotDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"size": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"volume_id": schema.StringAttribute{
							Computed: true,
						},
						"image_id": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"server_id": schema.StringAttribute{
							Computed: true,
						},
						"server_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (s *ServerSnapshotDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFServerSnapshotDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := s.client.SnapshotClient.ListServerSnapshots(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching server snapshots", err.Error())
		return
	}

	var tfData []models.ServerSnapshot
	for _, x := range apiResp {
		tfs := models.ServerSnapshot{
			Region:      data.Region,
			ID:          types.StringValue(x.ID),
			Size:        types.Int64Value(int64(x.Size)),
			Name:        types.StringValue(x.Name),
			Description: types.StringValue(x.Description),
			VolumeID:    types.StringValue(x.VolumeID),
			ImageID:     types.StringValue(x.ImageID),
			Status:      types.StringValue(x.Status),
			ServerID:    types.StringValue(x.ServerID),
			ServerName:  types.StringValue(x.ServerName),
		}
		tfData = append(tfData, tfs)
	}
	resp.Diagnostics.Append(data.SetSnapshots(ctx, tfData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewServerSnapshotDatasource() datasource.DataSource {
	return &ServerSnapshotDatasource{}
}
