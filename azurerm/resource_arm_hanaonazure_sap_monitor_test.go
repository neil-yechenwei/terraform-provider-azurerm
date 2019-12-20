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
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "Test"),
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
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "Prod"),
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
	template := testAccAzureRMHanaOnAzureSapMonitor_template()
	return fmt.Sprintf(`
%s

resource "azurerm_hanaonazure_sap_monitor" "test" {
  name                = "acctest-HanaOnAzureSapMonitor-%d"
  resource_group_name = "${data.azurerm_resource_group.test.name}"
  location            = "${data.azurerm_resource_group.test.location}"
  hana_host_name      = "10.0.0.6"
  hana_subnet_id      = "${data.azurerm_subnet.test.id}"
  hana_db_name        = "SYSTEMDB"
  hana_db_sql_port    = 30215
  hana_db_username    = "SYSTEM"
  hana_db_password    = ""

  tags = {
    ENV = "Test"
  }
}
`, template, rInt)
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
	template := testAccAzureRMHanaOnAzureSapMonitor_template()
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test"{
	name                = "acctest-msi-%d"
	resource_group_name = "${data.azurerm_resource_group.test.name}"
	location            = "${data.azurerm_resource_group.test.location}"
}

resource "azurerm_hanaonazure_sap_monitor" "test"{
	name                           = "acctest-HanaOnAzureSapMonitor-%d"
	resource_group_name            = "${data.azurerm_resource_group.test.name}"
	location                       = "${data.azurerm_resource_group.test.location}"
	hana_host_name                 = "10.0.0.6"
	hana_subnet_id                 = "${data.azurerm_subnet.test.id}"
	hana_db_name                   = "SYSTEMDB"
	hana_db_sql_port               = 30215
	hana_db_username               = "SYSTEM"
	key_vault_id                   = ""
	hana_db_password_key_vault_url = ""
	hana_db_credentials_msi_id     = "${azurerm_user_assigned_identity.test.id}"

    tags = {
      ENV = "Prod"
    }
}
`, template, rInt, rInt)
}

func testAccAzureRMHanaOnAzureSapMonitor_template() string {
	return fmt.Sprintf(`
data "azurerm_resource_group" "test" {
  name = "HanaOnAzure-SapMonitor-RG5"
}

data "azurerm_subnet" "test" {
  name                 = "hdb-subnet"
  resource_group_name  = "${data.azurerm_resource_group.test.name}"
  virtual_network_name = "PV1-vnet"
}
`)
}
