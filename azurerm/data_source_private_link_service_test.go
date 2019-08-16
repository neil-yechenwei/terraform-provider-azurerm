package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateLinkService_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceGroupName := fmt.Sprintf("acctestRg_%d", ri)
	location := testLocation()
	serviceName := fmt.Sprintf("acctestPls_%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPrivateLinkServiceBasic(resourceGroupName, location, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azurerm_private_link_service.test", "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPrivateLinkServiceBasic(resourceGroupName string, location string, serviceName string) string {
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

	  data "azurerm_private_link_service" "test" {
	  	  name = "${azurerm_private_link_service.test.name}"
		  resource_group_name = "${azurerm_resource_group.test.name}"
	  }
`, resourceGroupName, location, serviceName)
}
