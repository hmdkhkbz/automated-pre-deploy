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

type InstanceSnapshotDatasourceV2 struct {
	client *api.Client
}

func (i *InstanceSnapshotDatasourceV2) SetAPIClient(c *api.Client) {
	i.client = c
}

func (i *InstanceSnapshotDatasourceV2) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_snapshots_v2"
}

func (i *InstanceSnapshotDatasourceV2) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, i)
}

func (i *InstanceSnapshotDatasourceV2) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"details": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed: true,
						},
						"instance_name": schema.StringAttribute{
							Computed: true,
						},
						"snapshots": schema.ListNestedAttribute{
							Computed: true,
							Optional: true,
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
										Computed:    true,
										Optional:    true,
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

func (i *InstanceSnapshotDatasourceV2) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFInstanceSnapshotDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := i.client.BackupV2.ListInstanceSnapshots(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching instance snapshot list", err.Error())
		return
	}
	var tfData []models.TFInstanceSnapshotListItem
	for _, x := range list.Data {
		item := models.TFInstanceSnapshotListItem{
			InstanceID:   types.StringValue(x.InstanceID),
			InstanceName: types.StringValue(x.InstanceName),
		}
		details, err := i.client.BackupV2.InstanceSnapshotDetails(ctx, data.Region.ValueString(), x.InstanceID)
		if err != nil {
			resp.Diagnostics.AddError("error fetching instance snapshot details", err.Error())
			return
		}
		resp.Diagnostics.Append(item.SetSnapshots(ctx, details.Snapshots)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tfData = append(tfData, item)
	}

	resp.Diagnostics.Append(data.SetDetails(tfData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewInstanceSnapshotDatasourceV2() datasource.DataSource {
	return &InstanceSnapshotDatasourceV2{}
}
