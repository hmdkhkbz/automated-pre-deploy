package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type ServerVolume struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int32  `json:"size"`
	Type        string `json:"type"`
}

type Attachment struct {
	ID           string `json:"id"`
	Device       string `json:"device"`
	ServerID     string `json:"server_id"`
	ServerName   string `json:"server_name"`
	VolumeID     string `json:"volume_id"`
	AttachmentID string `json:"attachment_id"`
	AttachedAt   string `json:"attached_at"`
	HostName     string `json:"host_name"`
}

type VolumeDetails struct {
	ID             string       `json:"id"`
	Size           int32        `json:"size"`
	Status         string       `json:"status"`
	CreatedAt      string       `json:"created_at"`
	Description    string       `json:"description"`
	VolumeTypeName string       `json:"volume_type_name"`
	SnapshotID     string       `json:"snapshot_id"`
	SourceVolumeID string       `json:"source_volume_id"`
	Bootable       string       `json:"bootable"`
	Name           string       `json:"name"`
	Attachments    []Attachment `json:"attachments"`
}

type VolumeAttachDetach struct {
	ServerID string `json:"server_id"`
	VolumeID string `json:"volume_id"`
}

type VolumeClient struct {
	r *Requester
}

func NewVolumeClient(r *Requester) *VolumeClient {
	return &VolumeClient{
		r: r,
	}
}

func (v *VolumeClient) CreateVolume(ctx context.Context, region string, req *ServerVolume) (*VolumeDetails, error) {
	type createResp struct {
		Data    *VolumeDetails `json:"data"`
		Message string         `json:"message"`
	}
	url := fmt.Sprintf("%s/%s/volumes", basePath, region)

	data, err := v.r.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}
	var resp createResp

	err = json.Unmarshal(data, &resp)
	return resp.Data, nil
}

func (v *VolumeClient) DeleteVolume(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/volumes/%s", basePath, region, id)

	_, err := v.r.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (v *VolumeClient) AttachVolume(ctx context.Context, region string, req *VolumeAttachDetach) (*Attachment, error) {
	url := fmt.Sprintf("%s/%s/volumes/attach", basePath, region)

	data, err := v.r.DoRequest(ctx, "PATCH", url, req)
	if err != nil {
		return nil, err
	}
	var resp Attachment
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (v *VolumeClient) DetachVolume(ctx context.Context, region string, req *VolumeAttachDetach) error {
	url := fmt.Sprintf("%s/%s/volumes/detach", basePath, region)

	_, err := v.r.DoRequest(ctx, "PATCH", url, req)
	return err
}

func (v *VolumeClient) UpdateVolume(ctx context.Context, region, id string, req *ServerVolume) error {
	url := fmt.Sprintf("%s/%s/volumes/%s", basePath, region, id)

	_, err := v.r.DoRequest(ctx, "PATCH", url, req)
	return err
}

func (v *VolumeClient) ListVolumes(ctx context.Context, region string) ([]*VolumeDetails, error) {
	type listVolumeResponse struct {
		Data []*VolumeDetails `json:"data"`
	}

	url := fmt.Sprintf("%s/%s/volumes", basePath, region)
	data, err := v.r.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp listVolumeResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (v *VolumeClient) GetServerVolumes(ctx context.Context, region, serverID string) ([]string, error) {
	allV, err := v.ListVolumes(ctx, region)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, x := range allV {
		for _, a := range x.Attachments {
			if a.ServerID == serverID {
				ret = append(ret, x.ID)
				continue
			}
		}
	}
	return ret, nil
}

func (v *VolumeClient) GetVolume(ctx context.Context, region, id string) (*VolumeDetails, error) {
	l, err := v.ListVolumes(ctx, region)
	if err != nil {
		return nil, err
	}
	for _, v := range l {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "volume not found",
	}
}
