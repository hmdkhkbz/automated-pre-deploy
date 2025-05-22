package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type PlanItem struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	CpuCount         int     `json:"cpu_count"`
	Disk             int     `json:"disk"`
	DiskInBytes      int64   `json:"disk_in_bytes"`
	BandwidthInBytes int64   `json:"bandwidth_in_bytes"`
	Memory           int     `json:"memory"`
	MemoryInBytes    int64   `json:"memory_in_bytes"`
	PricePerHour     float64 `json:"price_per_hour"`
	PricePerMonth    float64 `json:"price_per_month"`
	Generation       string  `json:"generation"`
	Type             string  `json:"type"`
	Subtype          string  `json:"subtype"`
	BasePackage      string  `json:"base_package"`
	CpuShare         string  `json:"cpu_share"`
	PPS              []int   `json:"pps"`
	IOpsMaxHDD       int     `json:"iops_max_hdd"`
	IOpsMaxSSD       int     `json:"iops_max_ssd"`
	Off              string  `json:"off"`
	OffPercent       string  `json:"off_percent"`
	Throughput       int64   `json:"throughput"`
	Outbound         int64   `json:"outbound"`
}

type ServerFlavor struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	RAM   float64 `json:"ram"`
	Swap  string  `json:"swap"`
	VCPUs int32   `json:"vcpus"`
	Disk  int32   `json:"disk"`
}

type PlanList struct {
	Data []PlanItem `json:"data"`
}

type PlanClient struct {
	requester *Requester
}

func NewPlanClient(r *Requester) *PlanClient {
	return &PlanClient{
		requester: r,
	}
}

func (p *PlanClient) ListPlans(ctx context.Context, region string) (*PlanList, error) {
	uri := fmt.Sprintf("%s/%s/sizes", basePath, region)
	data, err := p.requester.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret PlanList
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
