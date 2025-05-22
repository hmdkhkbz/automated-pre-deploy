package api

import (
	"context"
	"encoding/json"
	"fmt"
)


type DedicatedServerClient struct{
	requester *Requester
}

func NewDedicatedServerClient(r *Requester) *DedicatedServerClient {
	return &DedicatedServerClient{
		requester: r,
	}
}

type DedicatedServerList struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TypeID      string   `json:"type_id"`
	Sockets     int      `json:"sockets"`
	VCPUs       int      `json:"vcpus"`
	VCPUsUsed   int      `json:"vcpus_used"`
	Memory      int      `json:"memory"`
	MemoryUsed  int      `json:"memory_used"`
	Disk        int      `json:"disk"`
	DiskUsed    int      `json:"disk_used"`
	Instances   int      `json:"instances"`
	Status      string   `json:"status"`
	ClusterName string   `json:"cluster_name"`
	CreatedAt   int64    `json:"created_at"`
	Labels      []string `json:"labels"`
}

func (i *DedicatedServerClient) ListDedicatedServers(ctx context.Context, region string) ([]DedicatedServerList, error) {
	type dedicatedServerListResponse struct {
		Data []DedicatedServerList `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/dedicated-servers/servers", basePath, region)

	data, err := i.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp dedicatedServerListResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}