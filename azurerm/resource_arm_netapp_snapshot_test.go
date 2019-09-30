package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetAppSnapshot_basic(t *testing.T) {
	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
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

func TestAccAzureRMNetAppSnapshot_complete(t *testing.T) {
	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "file_system_id"),
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

func TestAccAzureRMNetAppSnapshot_update(t *testing.T) {
	resourceName := "azurerm_netapp_snapshot.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetAppSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetAppSnapshot_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctestsnapshot"),
				),
			},
			{
				Config: testAccAzureRMNetAppSnapshot_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetAppSnapshotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "newacctestsnapshot"),
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

func testCheckAzureRMNetAppSnapshotExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NetApp Snapshot not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		volumeName := rs.Primary.Attributes["volume_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).netapp.SnapshotClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: NetApp Snapshot %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on netapp.SnapshotClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetAppSnapshotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).netapp.SnapshotClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_netapp_snapshot" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		poolName := rs.Primary.Attributes["pool_name"]
		volumeName := rs.Primary.Attributes["volume_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on netapp.SnapshotClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNetAppSnapshot_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestaccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestpool-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctestvolume-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  creation_token      = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.test.id}"
  usage_threshold     = "107374182400"
}

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctestsnapshot"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_name         = "${azurerm_netapp_volume.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

func testAccAzureRMNetAppSnapshot_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestaccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestpool-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctestvolume-%d"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  creation_token      = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.test.id}"
  usage_threshold     = "107374182400"
}

resource "azurerm_netapp_snapshot" "test" {
  name                = "newacctestsnapshot"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_name         = "${azurerm_netapp_volume.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  file_system_id      = "${azurerm_netapp_volume.test.file_system_id}"
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}
