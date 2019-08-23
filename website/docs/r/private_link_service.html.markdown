---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-resource-private-link-service"
description: |-
  Manages a Private Link Service.

---

# azurerm_private_link_service

Manages a Private Link Service.

## Example Usage

```hcl
resource "azurerm_vnet_private_link_service" "test" {
    name                                        = "neilterraformpls01"
    location                                    = "${azurerm_resource_group.test.location}"
    resource_group_name                         = "${azurerm_resource_group.test.name}"
    auto_approval_subscription_names            = ["00000000-0000-0000-0000-000000000000"]
    fqdns                                       = ["fqdntest"]
    ip_configuration                            = {
        name = "neilpubliciptest01"
        subnet_id                               = "${azurerm_subnet.test2.id}"
        private_ip_address                      = "13.77.105.127"
        private_ip_address_version              = "IPv4"
        private_ip_address_allocation           = "Static"
    }
    load_balancer_frontend_ip_configuration     = "${azurerm_lb.test.private_ip_addresses}"
    visibility_subscription_names               = ["00000000-0000-0000-0000-000000000000"]
    tags = {
      env = "test"
    }
}

resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "test" {
    name                = "acceptanceTestVirtualNetwork1"
    address_space       = ["10.0.0.0/16"]
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
 
resource "azurerm_subnet" "test" {
    name                 = "testsubnet"
    resource_group_name  = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix       = "10.0.1.0/24"
}
 
resource "azurerm_subnet" "test2" {
    name                 = "testsubnet"
    resource_group_name  = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix       = "10.0.1.0/24"
}
 
resource "azurerm_public_ip" "test" {
    name                = "acceptanceTestPublicIp1"
    location            = "West US 2"
    resource_group_name = "${azurerm_resource_group.test.name}"
    allocation_method   = "Static"
}
 
resource "azurerm_lb" "test" {
    name                = "TestLoadBalancer"
    location            = "West US 2"
    resource_group_name = "${azurerm_resource_group.test.name}"
 
    frontend_ip_configuration {
        name                 = "PublicIPAddress"
        public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service. Changing this forces a new resource to be created.

* `location` - (Optional) The location/region where the private link service is created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the private link service. Changing this forces a new resource to be created.

* `ip_configuration` - (Optional) An array of references to the private link service IP configuration.

* `load_balancer_frontend_ip_configuration` - (Optional) An array of references to the load balancer IP configurations. It only supports Internal Load Balancer.

* `auto_approval_subscription_names` - (Optional) The auto-approval list of the private link service.

* `fqdns` - (Optional) The list of Fqdn.

* `visibility_subscription_names` - (Optional) The visibility list of the private link service.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `ip_configuration` block supports:

* `name` - (Optional) The name of private link service ip configuration.

* `subnet_id` - (Optional) The reference of the subnet resource.

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `private_ip_address_allocation` - (Optional) The private IP address allocation method.

* `private_ip_address_version` - (Optional) Available from Api-Version 2016-03-30 onwards, it represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.

The `load_balancer_frontend_ip_configuration` block supports:

* `id` - (Optional) Resource ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Link Service.

## Import

Private Link Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_link_service.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.network/privateLinkServices/myprivatelinkservice1
```
