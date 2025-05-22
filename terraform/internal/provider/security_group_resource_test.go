package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
)

func TestAccSecurityGroupResource(t *testing.T) {
	sg := secG{
		region:      "ir-thr-fr1",
		name:        "test_acc",
		description: "for acceptance test",
		rules: []sgRule{
			{
				direction: "ingress",
				from:      "15000",
				to:        "20000",
				protocol:  "tcp",
			},
			{
				direction: "ingress",
				from:      "21000",
				to:        "22000",
				protocol:  "udp",
			},
			{
				direction: "ingress",
				protocol:  "udp",
				ip:        "192.168.0.240",
			},
		},
	}
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ArvanProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: sg.tfConfig(),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("arvan_security_group.test", "id")),
			},
		},
	})
}

type secG struct {
	region      string
	name        string
	description string
	rules       []sgRule
}

func (s *secG) tfConfig() string {
	rules := "["
	for _, r := range s.rules {
		rules += r.tfConfig() + ","
	}
	rules = strings.TrimSuffix(rules, ",")
	rules += "]"
	ret := fmt.Sprintf(`
variable "API_KEY" {
	type = string
}
provider "arvan" {
  api_key = var.API_KEY
}
resource "arvan_security_group" "test" {
	region = "%s"
	name = "%s"
	description = "%s"
	rules = %s
}
`, s.region, s.name, s.description, rules)
	return ret
}

type sgRule struct {
	ip        string
	from      string
	to        string
	direction string
	protocol  string
}

func (r *sgRule) tfConfig() string {
	s := fmt.Sprintf(`
{
	protocol = "%s"
	direction = "%s"`, r.protocol, r.direction)

	if r.from != "" {
		s = fmt.Sprintf(`%s
	port_from = "%s"`, s, r.from)
	}

	if r.to != "" {
		s = fmt.Sprintf(`%s
	port_to = "%s"`, s, r.to)
	}

	if r.ip != "" {
		s = fmt.Sprintf(`%s
	ip = "%s"`, s, r.ip)
	}
	s += "\n}"
	return s
}
