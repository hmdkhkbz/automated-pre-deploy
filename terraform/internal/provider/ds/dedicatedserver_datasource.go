package ds

import (
	"context"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DedicatedServerDatasource struct {
	client *api.Client
}

func (b *DedicatedServerDatasource) SetAPIClient(c *api.Client) {
	b.client = c
}

func (b *DedicatedServerDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_servers"
}

func (b *DedicatedServerDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, b)
}

func (b *DedicatedServerDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"dedicated_servers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type_id": schema.StringAttribute{
							Computed: true,
						},
						"sockets": schema.Int64Attribute{
							Computed: true,
						},
						"vcpus": schema.Int64Attribute{
							Computed: true,
						},
						"vcpus_used": schema.Int64Attribute{
							Computed: true,
						},
						"memory": schema.Int64Attribute{
							Computed: true,
						},
						"memory_used": schema.Int64Attribute{
							Computed: true,
						},
						"disk": schema.Int64Attribute{
							Computed: true,
						},
						"disk_used": schema.Int64Attribute{
							Computed: true,
						},
						"instances": schema.Int64Attribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"cluster_name": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"labels": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (b *DedicatedServerDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFDedicatedServer
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dedicatedServers, err := b.client.DedicatedServer.ListDedicatedServers(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching dedicated servers", err.Error())
		return
	}
	var tfDedicatedServers []models.DedicatedServerItem
	for _, x := range dedicatedServers {
		createdAt := time.Unix(0, x.CreatedAt)

		item := models.DedicatedServerItem{
			ID:          types.StringValue(x.ID),
			Name:        types.StringValue(x.Name),
			TypeID:      types.StringValue(x.TypeID),
			Sockets:     types.Int64Value(int64(x.Sockets)),
			VCPUs:       types.Int64Value(int64(x.VCPUs)),
			VCPUsUsed:   types.Int64Value(int64(x.VCPUsUsed)),
			Memory:      types.Int64Value(int64(x.Memory)),
			MemoryUsed:  types.Int64Value(int64(x.MemoryUsed)),
			Disk:        types.Int64Value(int64(x.Disk)),
			DiskUsed:    types.Int64Value(int64(x.DiskUsed)),
			Instances:   types.Int64Value(int64(x.Instances)),
			Status:      types.StringValue(x.Status),
			ClusterName: types.StringValue(x.ClusterName),
			CreatedAt:   types.StringValue(createdAt.String()),
		}

		resp.Diagnostics.Append(item.SetLabels(ctx, x.Labels)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tfDedicatedServers = append(tfDedicatedServers, item)
	}

	resp.Diagnostics.Append(data.SetDedicatedServers(tfDedicatedServers)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewDedicatedServerDatasource() datasource.DataSource {
	return &DedicatedServerDatasource{}
}
