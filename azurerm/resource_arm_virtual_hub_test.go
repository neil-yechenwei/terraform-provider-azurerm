package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHub_basic(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
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
func TestAccAzureRMVirtualHub_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
				),
			},
			{
				Config:      testAccAzureRMVirtualHub_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_virtual_hub"),
			},
		},
	})
}

func testCheckAzureRMVirtualHubDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.VirtualHubClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Hub still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMVirtualHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualHubName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Virtual Hub: %s", virtualHubName)
		}

		client := testAccProvider.Meta().(*ArmClient).network.VirtualHubClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualHubName)
		if err != nil {
			return fmt.Errorf("Bad: Get on virtualHubClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Hub %q (resource group: %q) does not exist", virtualHubName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMVirtualHub_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
	name                = "acctestvhub%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location            = "${azurerm_resource_group.test.location}"
	address_prefix      = "10.168.0.0/24"
	virtual_wan_id      = "${azurerm_virtual_wan.test.id}"
	tags = {
		"env" = "test"
	}
  }
`, rInt, location, rInt, rInt)
}

func testAccAzureRMVirtualHub_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualHub_basic(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub" "import" {
  name                = "${azurerm_virtual_hub.test.name}"
  resource_group_name = "${azurerm_virtual_hub.test.resource_group_name}"
}
`, template)
}
