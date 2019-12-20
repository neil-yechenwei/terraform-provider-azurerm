package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMRecoveryFabric_basic(t *testing.T) {
	resourceName := "azurerm_recovery_services_fabric.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryFabricDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryFabric_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryFabricExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMRecoveryFabric_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_fabric" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  name                = "acctest-fabric-%d"
  location            = "${azurerm_resource_group.test.location}"
}
`, rInt, location, rInt, rInt)
}

func testCheckAzureRMRecoveryFabricExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		state, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroupName := state.Primary.Attributes["resource_group_name"]
		vaultName := state.Primary.Attributes["recovery_vault_name"]
		fabricName := state.Primary.Attributes["name"]

		// Ensure fabric exists in API
		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.FabricClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, fabricName)
		if err != nil {
			return fmt.Errorf("Bad: Get on fabricClient: %+v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: fabric: %q does not exist", fabricName)
		}

		return nil
	}
}

func testCheckAzureRMRecoveryFabricDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_services_fabric" {
			continue
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		resourceGroup := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.FabricClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Recovery Services Fabric still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}
