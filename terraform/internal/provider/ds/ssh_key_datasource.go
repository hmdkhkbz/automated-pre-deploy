package ds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
)

type SSHKeyDatasource struct {
	client *api.Client
}

func (s *SSHKeyDatasource) SetAPIClient(client *api.Client) {
	s.client = client
}

func (s *SSHKeyDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_keys"
}

func (s *SSHKeyDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, s)
}

func (s *SSHKeyDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"keys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"public_key": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (s *SSHKeyDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFSSHKeyDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	keys, err := s.client.SSHClient.GetSSHKeys(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching ssh keys", err.Error())
		return
	}
	resp.Diagnostics.Append(data.SetKeys(ctx, keys)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewSSHKeyDatasource() datasource.DataSource {
	return &SSHKeyDatasource{}
}
