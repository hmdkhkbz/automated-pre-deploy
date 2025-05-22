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

type PlanDatasource struct {
	client *api.Client
}

func NewPlanDatasource() datasource.DataSource {
	return &PlanDatasource{}
}

func (p *PlanDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plans"
}

func (p *PlanDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"plans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"cpu_count": schema.Int64Attribute{
							Computed: true,
						},
						"disk": schema.Int64Attribute{
							Computed: true,
						},
						"disk_in_bytes": schema.Int64Attribute{
							Computed: true,
						},
						"bandwidth_in_bytes": schema.Int64Attribute{
							Computed: true,
						},
						"memory": schema.Int64Attribute{
							Computed: true,
						},
						"memory_in_bytes": schema.Int64Attribute{
							Computed: true,
						},
						"price_per_hour": schema.Float64Attribute{
							Computed: true,
						},
						"price_per_month": schema.Float64Attribute{
							Computed: true,
						},
						"generation": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"subtype": schema.StringAttribute{
							Computed: true,
						},
						"base_package": schema.StringAttribute{
							Computed: true,
						},
						"cpu_share": schema.StringAttribute{
							Computed: true,
						},
						"pps": schema.ListAttribute{
							Computed:    true,
							ElementType: types.Int64Type,
						},
						"iops_max_hdd": schema.Int64Attribute{
							Computed: true,
						},
						"iops_max_ssd": schema.Int64Attribute{
							Computed: true,
						},
						"off": schema.StringAttribute{
							Computed: true,
						},
						"off_percent": schema.StringAttribute{
							Computed: true,
						},
						"throughput": schema.Int64Attribute{
							Computed: true,
						},
						"outbound": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (p *PlanDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		utl.DataSourceConfigureError(&req, resp)
		return
	}
	p.client = client
}

func (p *PlanDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.TFPlanListDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plans, err := p.client.Pln.ListPlans(ctx, data.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching plans", err.Error())
		return
	}
	for _, x := range plans.Data {
		tfPlan := models.TFPlanItem{
			ID:               types.StringValue(x.ID),
			Name:             types.StringValue(x.Name),
			CPUCount:         types.Int64Value(int64(x.CpuCount)),
			Disk:             types.Int64Value(int64(x.Disk)),
			DiskInBytes:      types.Int64Value(x.DiskInBytes),
			BandwidthInBytes: types.Int64Value(x.BandwidthInBytes),
			Memory:           types.Int64Value(int64(x.Memory)),
			MemoryInBytes:    types.Int64Value(x.MemoryInBytes),
			PricePerHour:     types.Float64Value(x.PricePerHour),
			PricePerMonth:    types.Float64Value(x.PricePerMonth),
			Generation:       types.StringValue(x.Generation),
			Type:             types.StringValue(x.Type),
			Subtype:          types.StringValue(x.Subtype),
			BasePackage:      types.StringValue(x.BasePackage),
			CPUShare:         types.StringValue(x.CpuShare),
			PPS:              nil,
			IOPSMaxHDD:       types.Int64Value(int64(x.IOpsMaxHDD)),
			IOPSMaxSSD:       types.Int64Value(int64(x.IOpsMaxSSD)),
			Off:              types.StringValue(x.Off),
			OffPercent:       types.StringValue(x.OffPercent),
			Throughput:       types.Int64Value(x.Throughput),
			Outbound:         types.Int64Value(x.Outbound),
		}
		for _, pps := range x.PPS {
			tfPlan.PPS = append(tfPlan.PPS, types.Int64Value(int64(pps)))
		}
		data.Plans = append(data.Plans, tfPlan)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

}
