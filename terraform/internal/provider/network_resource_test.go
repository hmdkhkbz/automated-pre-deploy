package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
)

func TestAccNetworkResource(t *testing.T) {
	nw := accNetwork{
		region:      "ir-thr-fr1",
		name:        "test",
		description: "for acc test",
		dhcpRange: ipRange{
			start: "10.255.255.19",
			end:   "10.255.255.150",
		},
		cidr:          "10.255.255.0/24",
		dnsServers:    []string{"8.8.8.8", "1.1.1.1"},
		enableGateway: true,
		enableDHCP:    true,
		gatewayIP:     "10.255.255.1",
	}
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ArvanProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nw.tfConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("arvan_network.test", "id"),
				),
			},
		},
	})
}

type ipRange struct {
	start string
	end   string
}

type accNetwork struct {
	region        string
	name          string
	description   string
	dhcpRange     ipRange
	cidr          string
	dnsServers    []string
	enableDHCP    bool
	enableGateway bool
	gatewayIP     string
}

func (n *accNetwork) dnsServersToString() string {
	var ret string
	ret = "["
	for i := 0; i < len(n.dnsServers); i++ {
		ret += "\"" + n.dnsServers[i] + "\"" + ","

	}
	ret = strings.TrimSuffix(ret, ",")
	ret += "]"
	return ret
}

func (n *accNetwork) tfConfig() string {
	return fmt.Sprintf(`
variable "API_KEY" {
	type = string
}
provider "arvan" {
  api_key = var.API_KEY
}
resource "arvan_network" "test" {
	region = "%s"
	name = "%s"
	description = "%s"
	dhcp_range = {
		start = "%s"
		end = "%s"
	}
	dns_servers = %s
	enable_dhcp = %t
	enable_gateway = %t
	cidr = "%s"
	gateway_ip = "%s"
}
`, n.region, n.name, n.description, n.dhcpRange.start, n.dhcpRange.end, n.dnsServersToString(), n.enableDHCP, n.enableGateway, n.cidr, n.gatewayIP)
}
