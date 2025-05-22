package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccFloatingIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ArvanProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: floatingIPResourceConfig("ir-thr-fr1", "for test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("arvan_floating_ip.test", "id"),
				),
			},
		},
	})
}

func floatingIPResourceConfig(region string, description string) string {
	return fmt.Sprintf(`
variable "API_KEY" {
	type = string
}
provider "arvan" {
  api_key = var.API_KEY
}
resource "arvan_floating_ip" "test" {
	region = "%s"
	description = "%s"
}
`, region, description)
}
