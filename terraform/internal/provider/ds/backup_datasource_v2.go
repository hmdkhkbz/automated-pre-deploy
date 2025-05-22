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

type BackupV2Datasource struct {
	client *api.Client
}

func (b *BackupV2Datasource) SetAPIClient(c *api.Client) {
	b.client = c
}

func (b *BackupV2Datasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backups_v2"
}

func (b *BackupV2Datasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	misc.ConfigureDatasource(ctx, &req, resp, b)
}

func (b *BackupV2Datasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"backups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"backup_name": schema.StringAttribute{
							Computed: true,
						},
						"instance_id": schema.StringAttribute{
							Computed: true,
						},
						"instance_name": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"quota": schema.Int64Attribute{
							Computed: true,
						},
						"occupancy": schema.Int64Attribute{
							Computed: true,
						},
						"next_backup": schema.StringAttribute{
							Computed: true,
						},
						"labels": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						"details": schema.ListNestedAttribute{
							Computed: true,
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"provisioned_size": schema.Int64Attribute{
										Computed: true,
									},
									"used_size": schema.Float64Attribute{
										Computed: true,
									},
									"created_at": schema.StringAttribute{
										Computed: true,
									},
									"backup_id": schema.StringAttribute{
										Computed: true,
									},
									"status": schema.StringAttribute{
										Computed: true,
									},
									"slot_name": schema.StringAttribute{
										Computed: true,
									},
									"fail_reason": schema.StringAttribute{
										Computed: true,
										Optional: true,
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

func (b *BackupV2Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFBackupV2
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	backs, err := b.client.BackupV2.ListBackups(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching backups", err.Error())
		return
	}
	var tfBackups []models.BackupV2Item
	for _, x := range backs.Data {
		tfB := models.BackupV2Item{
			Name:         types.StringValue(x.BackupName),
			InstanceID:   types.StringValue(x.InstanceID),
			InstanceName: types.StringValue(x.InstanceName),
			Status:       types.StringValue(x.Status),
			Quota:        types.Int64Value(int64(x.Quota)),
			Occupancy:    types.Int64Value(int64(x.Occupancy)),
			NextBackup:   types.StringValue(x.NextBackup),
			Labels:       types.List{},
			Details:      types.List{},
		}
		resp.Diagnostics.Append(tfB.SetLabels(ctx, x.Labels)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiDetails, err := b.client.BackupV2.BackupDetails(ctx, data.Region.ValueString(), x.InstanceID)
		if err != nil {
			resp.Diagnostics.AddError("error fetching backup details", err.Error())
			return
		}
		resp.Diagnostics.Append(tfB.SetDetails(ctx, apiDetails.Data)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tfBackups = append(tfBackups, tfB)
	}

	resp.Diagnostics.Append(data.SetBackups(tfBackups)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func NewBackupV2Datasource() datasource.DataSource {
	return &BackupV2Datasource{}
}
