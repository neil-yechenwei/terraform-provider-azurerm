package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDataBoxCredential_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databox_credential", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataBoxCredential_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
		},
	})
}

func testAccDataSourceDataBoxCredential_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataBoxJob_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_databox_credential" "test" {
  databox_job_name    = azurerm_databox_job.test.name
  resource_group_name = azurerm_databox_job.test.resource_group_name
}
`, config)
}
