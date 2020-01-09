package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMManagedApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceManagedApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceManagedApplication_basic(data acceptance.TestData) string {
	config := testAccAzureRMManagedApplication_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_managed_application" "test" {
  resource_group_name = "${azurerm_managed_application.test.resource_group_name}"
  name                = "${azurerm_managed_application.test.name}"
}
`, config)
}
