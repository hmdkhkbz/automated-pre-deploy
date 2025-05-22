package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type ImageDistroItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	DistroName    string `json:"distribution_name"`
	OSDescription string `json:"os_description"`
	Disk          int    `json:"disk"`
	Ram           int    `json:"ram"`
	SSHKey        bool   `json:"ssh_key"`
	SSHPassword   bool   `json:"ssh_password"`
}

type ImgDoc struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type ServerImage struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	MinDisk   int32             `json:"min_disk"`
	MinRam    int32             `json:"min_ram"`
	OS        string            `json:"os"`
	OSVersion string            `json:"os_version"`
	Progress  int32             `json:"progress"`
	Size      int64             `json:"size"`
	Status    string            `json:"status"`
	Created   string            `json:"created"`
	UserName  string            `json:"username"`
	MetaData  map[string]string `json:"metadata"`
	Documents []*ImgDoc         `json:"documents"`
}

type PrivateImage struct {
	Abrak           string `json:"abrak"`
	AbrakID         string `json:"abrak_id"`
	Checksum        string `json:"checksum"`
	ContainerFormat string `json:"container_format"`
	CreatedAt       string `json:"created_at"`
	DiskFormat      string `json:"disk_format"`
	ID              string `json:"id"`
	ImageType       string `json:"image_type"`
	MinDisk         int    `json:"min_disk"`
	MinRam          int    `json:"min_ram"`
	Name            string `json:"name"`
	RealSize        int64  `json:"real_size"`
	Size            int64  `json:"size"`
	Status          string `json:"status"`
}

type ImageListItem struct {
	Name   string            `json:"name"`
	Images []ImageDistroItem `json:"images"`
}

type ImageListResponse struct {
	Data []ImageListItem `json:"data"`
}

type ImageClient struct {
	requester *Requester
}

func NewImageClient(r *Requester) *ImageClient {
	return &ImageClient{
		requester: r,
	}
}

func (i *ImageClient) ListImages(ctx context.Context, region, imgType string) (*ImageListResponse, error) {
	uri := fmt.Sprintf("%s/%s/images?type=%s", basePath, region, imgType)
	resp, err := i.requester.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}

	var ret ImageListResponse
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil

}

func (i *ImageClient) ListPrivateImages(ctx context.Context, region string) ([]PrivateImage, error) {
	uri := fmt.Sprintf("%s/%s/images?type=private", basePath, region)
	resp, err := i.requester.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret DataResponse[[]PrivateImage]
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	return ret.Data, nil
}

func (i *ImageClient) GetPrivateImageByID(ctx context.Context, region, id string) (*PrivateImage, error) {
	images, err := i.ListPrivateImages(ctx, region)
	if err != nil {
		return nil, err
	}
	for _, x := range images {
		if x.ID == id {
			return &x, nil
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "private image not found",
	}
}
