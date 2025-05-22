package ds

import (
	"context"
	"terraform-provider-hashicups-pf/internal/api"

	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServerGroupDatasource struct {
	client *api.Client
}

func (s *ServerGroupDatasource) SetAPIClient(c *api.Client) {
	s.client = c
}

func (s *ServerGroupDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_groups"
}

func (s *ServerGroupDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, s)
}

func (s *ServerGroupDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"server_groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"policies": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"members": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (s *ServerGroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFServerGroupDatasourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := s.client.ServerGroup.ListServerGroups(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching server groups", err.Error())
		return
	}

	resp.Diagnostics.Append(data.SetServerGroups(ctx, apiResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewServerGroupDatasource() datasource.DataSource {
	return &ServerGroupDatasource{}
}