package api

import (
	"context"
	"encoding/json"
	"fmt"
)


type ServerGroupClient struct{
	requester *Requester
}

func NewServerGroupClient(r *Requester) *ServerGroupClient {
	return &ServerGroupClient{
		requester: r,
	}
}

type ServerGroupDetail struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Policies []string `json:"policies"`
	Members  []string `json:"members"`
}

func (i *ServerGroupClient) ListServerGroups(ctx context.Context, region string) ([]ServerGroupDetail, error) {
	type serverGroupListResponse struct {
		Data []ServerGroupDetail `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/server-groups", basePath, region)

	data, err := i.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp serverGroupListResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}