package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type FloatIPResponse struct {
	ID                string        `json:"id"`
	Status            string        `json:"status"`
	CreatedAt         string        `json:"created_at"`
	Description       string        `json:"description"`
	FixedIPAddress    string        `json:"fixed_ip_address"`
	FloatingIPAddress string        `json:"floating_ip_address"`
	FloatingNetworkID string        `json:"floating_network_id"`
	PortID            string        `json:"port_id"`
	RevisionNumber    string        `json:"revision_number"`
	RouterID          string        `json:"router_id"`
	Tags              []string      `json:"tags"`
	UpdatedAt         string        `json:"updated_at"`
	Server            *ServerDetail `json:"server"`
}

type FloatIPAttachRequest struct {
	ServerID string `json:"server_id"`
	PortID   string `json:"port_id"`
	SubnetID string `json:"subnet_id"`
}

type ServerIPInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	CreationDate string    `json:"creation_date"`
	Status       string    `json:"status"`
	HasPublicIP  bool      `json:"has_public_ip"`
	IPData       []*IPInfo `json:"ip_data"`
}

type IPInfo struct {
	SubnetID       string `json:"subnet_id"`
	Type           string `json:"type"`
	PortID         string `json:"port_id"`
	Address        string `json:"address"`
	GatewayAddress string `json:"gateway_address"`
}

type AttachReq struct {
	ServerID string `json:"server_id"`
	SubnetID string `json:"subnet_id"`
	PortID   string `json:"port_id"`
}

type FloatingIPClient struct {
	requester *Requester
}

func NewFloatingIPClient(r *Requester) *FloatingIPClient {
	return &FloatingIPClient{
		requester: r,
	}
}

func (f *FloatingIPClient) CreateFloatingIP(ctx context.Context, region, description string) (*FloatIPResponse, error) {
	type createReq struct {
		Description string `json:"description"`
	}
	type dataResponse struct {
		Data FloatIPResponse `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/float-ips", basePath, region)
	data, err := f.requester.DoRequest(ctx, "POST", url, &createReq{description})
	if err != nil {
		return nil, err
	}
	var resp dataResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (f *FloatingIPClient) GetAllFloatingIPs(ctx context.Context, region string) ([]*FloatIPResponse, error) {
	type response struct {
		Data []*FloatIPResponse `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/float-ips", basePath, region)
	data, err := f.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (f *FloatingIPClient) GetFloatingIP(ctx context.Context, region, floatingID string) (*FloatIPResponse, error) {
	data, err := f.GetAllFloatingIPs(ctx, region)
	if err != nil {
		return nil, err
	}
	for _, x := range data {
		if x.ID == floatingID {
			return x, nil
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "floating ip not found",
	}
}

func (f *FloatingIPClient) DeleteFloatingIP(ctx context.Context, region, floatingID string) error {
	url := fmt.Sprintf("%s/%s/float-ips/%s", basePath, region, floatingID)
	_, err := f.requester.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (f *FloatingIPClient) GetServerIPInfo(ctx context.Context, region string) (map[string]*ServerIPInfo, error) {
	type dataResponse struct {
		Data []*ServerIPInfo `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/float-ips/ips", basePath, region)
	data, err := f.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp dataResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*ServerIPInfo)
	for _, x := range resp.Data {
		ret[x.ID] = x
	}
	return ret, nil
}

func (f *FloatingIPClient) AttachFloatingIP(ctx context.Context, region, floatingID string, req *AttachReq) error {
	url := fmt.Sprintf("%s/%s/float-ips/%s/attach", basePath, region, floatingID)

	_, err := f.requester.DoRequest(ctx, "PATCH", url, req)
	return err
}

func (f *FloatingIPClient) DetachFloatingIP(ctx context.Context, region, portID string) error {
	type detachReq struct {
		PortID string `json:"port_id"`
	}
	url := fmt.Sprintf("%s/%s/float-ips/detach", basePath, region)
	req := &detachReq{
		PortID: portID,
	}
	_, err := f.requester.DoRequest(ctx, "PATCH", url, req)
	return err
}
