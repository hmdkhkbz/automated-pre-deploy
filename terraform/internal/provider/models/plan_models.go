package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type TFPlanItem struct {
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	CPUCount         types.Int64   `tfsdk:"cpu_count"`
	Disk             types.Int64   `tfsdk:"disk"`
	DiskInBytes      types.Int64   `tfsdk:"disk_in_bytes"`
	BandwidthInBytes types.Int64   `tfsdk:"bandwidth_in_bytes"`
	Memory           types.Int64   `tfsdk:"memory"`
	MemoryInBytes    types.Int64   `tfsdk:"memory_in_bytes"`
	PricePerHour     types.Float64 `tfsdk:"price_per_hour"`
	PricePerMonth    types.Float64 `tfsdk:"price_per_month"`
	Generation       types.String  `tfsdk:"generation"`
	Type             types.String  `tfsdk:"type"`
	Subtype          types.String  `tfsdk:"subtype"`
	BasePackage      types.String  `tfsdk:"base_package"`
	CPUShare         types.String  `tfsdk:"cpu_share"`
	PPS              []types.Int64 `tfsdk:"pps"`
	IOPSMaxHDD       types.Int64   `tfsdk:"iops_max_hdd"`
	IOPSMaxSSD       types.Int64   `tfsdk:"iops_max_ssd"`
	Off              types.String  `tfsdk:"off"`
	OffPercent       types.String  `tfsdk:"off_percent"`
	Throughput       types.Int64   `tfsdk:"throughput"`
	Outbound         types.Int64   `tfsdk:"outbound"`
}

type TFPlanListDataModel struct {
	Region types.String `tfsdk:"region"`
	Plans  []TFPlanItem `tfsdk:"plans"`
}
