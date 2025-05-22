package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type FirewallV2Client struct {
	requester *Requester
}

func NewFirewallV2Clien(r *Requester) *FirewallV2Client {
	return &FirewallV2Client{
		requester: r,
	}
}

type SecurityGroupV2ConnectedInstance struct {
	InstanceID     string   `json:"instance_id"`
	InstanceName   string   `json:"instance_name"`
	ConnectedPorts []string `json:"connected_ports"`
}

func (s *FirewallV2Client) GetFirewallConnectedInstances(ctx context.Context, region, groupID string) ([]SecurityGroupV2ConnectedInstance, error) {
	type response struct {
		Data []SecurityGroupV2ConnectedInstance `json:"data"`
	}
	url := fmt.Sprintf("%s/firewall/%s/%s/instance/list", bpV2, region, groupID)

	data, err := s.requester.DoRequest(ctx, "GET", url, nil)
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

func (s *FirewallV2Client) DetachInstancesFromFirewall(ctx context.Context, region, groupID string, instanceIDs []string) error {
	type req struct {
		InstanceIDs []string `json:"instance_ids"`
	}
	url := fmt.Sprintf("%s/firewall/%s/%s/detach-instance", bpV2, region, groupID)

	_, err := s.requester.DoRequest(ctx, "POST", url, req{InstanceIDs: instanceIDs})
	if err != nil {
		return err
	}

	return nil
}