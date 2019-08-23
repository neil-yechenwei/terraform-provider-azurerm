package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMPrivateLinkService_basic(t *testing.T) {
	resourceName := "azurerm_private_link_service.test"
	ri := tf.AccRandTimeInt()
	resourceGroupName := fmt.Sprintf("acctestRg_%d", ri)
	location := testLocation()
	serviceName := fmt.Sprintf("acctestPls_%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateLinkServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateLinkService_basic(resourceGroupName, location, serviceName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateLinkServiceExists(resourceName),
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

func testCheckAzureRMPrivateLinkServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		privateLinkService := rs.Primary.Attributes["name"]
		resourceGroupName := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).network.PrivateLinkServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroupName, privateLinkService, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on privateLinkServiceClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: private link service: %q does not exist", privateLinkService)
		}

		return nil
	}
}

func testCheckAzureRMPrivateLinkServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.PrivateLinkServiceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_link_service" {
			continue
		}

		privateLinkService := rs.Primary.Attributes["name"]
		resourceGroupName := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroupName, privateLinkService, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Private Link Service still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMPrivateLinkService_basic(resourceGroupName string, location string, serviceName string) string {
	return fmt.Sprintf(`
	resource "azurerm_resource_group" "test" {
		name     = "%s"
		location = "%s"
	  }
	  
	  resource "azurerm_private_link_service" "test" {
		  name = "%s"
		  location = "${azurerm_resource_group.test.location}"
		  resource_group_name = "${azurerm_resource_group.test.name}"
		  fqdns = ["testFqdns"]
		  ip_configuration {
			name = "${azurerm_public_ip.test.name}"
			subnet_id = "${azurerm_subnet.test.id}"
			private_ip_address = "10.0.1.17"
			private_ip_address_version = "IPv4"
			private_ip_address_allocation = "Static"
		  }
		  load_balancer_frontend_ip_configuration {
			 id = "${azurerm_lb.test.frontend_ip_configuration.0.id}"
		  }
		  tags = {
			env = "test"
		  }
	  }
	  
	  resource "azurerm_virtual_network" "test" {
		  name                = "testVnet"
		  address_space       = ["10.0.0.0/16"]
		  location            = "${azurerm_resource_group.test.location}"
		  resource_group_name = "${azurerm_resource_group.test.name}"
	  }
	  
	  resource "azurerm_subnet" "test" {
		  name                 = "testSubnet"
		  resource_group_name  = "${azurerm_resource_group.test.name}"
		  virtual_network_name = "${azurerm_virtual_network.test.name}"
		  address_prefix       = "10.0.1.0/24"
		  private_link_service_network_policies = "Disabled"
	  }
	  
	  resource "azurerm_public_ip" "test" {
		  name                = "testPip"
		  sku                 = "Standard"
		  location            = "${azurerm_resource_group.test.location}"
		  resource_group_name = "${azurerm_resource_group.test.name}"
		  allocation_method   = "Static"
	  }
	  
	  resource "azurerm_lb" "test" {
		  name                = "testLb"
		  sku                 = "Standard"
		  location            = "${azurerm_resource_group.test.location}"
		  resource_group_name = "${azurerm_resource_group.test.name}"
		  frontend_ip_configuration {
			  name                 = "${azurerm_public_ip.test.name}"
			  public_ip_address_id = "${azurerm_public_ip.test.id}"
		  }
	  }
`, resourceGroupName, location, serviceName)
}
