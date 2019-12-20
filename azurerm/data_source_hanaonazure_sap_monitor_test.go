package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMHanaOnAzureSapMonitor_basic(t *testing.T) {
	dataSourceName := "data.azurerm_hanaonazure_sap_monitor.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHanaOnAzureSapMonitor_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceHanaOnAzureSapMonitor_basic(rInt int, location string) string {
	config := testAccAzureRMHanaOnAzureSapMonitor_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_hanaonazure_sap_monitor" "test" {
  name                = "${azurerm_hanaonazure_sap_monitor.test.name}"
  resource_group_name = "${azurerm_hanaonazure_sap_monitor.test.resource_group_name}"
}
`, config)
}
