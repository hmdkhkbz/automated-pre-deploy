package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type SnapshotRequest struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type SnapshotResponse struct {
	ID                string `json:"id"`
	Size              int32  `json:"size"`
	Status            string `json:"status"`
	VolumeID          string `json:"volume_id"`
	VolumeName        string `json:"volume_name"`
	Description       string `json:"description"`
	Name              string `json:"name"`
	CreatedAt         string `json:"created_at"`
	ServerID          string `json:"server_id"`
	ServerName        string `json:"server_name"`
	ImageID           string `json:"image_id"`
	RevertedOn        string `json:"reverted_on"`
	RealSize          int    `json:"real_size"`
	RealSizeAvailable bool   `json:"real_size_status"`
	Type              string `json:"type"`
}

type UpdateVolumeSnapshot struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	InstanceName string `json:"instance_name,omitempty"`
	VolumeName   string `json:"volume_name,omitempty"`
}

type SnapshotClient struct {
	requester *Requester
}

func NewSnapshotClient(r *Requester) *SnapshotClient {
	return &SnapshotClient{requester: r}
}

func (s *SnapshotClient) CreateVolumeSnapshot(ctx context.Context, region, volumeID string, req *SnapshotRequest) (*SnapshotResponse, error) {
	url := fmt.Sprintf("%s/%s/snapshots/volumes/%s/", basePath, region, volumeID)
	data, err := s.requester.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[SnapshotResponse]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SnapshotClient) CreateServerSnapshot(ctx context.Context, region, serverID string, req *SnapshotRequest) (*SnapshotResponse, error) {
	url := fmt.Sprintf("%s/%s/volumes/%s/snapshot", basePath, region, serverID)
	data, err := s.requester.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[SnapshotResponse]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SnapshotClient) UpdateVolumeSnapshot(ctx context.Context, region, snapshotID string, req *UpdateVolumeSnapshot) (*SnapshotResponse, error) {
	url := fmt.Sprintf("%s/%s/volumes/%s/snapshot", basePath, region, snapshotID)
	data, err := s.requester.DoRequest(ctx, "PUT", url, req)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[SnapshotResponse]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (s *SnapshotClient) ListSnapshots(ctx context.Context, region string) ([]SnapshotResponse, error) {
	url := fmt.Sprintf("%s/%s/volumes/snapshots", basePath, region)
	data, err := s.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp DataResponse[[]SnapshotResponse]
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *SnapshotClient) listByType(ctx context.Context, region, sType string) ([]SnapshotResponse, error) {
	all, err := s.ListSnapshots(ctx, region)
	if err != nil {
		return nil, err
	}
	var ret []SnapshotResponse
	for _, x := range all {
		if x.Type == sType {
			ret = append(ret, x)
		}
	}
	return ret, nil
}

func (s *SnapshotClient) ListVolumeSnapshots(ctx context.Context, region string) ([]SnapshotResponse, error) {
	return s.listByType(ctx, region, "VOLUME")
}

func (s *SnapshotClient) ListServerSnapshots(ctx context.Context, region string) ([]SnapshotResponse, error) {
	return s.listByType(ctx, region, "SERVER")
}

func (s *SnapshotClient) GetSnapshotByID(ctx context.Context, region, snapshotID string) (*SnapshotResponse, error) {
	all, err := s.ListSnapshots(ctx, region)
	if err != nil {
		return nil, err
	}

	for _, x := range all {
		if x.ID == snapshotID {
			return &x, nil
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "snapshot not found",
		URL:     "",
		Errors:  nil,
	}
}

func (s *SnapshotClient) DeleteSnapshot(ctx context.Context, region, snapshotID string) error {
	url := fmt.Sprintf("%s/%s/volumes/%s/snapshot", basePath, region, snapshotID)
	_, err := s.requester.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (s *SnapshotClient) Revert(ctx context.Context, region, serverID, snapshotID string) error {
	type revertReq struct {
		SnapshotID string `json:"snapshot_id"`
	}
	url := fmt.Sprintf("%s/%s/snapshots/%s/revert", basePath, region, serverID)
	req := &revertReq{
		SnapshotID: snapshotID,
	}
	_, err := s.requester.DoRequest(ctx, "PUT", url, req)
	return err
}

func (s *SnapshotClient) CreatePersonalImage(ctx context.Context, region, name, snapshotID string) (*VolumeDetails, error) {
	type cPersonalReq struct {
		Name string `json:"name"`
	}
	url := fmt.Sprintf("%s/%s/volumes/snapshots/%s/os-volume", basePath, region, snapshotID)
	req := cPersonalReq{
		Name: name,
	}
	var ret DataResponse[VolumeDetails]
	resp, err := s.requester.DoRequest(ctx, "POST", url, &req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	return &ret.Data, nil
}
