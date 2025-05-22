package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type ServerAddress struct {
	MAC      string `json:"mac_addr"`
	Version  string `json:"version"`
	Addr     string `json:"addr"`
	Type     string `json:"type"`
	IsPublic bool   `json:"is_public"`
}

type PortSecGroupData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type FullIP struct {
	FloatIP             interface{}        `json:"float_ip"`
	IP                  string             `json:"ip"`
	MacAddress          string             `json:"mac_address"`
	PortID              string             `json:"port_id"`
	PortSecurityEnabled bool               `json:"port_security_enabled"`
	PTR                 interface{}        `json:"ptr"`
	Public              bool               `json:"public"`
	SubnetID            string             `json:"subnet_id"`
	SubnetName          string             `json:"subnet_name"`
	Version             string             `json:"version"`
	SecurityGroups      []PortSecGroupData `json:"security_groups"`
}

type PublicIP struct {
	SubnetID  string `json:"subnet_id"`
	IPAddress string `json:"ip_address"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type HostRoute struct {
	DestinationCIDR string `json:"destination"`
	NextHop         string `json:"nexthop"`
}

type NetworkServer struct {
	Addresses      map[string][]*ServerAddress `json:"addresses"`
	ID             string                      `json:"id"`
	Name           string                      `json:"name"`
	IPs            []*FullIP                   `json:"ips"`
	PublicIPs      []*PublicIP                 `json:"public_ip"`
	SecurityGroups []string                    `json:"security_groups"`
}

type Subnet struct {
	Name          string           `json:"name"`
	DHCPRange     string           `json:"dhcp"`        //192.168.1.2,192.168.1.140
	DNSServers    string           `json:"dns_servers"` //8.8.8.8\n1.1.1.1
	EnableDHCP    bool             `json:"enable_dhcp"`
	EnableGateway bool             `json:"enable_gateway"`
	NetworkID     string           `json:"network_id"`
	SubnetGateway string           `json:"subnet_gateway"`
	SubnetID      string           `json:"subnet_id"`
	CIDR          string           `json:"subnet_ip"` //192.168.1.1/24
	Description   string           `json:"description"`
	Servers       []*NetworkServer `json:"servers"`
}

type SubnetDetails struct {
	ID string `json:"id"`

	// UUID of the parent network.
	NetworkID string `json:"network_id"`

	// Human-readable name for the subnet. Might not be unique.
	Name string `json:"name"`

	// Description for the subnet.
	Description string `json:"description"`

	// IP version, either `4' or `6'.
	IPVersion string `json:"ip_version"`

	// CIDR representing IP range for this subnet, based on IP version.
	CIDR string `json:"cidr"`

	// Default gateway used by devices in this subnet.
	GatewayIP *string `json:"gateway_ip"`

	// DNS name servers used by hosts in this subnet.
	DNSNameservers []string `json:"dns_nameservers"`

	// Sub-ranges of CIDR available for dynamic allocation to ports.
	// See AllocationPool.
	AllocationPools []AllocationPool `json:"allocation_pools"`

	// Routes that should be used by devices with IPs from this subnet
	// (not including local subnet route).
	HostRoutes []HostRoute `json:"host_routes"`

	// Specifies whether DHCP is enabled for this subnet or not.
	EnableDHCP bool `json:"enable_dhcp"`

	// TenantID is the project owner of the subnet.
	TenantID string `json:"tenant_id"`

	// ProjectID is the project owner of the subnet.
	ProjectID string `json:"project_id"`

	// The IPv6 address modes specifies mechanisms for assigning IPv6 IP addresses.
	IPv6AddressMode string `json:"ipv6_address_mode"`

	// The IPv6 router advertisement specifies whether the networking service
	// should transmit ICMPv6 packets.
	IPv6RAMode string `json:"ipv6_ra_mode"`

	// SubnetPoolID is the id of the subnet pool associated with the subnet.
	SubnetPoolID   string           `json:"subnetpool_id"`
	ServiceType    []string         `json:"service_types"`
	RevisionNumber int              `json:"revision_number"`
	Tags           []string         `json:"tags"`
	Servers        []*NetworkServer `json:"servers"`
}

type Network struct {
	ID                    string           `json:"id"`
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	AdminStateUp          bool             `json:"admin_state_up"`
	Shared                bool             `json:"shared"`
	Status                string           `json:"status"`
	Subnets               []*SubnetDetails `json:"subnets"`
	TenantID              string           `json:"tenant_id"`
	DHCPIP                string           `json:"dhcp_ip"`
	UpdatedAt             string           `json:"updated_at"`
	CreatedAt             string           `json:"created_at"`
	IPV4AddressScope      string           `json:"ipv4_address_scope"`
	IPV6AddressScope      string           `json:"ipv6_address_scope"`
	QOSPolicyID           string           `json:"qos_policy_id"`
	ReviosionNumber       *int             `json:"revision_number"`
	RouterExternal        *bool            `json:"router:external"`
	MTU                   int              `json:"mtu"`
	PortSecurityEnabled   bool             `json:"port_security_enabled"`
	AvailabilityZoneHints []string         `json:"availability_zone_hints"`
	AvailabilityZones     []string         `json:"availability_zones"`
	Tags                  []string         `json:"tags"`
}

type AttachServerToNetworkRequest struct {
	ServerID           string `json:"server_id"`
	IP                 string `json:"ip"`
	SubnetID           string `json:"subnet_id"`
	EnablePortSecurity bool   `json:"enablePortSecurity"`
}

type AttachedPort struct {
	AdminStateUp    bool   `json:"admin_state_up"`
	IsRegionNetwork bool   `json:"is_region_network"`
	Status          string `json:"status"`
	MacAddr         string `json:"mac_addr"`
	IPAddress       string `json:"ip_address"`
	ID              string `json:"id"`
	NetworkID       string `json:"network_id"`
	DeviceID        string `json:"device_id"`
	SubnetID        string `json:"subnet_id"`
}

type SubnetClient struct {
	requester *Requester
}

func NewSubnetClient(r *Requester) *SubnetClient {
	return &SubnetClient{
		requester: r,
	}
}

func (s *SubnetClient) CreatePrivateNetwork(ctx context.Context, region string, req *Subnet) (*SubnetDetails, error) {
	type subnetResponse struct {
		Data    *SubnetDetails `json:"data"`
		Message string         `json:"message"`
	}
	url := fmt.Sprintf("%s/%s/subnets", basePath, region)

	data, err := s.requester.DoRequest(ctx, "POST", url, req)
	if err != nil {
		return nil, err
	}
	var resp subnetResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *SubnetClient) UpdatePrivateNetwork(ctx context.Context, region string, req *Subnet) error {
	url := fmt.Sprintf("%s/%s/subnets", basePath, region)
	_, err := s.requester.DoRequest(ctx, "PATCH", url, req)
	return err
}

func (s *SubnetClient) DeletePrivateNetwork(ctx context.Context, region, id string) error {
	url := fmt.Sprintf("%s/%s/subnets/%s", basePath, region, id)
	_, err := s.requester.DoRequest(ctx, "DELETE", url, nil)
	return err
}

func (s *SubnetClient) GetPrivateNetwork(ctx context.Context, region, id string) (*SubnetDetails, error) {
	type subnetResponse struct {
		Data *SubnetDetails `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/subnets/%s", basePath, region, id)

	data, err := s.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp subnetResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *SubnetClient) GetAllNetworks(ctx context.Context, region string) (map[string]*Network, error) {
	type networksResponse struct {
		Data []*Network `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/networks", basePath, region)

	data, err := s.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp networksResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Network)
	for _, x := range resp.Data {
		m[x.ID] = x
	}
	return m, nil
}

func (s *SubnetClient) GetNetworkSubnet(ctx context.Context, region, netID string) (*SubnetDetails, error) {
	nets, err := s.GetAllNetworks(ctx, region)
	if err != nil {
		return nil, err
	}
	n, ok := nets[netID]
	if !ok {
		return nil, &ResponseError{
			Code:    404,
			Message: "network not found",
		}
	}
	if len(n.Subnets) == 0 {
		return nil, &ResponseError{
			Code:    404,
			Message: "network has no subnet",
		}
	}
	return n.Subnets[0], nil
}

func (s *SubnetClient) GetAllNetworksByName(ctx context.Context, region string) (map[string]*Network, error) {
	type networksResponse struct {
		Data []*Network `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/networks", basePath, region)

	data, err := s.requester.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var resp networksResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*Network)
	for _, x := range resp.Data {
		m[x.Name] = x
	}
	return m, nil
}

func (s *SubnetClient) DetachServerFromNetwork(ctx context.Context, region, portID, serverId string) error {
	type serverIdBody struct {
		ServerId string `json:"server_id"`
	}
	url := fmt.Sprintf("%s/%s/networks/%s/detach", basePath, region, portID)
	req := serverIdBody{
		ServerId: serverId,
	}
	_, err := s.requester.DoRequest(ctx, "PATCH", url, &req)
	return err
}

func (s *SubnetClient) AttachServerToNetwork(ctx context.Context, region, networkID string, req *AttachServerToNetworkRequest) (*AttachedPort, error) {
	type attachResponse struct {
		Data    *AttachedPort `json:"data"`
		Message string        `json:"message"`
	}
	url := fmt.Sprintf("%s/%s/networks/%s/attach", basePath, region, networkID)
	data, err := s.requester.DoRequest(ctx, "PATCH", url, req)
	if err != nil {
		return nil, err
	}
	var resp attachResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (s *SubnetClient) DisablePortSecurity(ctx context.Context, region, networkID, portID string) error {
	type disableReq struct {
		NetworkID string `json:"network_id"`
	}
	url := fmt.Sprintf("%s/%s/ports/%s/disablePortSecurity", basePath, region, portID)
	_, err := s.requester.DoRequest(ctx, "PATCH", url, &disableReq{networkID})
	return err
}

func (s *SubnetClient) EnablePortSecurity(ctx context.Context, region, networkID, portID string) error {
	type enableReq struct {
		NetworkID string `json:"network_id"`
	}
	url := fmt.Sprintf("%s/%s/ports/%s/enablePortSecurity", basePath, region, portID)
	_, err := s.requester.DoRequest(ctx, "PATCH", url, &enableReq{networkID})
	return err
}
