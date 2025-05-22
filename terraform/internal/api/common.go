package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"time"
)

const (
	basePath       = "https://napi.arvancloud.ir/ecc/v1/regions"
	DefaultTimeout = 1 * time.Minute
	basePathV2     = "https://napi.arvancloud.ir/ecc/v2/ssc"
	bpV2           = "https://napi.arvancloud.ir/ecc/v2"
)

var (
	ErrTimeout = errors.New("operation timed out")
)

type ResponseError struct {
	Code    int
	Message string
	URL     string
	Errors  []string
}

type DataResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("an error ocurred calling endpoint: %s. status code: %d. message: %s", r.URL, r.Code, r.Message)
}

type Requester struct {
	client *http.Client
	apiKey string
}

func NewRequester(c *http.Client, apiKey string) *Requester {
	return &Requester{
		client: c,
		apiKey: apiKey,
	}
}

type errorResponse struct {
	Code    int      `json:"code,omitempty"`
	Status  int      `json:"status,omitempty"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func (r *Requester) DoRequest(ctx context.Context, method, uri string, input interface{}) ([]byte, error) {

	var body io.Reader
	if input != nil {
		js, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(js)
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", r.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, &ResponseError{
			Code:    resp.StatusCode,
			URL:     uri,
			Message: "internal server error",
		}
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		var errResp errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err == nil {
			return nil, &ResponseError{
				Code:    resp.StatusCode,
				URL:     uri,
				Message: errResp.Message,
				Errors:  errResp.Errors,
			}
		}
		return nil, &ResponseError{
			Code:    resp.StatusCode,
			URL:     uri,
			Message: resp.Status,
		}

	}

	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//tflog.Info(ctx, "RESPONSE FROM SERVER", map[string]interface{}{"data": string(ret)})
	return ret, nil

}

type Client struct {
	Img             *ImageClient
	Pln             *PlanClient
	Instance        *InstanceClient
	Volume          *VolumeClient
	Subnet          *SubnetClient
	Firewall        *SecurityGroupClient
	FIPClient       *FloatingIPClient
	SnapshotClient  *SnapshotClient
	SSHClient       *SSHKeyClient
	VolumeV2        *VolumeV2Client
	BackupV2        *BackupV2Client
	ServerGroup     *ServerGroupClient
	DedicatedServer *DedicatedServerClient
	FirewallV2      *FirewallV2Client
}

func NewClient(apiKey string) *Client {
	c := &http.Client{
		Timeout: DefaultTimeout,
	}
	r := NewRequester(c, apiKey)
	imgC := NewImageClient(r)
	plnC := NewPlanClient(r)
	instanceC := NewInstanceClient(r)
	volC := NewVolumeClient(r)
	subnetC := NewSubnetClient(r)
	fwC := NewSecurityGroupClient(r)
	fwv2C := NewFirewallV2Clien(r)
	fipC := NewFloatingIPClient(r)
	snpC := NewSnapshotClient(r)
	sshC := NewSSHKeyClient(r)
	vv2C := NewVolumeV2Client(r)
	bV2 := NewBackupV2Client(r)
	serverGroupC := NewServerGroupClient(r)
	dsClient := NewDedicatedServerClient(r)
	ret := &Client{
		Img:             imgC,
		Pln:             plnC,
		Instance:        instanceC,
		Volume:          volC,
		Subnet:          subnetC,
		Firewall:        fwC,
		FIPClient:       fipC,
		SnapshotClient:  snpC,
		SSHClient:       sshC,
		VolumeV2:        vv2C,
		BackupV2:        bV2,
		FirewallV2:      fwv2C,
		ServerGroup:     serverGroupC,
		DedicatedServer: dsClient,
	}
	return ret
}

func (c *Client) WaitForCondition(ctx context.Context, timeout time.Duration, condition func() (bool, error)) error {
	tick := time.NewTicker(5 * time.Second)
	timeoutTimer := time.NewTimer(timeout)
	defer func() {
		tick.Stop()
		timeoutTimer.Stop()
	}()
	for {
		select {
		case <-tick.C:
			ok, err := condition()
			if err != nil {
				tflog.Trace(ctx, "error checking condition", map[string]interface{}{"err": err})
				continue
			}
			if ok {
				return nil
			}
		case <-timeoutTimer.C:
			return ErrTimeout
		}
	}
}

type NetworkAttachment struct {
	NetworkID           string
	SubnetID            string
	IP                  string
	PortID              string
	IsPublic            bool
	PortSecurityEnabled bool
}

type FloatingIPInfo struct {
	ID               string
	PrivateNetworkID string
}

func (c *Client) GetServerFloatingIPInfo(ctx context.Context, region, serverID string) (*FloatingIPInfo, error) {
	fIPList, err := c.FIPClient.GetAllFloatingIPs(ctx, region)
	if err != nil {
		return nil, err
	}
	for _, f := range fIPList {
		if f.Server != nil && f.Server.ID == serverID {
			if f.Server.Addresses != nil {
				for name, addr := range f.Server.Addresses {
					for _, x := range addr {
						if x.Addr == f.FloatingIPAddress {
							netList, err := c.Subnet.GetAllNetworksByName(ctx, region)
							if err != nil {
								return nil, err
							}
							if n, ok := netList[name]; ok {
								return &FloatingIPInfo{
									ID:               f.ID,
									PrivateNetworkID: n.ID,
								}, nil
							}
						}
					}

				}
			}
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "server has no floating ip",
	}
}

func (c *Client) GetNetworkAttachments(ctx context.Context, detail *ServerDetail, region, serverID string) (map[string]NetworkAttachment, error) {
	if detail.Addresses != nil {
		networks, err := c.Subnet.GetAllNetworksByName(ctx, region)
		if err != nil {
			return nil, err
		}
		ret := make(map[string]NetworkAttachment)

		for name, _ := range detail.Addresses {

			if n, ok := networks[name]; ok {

				if len(n.Subnets) > 0 {
					attachment := NetworkAttachment{
						NetworkID: n.Subnets[0].NetworkID,
						SubnetID:  "",
						IP:        "",
						PortID:    "",
					}
				outer:
					for _, s := range n.Subnets[0].Servers {
						for _, i := range s.IPs {
							if i.SubnetID == n.Subnets[0].ID && s.ID == serverID {
								attachment.SubnetID = i.SubnetID
								attachment.PortID = i.PortID
								attachment.IP = i.IP
								attachment.IsPublic = i.Public
								attachment.PortSecurityEnabled = i.PortSecurityEnabled
								break outer
							}

						}

					}
					ret[n.ID] = attachment

				}

			}

		}
		return ret, nil
	}
	return nil, errors.New("unexpexted result")

}

func (c *Client) FillNetworkData(ctx context.Context, netIds []string, region, serverID string) (map[string]NetworkAttachment, error) {
	networks, err := c.Subnet.GetAllNetworks(ctx, region)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]NetworkAttachment)

	for idx := 0; idx < len(netIds); idx++ {
		tflog.Info(ctx, netIds[idx])
		if n, ok := networks[netIds[idx]]; ok && len(n.Subnets) > 0 {

			attachment := NetworkAttachment{
				NetworkID: n.Subnets[0].NetworkID,
				SubnetID:  n.Subnets[0].ID,
				IP:        "",
				PortID:    "",
			}

			if len(n.Subnets) > 0 {
			outer:
				for _, s := range n.Subnets[0].Servers {

					for _, i := range s.IPs {
						if i.SubnetID == n.Subnets[0].ID && s.ID == serverID {
							tflog.Warn(ctx, "SERVERIDS", map[string]interface{}{"IDS": []string{s.ID, serverID, i.PortID, i.IP}})
							//tflog.Info(ctx, "SETTING_PORT_IP", map[string]interface{}{"NETWORK_NMAE": n.Name})
							attachment.PortID = i.PortID
							attachment.IP = i.IP
							attachment.IsPublic = i.Public
							attachment.PortSecurityEnabled = i.PortSecurityEnabled
							break outer
						}

					}
				}
				ret[n.ID] = attachment
			}

		}

	}

	return ret, nil
}
