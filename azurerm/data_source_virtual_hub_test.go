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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMVirtualHub_basic(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub" "test" {
  name                = "${azurerm_virtual_hub.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}
