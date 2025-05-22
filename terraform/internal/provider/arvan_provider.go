package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/ds"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/provider/rs"
)

type ArvanProvider struct {
	version string
}

func NewArvanProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ArvanProvider{version: version}
	}
}

func (p *ArvanProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "arvan"
	resp.Version = p.version
}

func (p *ArvanProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "An API key that is acquired from MachineUser menu in ArvanCloud's panel",
			},
		},
	}
}

func (p *ArvanProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data models.ArvanProviderDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiC := api.NewClient(data.ApiKey.ValueString())
	resp.ResourceData = apiC
	resp.DataSourceData = apiC

}

func (p *ArvanProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		rs.NewInstanceResource,
		rs.NewVolumeResource,
		rs.NewNetworkResource,
		rs.NewSecurityGroupResource,
		rs.NewFloatingIPResource,
		rs.NewVolumeSnapshotResource,
		rs.NewServerSnapshotResource,
		rs.NewPersonalImageResource,
		rs.NewVolumeV2Resource,
		rs.NewVolumeSnapshotV2Resource,
		rs.NewInstanceSnapshotResource,
	}
}

func (p *ArvanProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ds.NewImageDistroDatasource,
		ds.NewPlanDatasource,
		ds.NewInstanceDatasource,
		ds.NewVolumeDataSource,
		ds.NewFloatingIPDataSource,
		ds.NewSecurityGroupDatasource,
		ds.NewNetworkDatasource,
		ds.NewSnapshotDatasource,
		ds.NewServerSnapshotDatasource,
		ds.NewSSHKeyDatasource,
		ds.NewVolumeV2Datasource,
		ds.NewBackupV2Datasource,
		ds.NewVolumeSnapshotV2Datasource,
		ds.NewInstanceSnapshotDatasourceV2,
		ds.NewServerGroupDatasource,
		ds.NewDedicatedServerDatasource,
	}
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ArvanProvider{
			version: version,
		}
	}
}
