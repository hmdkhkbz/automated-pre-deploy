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

type NetworkDatasource struct {
	client *api.Client
}

func (n *NetworkDatasource) SetAPIClient(c *api.Client) {
	n.client = c
}

func (n *NetworkDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks"
}

func (n *NetworkDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, n)
}

func (n *NetworkDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"public": schema.BoolAttribute{
						Optional: true,
					},
					"name": schema.StringAttribute{
						Optional: true,
					},
					"network_id": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"networks": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"subnet_id": schema.StringAttribute{
							Computed: true,
						},
						"network_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"dhcp_range": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"start": schema.StringAttribute{
									Computed: true,
								},
								"end": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"dns_servers": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"enable_dhcp": schema.BoolAttribute{
							Computed: true,
						},
						"enable_gateway": schema.BoolAttribute{
							Computed: true,
						},
						"gateway_ip": schema.StringAttribute{
							Computed: true,
						},
						"cidr": schema.StringAttribute{
							Computed: true,
						},
						"shared": schema.BoolAttribute{
							Computed: true,
						},
						"mtu": schema.Int64Attribute{
							Computed: true,
						},
						"servers": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"ip": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"port_id": schema.StringAttribute{
										Computed: true,
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

func (n *NetworkDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFNetworkDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := n.client.Subnet.GetAllNetworks(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching networks", err.Error())
		return
	}

	var tfNets []models.TFSubnetDatasourceModel
	for _, v := range apiResp {
		tfSub, d := models.TFSubnetDatasourceModelFromAPINetwork(ctx, v)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		tfNets = append(tfNets, *tfSub)
	}
	resp.Diagnostics.Append(data.SetNetworks(ctx, tfNets)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewNetworkDatasource() datasource.DataSource {
	return &NetworkDatasource{}
}
