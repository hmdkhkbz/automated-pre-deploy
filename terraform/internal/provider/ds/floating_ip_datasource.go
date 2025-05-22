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

type FloatingIPDatasource struct {
	client *api.Client
}

func (f *FloatingIPDatasource) SetAPIClient(client *api.Client) {
	f.client = client
}

func (f *FloatingIPDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_floating_ips"
}

func (f *FloatingIPDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, f)
}

func (f *FloatingIPDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"floating_ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"fixed_ip": schema.StringAttribute{
							Computed: true,
						},
						"floating_ip": schema.StringAttribute{
							Computed: true,
						},
						"port_id": schema.StringAttribute{
							Computed: true,
						},
						"attached_instance_id": schema.StringAttribute{
							Computed: true,
						},
						"attached_instance_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (f *FloatingIPDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFFloatingIPDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := f.client.FIPClient.GetAllFloatingIPs(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching floating ips", err.Error())
		return
	}
	var tfIPList []models.TFFloatingIPDataSourceItem

	for _, x := range apiResp {
		ip := models.TFFloatingIPDataSourceItem{
			ID:          types.StringValue(x.ID),
			Status:      types.StringValue(x.Status),
			Description: types.StringValue(x.Description),
			FixedIP:     types.StringValue(x.FixedIPAddress),
			FloatingIP:  types.StringValue(x.FloatingIPAddress),
			PortID:      types.StringValue(x.PortID),
		}
		if x.Server != nil {
			ip.AttachedInstanceName = types.StringValue(x.Server.Name)
			ip.AttachedInstanceID = types.StringValue(x.Server.ID)
		}
		tfIPList = append(tfIPList, ip)
	}

	resp.Diagnostics.Append(data.SetFloatingIPs(ctx, tfIPList)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewFloatingIPDataSource() datasource.DataSource {
	return &FloatingIPDatasource{}
}
