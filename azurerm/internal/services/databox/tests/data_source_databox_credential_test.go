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
				Config: testAccDataSourceDataBoxCredential_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "databox_job_name"),
				),
			},
		},
	})
}

func testAccDataSourceDataBoxCredential_basic() string {
	return fmt.Sprintf(`
locals {
  databox_job_name    = "TJ-636646322037905056"
  resource_group_name = "bvttoolrg6"
}

data "azurerm_databox_credential" "test" {
  databox_job_name    = local.databox_job_name
  resource_group_name = local.resource_group_name
}
`)
}
