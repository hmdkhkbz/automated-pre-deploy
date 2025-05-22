package ds

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/utl"
)

type ImageDistroDatasource struct {
	client *api.Client
}

func NewImageDistroDatasource() datasource.DataSource {
	return &ImageDistroDatasource{}
}

func (i *ImageDistroDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

func (i *ImageDistroDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"image_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "image type",
			},
			"region": schema.StringAttribute{
				Required: true,
			},
			"distributions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"distro_name": schema.StringAttribute{
							Computed: true,
						},
						"disk": schema.Int64Attribute{
							Computed: true,
						},
						"ram": schema.Int64Attribute{
							Computed: true,
						},
						"ssh_key": schema.BoolAttribute{
							Computed: true,
						},
						"password": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (i *ImageDistroDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		utl.DataSourceConfigureError(&req, resp)
		return
	}
	i.client = client
}

func (i *ImageDistroDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ImageDistroListModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiRet, err := i.client.Img.ListImages(ctx, data.Region.ValueString(), data.ImgType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error fetching image list",
			err.Error(),
		)
		return
	}

	for _, group := range apiRet.Data {
		for _, x := range group.Images {
			data.Distributions = append(data.Distributions, models.ImageItem{
				ID:         types.StringValue(x.ID),
				Name:       types.StringValue(x.Name),
				DistroName: types.StringValue(x.DistroName),
				Disk:       types.Int64Value(int64(x.Disk)),
				Ram:        types.Int64Value(int64(x.Ram)),
				SSHKey:     types.BoolValue(x.SSHKey),
				Password:   types.BoolValue(x.SSHPassword),
			})
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
