package azurerm

import (
	"fmt"
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
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(resourceName, "route_table.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.0", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.1", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.next_hop_ip_address", "10.0.4.5"),
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

func TestAccAzureRMVirtualHub_complete(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.name", "testConnection"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_network_connections.0.remote_virtual_network_id"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.allow_hub_to_remote_vnet_transit", "false"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.allow_remote_vnet_to_use_hub_vnet_gateways", "false"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.enable_internet_security", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_table.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.0", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.1", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.next_hop_ip_address", "10.0.4.6"),
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

func TestAccAzureRMVirtualHub_update(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(resourceName, "route_table.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.0", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.1", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.next_hop_ip_address", "10.0.4.5"),
				),
			},
			{
				Config: testAccAzureRMVirtualHub_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.name", "testConnection"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_network_connections.0.remote_virtual_network_id"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.allow_hub_to_remote_vnet_transit", "false"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.allow_remote_vnet_to_use_hub_vnet_gateways", "false"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.0.enable_internet_security", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_table.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.0", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.1", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.next_hop_ip_address", "10.0.4.6"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				Config: testAccAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_wan_id"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connections.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "route_table.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.0", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.address_prefixes.1", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route_table.0.routes.0.next_hop_ip_address", "10.0.4.5"),
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

func testCheckAzureRMVirtualHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.VirtualHubClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
		}

		return nil
	}
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

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHub_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
  name           = "acctestvirtualhub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location       = "${azurerm_resource_group.test.location}"
  address_prefix = "10.0.1.0/24"

  virtual_wan_id  = "${azurerm_virtual_wan.test.id}"

  route_table {
    routes {
      address_prefixes    = ["10.0.2.0/24", "10.0.3.0/24"]
      next_hop_ip_address = "10.0.4.5"
    }
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMVirtualHub_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.5.1.0/24"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
  name           = "acctestvirtualhub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location       = "${azurerm_resource_group.test.location}"
  address_prefix = "10.0.1.0/24"

  virtual_wan_id  = "${azurerm_virtual_wan.test.id}"

  virtual_network_connections {
	name = "testConnection"
	
    remote_virtual_network_id = "${azurerm_virtual_network.test.id}"

    allow_hub_to_remote_vnet_transit           = "false"
    allow_remote_vnet_to_use_hub_vnet_gateways = "false"
    enable_internet_security                   = "false"
  }

  route_table {
    routes {
      address_prefixes    = ["10.0.2.0/24", "10.0.3.0/24"]
      next_hop_ip_address = "10.0.4.6"
    }
  }

  tags = {
    env = "test"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}
