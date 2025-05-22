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

type VolumeDatasource struct {
	client *api.Client
}

func (v *VolumeDatasource) SetAPIClient(client *api.Client) {
	v.client = client
}

func (v *VolumeDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volumes"
}

func (v *VolumeDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, v)
}

func (v *VolumeDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"volumes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"size": schema.Int64Attribute{
							Computed: true,
						},
						"attachments": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"device": schema.StringAttribute{
										Computed: true,
									},
									"volume_id": schema.StringAttribute{
										Computed: true,
									},
									"attachment_id": schema.StringAttribute{
										Computed: true,
									},
									"server_id": schema.StringAttribute{
										Computed: true,
									},
									"server_name": schema.StringAttribute{
										Computed: true,
									},
									"host_name": schema.StringAttribute{
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

func (v *VolumeDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFVolumeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, err := v.client.Volume.ListVolumes(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching volumes", err.Error())
		return
	}

	var tfVols []models.TFVolumeModel
	for _, x := range apiResp {
		tfV := models.TFVolumeModel{
			Region:      data.Region,
			ID:          types.StringValue(x.ID),
			Name:        types.StringValue(x.Name),
			Description: types.StringValue(x.Description),
			Size:        types.Int64Value(int64(x.Size)),
		}
		var atts []models.TFVolumeAttachment
		for _, a := range x.Attachments {
			tfAtt := models.TFVolumeAttachment{
				ID:           types.StringValue(a.ID),
				Device:       types.StringValue(a.Device),
				VolumeID:     types.StringValue(a.VolumeID),
				AttachmentID: types.StringValue(a.AttachmentID),
				ServerID:     types.StringValue(a.ServerID),
				ServerName:   types.StringValue(a.ServerName),
				HostName:     types.StringValue(a.HostName),
			}
			atts = append(atts, tfAtt)
		}
		resp.Diagnostics.Append(tfV.SetAttachments(ctx, atts)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tfVols = append(tfVols, tfV)
	}
	resp.Diagnostics.Append(data.SetVolumes(ctx, tfVols)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewVolumeDataSource() datasource.DataSource {
	return &VolumeDatasource{}
}
