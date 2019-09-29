package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMNetAppVolume_basic(t *testing.T) {
	dataSourceName := "data.azurerm_netapp_volume.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetAppVolume_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "creation_token"),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "usage_threshold"),
				),
			},
		},
	})
}

func testAccDataSourceNetAppVolume_basic(rInt int, location string) string {
	config := testAccAzureRMNetAppVolume_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume" "test" {
  resource_group_name = "${azurerm_netapp_volume.test.resource_group_name}"
  account_name        = "${azurerm_netapp_volume.test.account_name}"
  pool_name           = "${azurerm_netapp_volume.test.pool_name}"
  name                = "${azurerm_netapp_volume.test.name}"
}
`, config)
}
