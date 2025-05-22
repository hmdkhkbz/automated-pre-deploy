package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type VolumeV2CreateRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type VolumeV2DeleteRequest struct {
	VolumeIDs []string `json:"volume_ids"`
}

type Details struct {
	Data []DetailsData `json:"data"`
}

type DetailsData struct {
	CreatedAt int64  `json:"created_at"`
	AZ        string `json:"availability_zone"`
	Path      string `json:"path"`
	IOPS      int    `json:"iops"`
}

type List struct {
	Data []ListData `json:"data"`
}

type ListData struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Size         int      `json:"size"`
	InstanceName string   `json:"instance_name"`
	VolumeType   string   `json:"type"`
	Status       string   `json:"status"`
	Labels       []string `json:"labels"`
}

type Inquiry struct {
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Status       string `json:"status"`
	InstanceName string `json:"instance_name"`
}

type EditLabelsRequest struct {
	Labels []string `json:"labels"`
}

type EditNameRequest struct {
	Name string `json:"name"`
}

type AttachRequest struct {
	InstanceID string `json:"instance_id"`
	MountPoint string `json:"device"`
}

type ResizeRequest struct {
	Size int `json:"size"`
}

type VolumeV2Response struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	VolumeID string `json:"volume_id,omitempty"`
}

type VolumeV2Client struct {
	r *Requester
}

func NewVolumeV2Client(r *Requester) *VolumeV2Client {
	ret := &VolumeV2Client{
		r: r,
	}
	return ret
}

func (v2 *VolumeV2Client) Create(ctx context.Context, region string, req *VolumeV2CreateRequest) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/create", bpV2, region)
	data, err := v2.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) Delete(ctx context.Context, region string, req *VolumeV2DeleteRequest) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/delete", bpV2, region)
	data, err := v2.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) Details(ctx context.Context, region, volumeID string) (*Details, error) {
	uri := fmt.Sprintf("%s/volume/%s/details/%s", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret Details
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil

}

func (v2 *VolumeV2Client) List(ctx context.Context, region string) (*List, error) {
	uri := fmt.Sprintf("%s/volume/%s/list", bpV2, region)
	data, err := v2.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret List
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) Inquire(ctx context.Context, region, volumeID string) (*Inquiry, error) {
	uri := fmt.Sprintf("%s/volume/%s/inquiry/%s", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret Inquiry
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) EditLabels(ctx context.Context, region, volumeID string, req *EditLabelsRequest) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/%s/labels", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "PUT", uri, req)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) EditName(ctx context.Context, region, volumeID string, req *EditNameRequest) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/%s/name", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "PUT", uri, req)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) Attach(ctx context.Context, region, volumeID string, req *AttachRequest) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/%s/attach", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) Detach(ctx context.Context, region, volumeID string) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/%s/detach", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "POST", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) Resize(ctx context.Context, region, volumeID string, req *ResizeRequest) (*VolumeV2Response, error) {
	uri := fmt.Sprintf("%s/volume/%s/resize/%s", bpV2, region, volumeID)
	data, err := v2.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret VolumeV2Response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (v2 *VolumeV2Client) GetVolumeByID(ctx context.Context, region, volumeID string) (*ListData, error) {
	vols, err := v2.List(ctx, region)
	if err != nil {
		return nil, err
	}

	for _, x := range vols.Data {
		if x.ID == volumeID {
			return &x, nil
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "volume not found",
	}
}
