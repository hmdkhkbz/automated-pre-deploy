package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccVolumeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ArvanProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: volumeResourceConfig("ir-thr-fr1", "acc", 9),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("arvan_volume.test", "id"),
				),
			},
		},
	})
}

func volumeResourceConfig(region string, name string, size int) string {
	return fmt.Sprintf(`
variable "API_KEY" {
	type = string
}
provider "arvan" {
  api_key = var.API_KEY
}
resource "arvan_volume" "test" {
	region = "%s"
	name = "%s"
	description = "for acceptance test"
	size = %d
}
`, region, name, size)
}
