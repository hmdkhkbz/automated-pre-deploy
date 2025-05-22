package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type Rule struct {
	ID          string `json:"id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	EtherType   string `json:"ether_type,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	IP          string `json:"ip"`
	PortStart   int32  `json:"port_start"`
	PortEnd     int32  `json:"port_end"`
	Protocol    string `json:"protocol"`
}

type RuleRequest struct {
	ID          string   `json:"id,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
	Description string   `json:"description"`
	Direction   string   `json:"direction"`
	EtherType   string   `json:"ether_type,omitempty"`
	GroupID     string   `json:"group_id,omitempty"`
	IP          []string `json:"ips"`
	PortStart   string   `json:"port_from"`
	PortEnd     string   `json:"port_to"`
	Protocol    string   `json:"protocol"`
}

type SecurityGroup struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	ReadOnly    bool     `json:"readonly"`
	Default     bool     `json:"default"`
	RealName    string   `json:"real_name"`
	Rules       []*Rule  `json:"rules"`
	IPAddresses []string `json:"ip_addresses,omitempty"`
}

type SecurityGroupClient struct {
	requester *Requester
}

func NewSecurityGroupClient(r *Requester) *SecurityGroupClient {
	return &SecurityGroupClient{
		requester: r,
	}
}

func (s *SecurityGroupClient) CreateSecurityGroup(ctx context.Context, region, name, description string) (*SecurityGroup, error) {
	type createReq struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type response struct {
		Data    *SecurityGroup `json:"data"`
		Message string         `json:"message"`
	}
	url := fmt.Sprintf("%s/%s/securities", basePath, region)

	data, err := s.requester.DoRequest(ctx, "POST", url, &createReq{Name: name, Description: description})
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

func (s *SecurityGroupClient) DeleteSecurityGroup(ctx context.Context, region, groupID string) error {
	url := fmt.Sprintf("%s/%s/securities/%s", basePath, region, groupID)
	_, err := s.requester.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (s *SecurityGroupClient) GetSecurityGroupByID(ctx context.Context, region, groupID string) (*SecurityGroup, error) {
	type response struct {
		Data *SecurityGroup
	}
	url := fmt.Sprintf("%s/%s/securities/security-rules/%s", basePath, region, groupID)

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

func (s *SecurityGroupClient) CreateRule(ctx context.Context, region, groupID string, req *RuleRequest) error {
	url := fmt.Sprintf("%s/%s/securities/security-rules/%s", basePath, region, groupID)

	_, err := s.requester.DoRequest(ctx, "POST", url, req)
	return err
}

func (s *SecurityGroupClient) DeleteRule(ctx context.Context, region, ruleID string) error {
	url := fmt.Sprintf("%s/%s/securities/security-rules/%s", basePath, region, ruleID)
	_, err := s.requester.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (s *SecurityGroupClient) AddServerToGroup(ctx context.Context, region, serverID, groupID string) error {
	type addReq struct {
		GroupID string `json:"security_group_id"`
	}
	url := fmt.Sprintf("%s/%s/servers/%s/add-security-group", basePath, region, serverID)

	_, err := s.requester.DoRequest(ctx, "POST", url, &addReq{GroupID: groupID})
	return err
}

func (s *SecurityGroupClient) RemoveServerFromGroup(ctx context.Context, region, serverID, groupID string) error {
	type addReq struct {
		GroupID string `json:"security_group_id"`
	}
	url := fmt.Sprintf("%s/%s/servers/%s/remove-security-group", basePath, region, serverID)

	_, err := s.requester.DoRequest(ctx, "POST", url, &addReq{GroupID: groupID})
	return err
}

func (s *SecurityGroupClient) GetAllSecurityGroups(ctx context.Context, region string) ([]*SecurityGroup, error) {
	type dataResponse struct {
		Data []*SecurityGroup `json:"data"`
	}

	url := fmt.Sprintf("%s/%s/securities", basePath, region)
	data, err := s.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp dataResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
