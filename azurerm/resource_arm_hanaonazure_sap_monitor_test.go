package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMHanaOnAzureSapMonitor_basic(t *testing.T) {
	resourceName := "azurerm_hanaonazure_sap_monitor.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHanaOnAzureSapMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHanaOnAzureSapMonitor_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHanaOnAzureSapMonitorExists(resourceName),
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

func TestAccAzureRMHanaOnAzureSapMonitor_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_hanaonazure_sap_monitor.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHanaOnAzureSapMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHanaOnAzureSapMonitor_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHanaOnAzureSapMonitorExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMHanaOnAzureSapMonitor_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_hanaonazure_sap_monitor"),
			},
		},
	})
}

func TestAccAzureRMHanaOnAzureSapMonitor_complete(t *testing.T) {
	resourceName := "azurerm_hanaonazure_sap_monitor.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHanaOnAzureSapMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHanaOnAzureSapMonitor_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHanaOnAzureSapMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hana_db_sql_port", "30815"),
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

func TestAccAzureRMHanaOnAzureSapMonitor_update(t *testing.T) {
	resourceName := "azurerm_hanaonazure_sap_monitor.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHanaOnAzureSapMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHanaOnAzureSapMonitor_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHanaOnAzureSapMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hana_db_sql_port", "30815"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMHanaOnAzureSapMonitor_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHanaOnAzureSapMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "hana_db_sql_port", "30815"),
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

func testCheckAzureRMHanaOnAzureSapMonitorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Sap Monitor not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).HanaOnAzure.SapMonitorClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Sap Monitor %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on sapMonitorsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMHanaOnAzureSapMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).HanaOnAzure.SapMonitorClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hanaonazure_sap_monitor" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on sapMonitorsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMHanaOnAzureSapMonitor_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hanaonazure-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hanaonazure_sap_monitor" "test" {
  name                = "acctest-HanaOnAzureSapMonitor-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  hana_db_username    = "SYSTEM"
  hana_db_sql_port    = 30815
  hana_host_name      = "10.0.0.6"
  hana_db_name        = "SYSTEMDB"
  hana_db_password    = "Manager1"
  hana_subnet_id      = "${azurerm_subnet.test.id}"
}
`, rInt, location, rInt)
}

func testAccAzureRMHanaOnAzureSapMonitor_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hanaonazure_sap_monitor" "import" {
  name                = "${azurerm_hanaonazure_sap_monitor.test.name}"
  location            = "${azurerm_hanaonazure_sap_monitor.test.location}"
  resource_group_name = "${azurerm_hanaonazure_sap_monitor.test.name}"
}
`, testAccAzureRMHanaOnAzureSapMonitor_basic(rInt, location))
}

func testAccAzureRMHanaOnAzureSapMonitor_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hanaonazure-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hanaonazure_sap_monitor" "test" {
  name                = "acctest-HanaOnAzureSapMonitor-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  hana_db_username    = "SYSTEM"
  hana_db_sql_port    = 30815
  hana_host_name      = "10.0.0.6"
  hana_db_name        = "SYSTEMDB"
  hana_db_password    = "Manager1"
  hana_subnet_id      = "${azurerm_subnet.test.id}"
}
`, rInt, location, rInt)
}
