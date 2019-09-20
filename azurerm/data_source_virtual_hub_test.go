package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMVirtualHub_basic(t *testing.T) {
	dataSourceName := "data.azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table.0.routes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table.0.routes.0.address_prefixes.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table.0.routes.0.address_prefixes.0", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table.0.routes.0.address_prefixes.1", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table.0.routes.0.next_hop_ip_address", "10.0.4.5"),
				),
			},
		},
	})
}

func testAccDataSourceVirtualHub_basic(rInt int, location string) string {
	config := testAccAzureRMVirtualHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub" "test" {
  resource_group_name = "${azurerm_virtual_hub.test.resource_group_name}"
  name           = "${azurerm_virtual_hub.test.name}"
}
`, config)
}
