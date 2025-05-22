package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type SecGroupName struct {
	Name string `json:"name"`
}

type InstanceCreateRequest struct {
	Name              string         `json:"name"`
	Count             int            `json:"count"`
	NetworkID         string         `json:"network_id"`
	ImageID           string         `json:"image_id"`
	NetworkIDs        []string       `json:"network_ids"`
	FlavorID          string         `json:"flavor_id"`
	SecurityGroups    []SecGroupName `json:"security_groups"`
	SSHKey            bool           `json:"ssh_key"`
	CreateType        string         `json:"create_type"`
	DiskSize          int            `json:"disk_size"`
	InitScript        string         `json:"init_script"`
	HAEnabled         *bool          `json:"ha_enabled"`
	ServerVolumes     []ServerVolume `json:"server_volumes"`
	IsSandbox         bool           `json:"is_sandbox"`
	OSVolumeID        string         `json:"os_volume_id"`
	KeyName           interface{}    `json:"key_name"`
	ServerGroupID     string         `json:"server_group_id"`
	DedicatedServerID string         `json:"dedicated_server_id"`
	SnapshotID        string         `json:"snapshot_id"`
	EnableIPv4        bool           `json:"enable_ipv4"`
	EnableIPv6        bool           `json:"enable_ipv6"`
}

type ServerDetail struct {
	ID                string                      `json:"id"`
	TaskID            string                      `json:"task_id"`
	Name              string                      `json:"name"`
	Flavor            *ServerFlavor               `json:"flavor"`
	Status            string                      `json:"status"`
	Image             *ServerImage                `json:"image"`
	Created           string                      `json:"created"`
	Password          string                      `json:"password"`
	TaskState         *string                     `json:"task_state"`
	KeyName           string                      `json:"key_name"`
	ArNext            string                      `json:"ar_next"`
	SecurityGroups    []*SecurityGroup            `json:"security_groups"`
	Addresses         map[string][]*ServerAddress `json:"addresses"`
	Tags              []*Tag                      `json:"tags"`
	HAEnabled         bool                        `json:"ha_enabled"`
	ClusterID         string                      `json:"cluster_id"`
	DedicatedServerID string                      `json:"dedicated_server_id"`
}

type InstanceCreateResponse struct {
	Message string       `json:"message"`
	Data    ServerDetail `json:"data"`
}

type InstanceClient struct {
	requester *Requester
}

func NewInstanceClient(r *Requester) *InstanceClient {
	return &InstanceClient{
		requester: r,
	}
}

func (i *InstanceClient) CreateInstance(ctx context.Context, region string, req *InstanceCreateRequest) (*InstanceCreateResponse, error) {
	url := fmt.Sprintf("%s/%s/servers", basePath, region)

	data, err := i.requester.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}
	var resp InstanceCreateResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (i *InstanceClient) CreateInstanceAsync(ctx context.Context, region string, req *InstanceCreateRequest) (*InstanceCreateResponse, error) {
	url := fmt.Sprintf("%s/%s/servers?async=true", basePath, region)

	data, err := i.requester.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}
	var resp InstanceCreateResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (i *InstanceClient) InquiryInstance(ctx context.Context, region, taskID string) (*ServerDetail, error) {
	type getServerResponse struct {
		Data *ServerDetail `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/servers/inquiry/%s", basePath, region, taskID)

	data, err := i.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp getServerResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil

}

func (i *InstanceClient) ListInstances(ctx context.Context, region string) ([]ServerDetail, error) {
	type instanceListResponse struct {
		Data []ServerDetail `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/servers", basePath, region)

	data, err := i.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp instanceListResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (i *InstanceClient) GetInstance(ctx context.Context, region, id string) (*ServerDetail, error) {
	type getServerResponse struct {
		Data *ServerDetail `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/servers/%s", basePath, region, id)

	data, err := i.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp getServerResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil

}

func (i *InstanceClient) DeleteInstance(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/servers/%s?forceDelete=true", basePath, region, id)
	_, err := i.requester.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (i *InstanceClient) ResizeInstance(ctx context.Context, region, id, flavorID string) error {
	type resizeReq struct {
		FlavorID string `json:"flavor_id"`
	}
	url := fmt.Sprintf("%s/%s/servers/%s/resize", basePath, region, id)
	req := resizeReq{
		FlavorID: flavorID,
	}
	_, err := i.requester.DoRequest(ctx, "POST", url, &req)
	return err

}

func (i *InstanceClient) ResizeRootVolume(ctx context.Context, region, id string, newSize int64) error {
	type resizeRootReq struct {
		NewSize int64 `json:"new_size"`
	}
	url := fmt.Sprintf("%s/%s/servers/%s/resizeRoot", basePath, region, id)
	req := resizeRootReq{
		NewSize: newSize,
	}
	_, err := i.requester.DoRequest(ctx, "PUT", url, &req)
	return err
}

func (i *InstanceClient) PowerOffInstance(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/servers/%s/power-off", basePath, region, id)
	_, err := i.requester.DoRequest(ctx, "POST", url, nil)
	return err
}

func (i *InstanceClient) PowerOnInstance(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/servers/%s/power-on", basePath, region, id)
	_, err := i.requester.DoRequest(ctx, "POST", url, nil)
	return err
}

func (i *InstanceClient) RebootInstance(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/servers/%s/reboot", basePath, region, id)
	_, err := i.requester.DoRequest(ctx, "POST", url, nil)
	return err
}

func (i *InstanceClient) HardRebootInstance(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/servers/%s/hard-reboot", basePath, region, id)
	_, err := i.requester.DoRequest(ctx, "POST", url, nil)
	return err
}

func (i *InstanceClient) RenameInstance(ctx context.Context, region, id, name string) error {
	type renameReq struct {
		Name string `json:"name"`
	}
	url := fmt.Sprintf("%s/%s/servers/%s/rename", basePath, region, id)
	req := renameReq{
		Name: name,
	}
	_, err := i.requester.DoRequest(ctx, "POST", url, &req)
	return err
}
