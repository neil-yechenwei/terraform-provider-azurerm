package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppPool_basic(t *testing.T) {
	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
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

func TestAccAzureRMNetAppPool_complete(t *testing.T) {
	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_level", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
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

func TestAccAzureRMNetAppPool_update(t *testing.T) {
	resourceName := "azurerm_netapp_pool.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppPool_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMNetAppPool_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_level", "Premium"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
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

func testCheckAzureRMNetAppPoolExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Pool not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).netapp.PoolClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: NetApp Pool %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on netapp.PoolClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetAppPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).netapp.PoolClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_netapp_pool" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on netapp.PoolClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNetAppPool_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestnetappaccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestnetapppool-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMNetAppPool_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestnetappaccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestnetapppool-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"

  tags = {
    env = "test"
  }
}
`, rInt, location, rInt, rInt)
}
