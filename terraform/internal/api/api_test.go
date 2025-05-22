package api

import (
	"context"
	"strings"
	"testing"
)

func TestCreatePrivateNetwork(t *testing.T) {
	c := NewSubnetClient(getRequester())

	apiReq := Subnet{
		Name:          "terraform_private",
		EnableDHCP:    true,
		EnableGateway: true,
		SubnetGateway: "10.255.255.21",
		CIDR:          "10.255.255.0/24",
		Description:   "",
	}
	dhcpRange := []string{"10.255.255.20", "10.255.255.150"}
	apiReq.DHCPRange = strings.Join(dhcpRange, ",")
	dnsServers := ""
	d := []string{"8.8.8.8", "1.1.1.1"}
	for idx, x := range d {
		dnsServers += x
		if idx < len(d)-1 {
			dnsServers += "\n"
		}
	}
	apiReq.DNSServers = dnsServers

	t.Log(apiReq)

	resp, err := c.CreatePrivateNetwork(context.Background(), "ir-thr-fr1", &apiReq)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.NetworkID)
}

func TestCreateSecurityGroupRule(t *testing.T) {
	c := NewClient("apikey 82e8ceae-b333-5fc9-8de8-53961f71fa65")
	req := RuleRequest{
		Direction: "egress",
		GroupID:   "86526aae-c2c1-422a-b506-dd9c310736d3",
		Protocol:  "udp",
		PortEnd:   "20000",
		PortStart: "18000",
		//IP:          []string{"any"},
		Description: "test",
	}

	err := c.Firewall.CreateRule(context.Background(), "ir-thr-fr1", "86526aae-c2c1-422a-b506-dd9c310736d3", &req)
	if err != nil {
		t.Fatal(err)
	}
	g, err := c.Firewall.GetSecurityGroupByID(context.Background(), "ir-thr-fr1", "86526aae-c2c1-422a-b506-dd9c310736d3")
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range g.Rules {
		t.Log(*x)
	}

}

func TestDeleteInstanceWithSnapshot(t *testing.T) {
	c := NewClient("apikey 82e8ceae-b333-5fc9-8de8-53961f71fa65")
	err := c.Instance.DeleteInstance(context.Background(), "ir-thr-fr1", "d971a7a2-f4e0-4140-8bfc-89258d191360")
	if err != nil {
		t.Fatal(err)
	}
}
