---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-resource-private-link-service"
description: |-
  Manage Azure PrivateLinkService instance.
---

# azurerm_private_link_service

Manage Azure PrivateLinkService instance.


## Private Link Service Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acctestRG"
  location = "Eastus2"
}

resource "azurerm_virtual_network" "example" {
  name                = "acctestvnet-%d"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "acctestsnet-%d"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.5.1.0/24"
}

resource "azurerm_public_ip" "example" {
  name                = "acctestpip-%d"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "acctestlb-%d"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  frontend_ip_configuration {
    name                 = "${azurerm_public_ip.example.name}"
    public_ip_address_id = "${azurerm_public_ip.example.id}"
  }
}

resource "azurerm_private_link_service" "example" {
  name           = "acctestpls-%d"
  location       = "${azurerm_resource_group.example.location}"
  resource_group = "${azurerm_resource_group.example.name}"
  fqdns          = ["testFqdns"]

  ip_configurations {
    name = "${azurerm_public_ip.example.name}"

    subnet {
      id = "${azurerm_subnet.example.id}"
    }

    private_ip_address          = "10.5.1.17"
    private_ip_address_version  = "IPv4"
    private_ip_allocation_method = "Static"
  }

  load_balancer_frontend_ip_configurations {
    id = "${azurerm_lb.example.frontend_ip_configuration.0.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service. Changing this forces a new resource to be created.

* `resource_group` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `auto_approval` - (Optional) One `auto_approval` block defined below.

* `fqdns` - (Optional) The list of Fqdn.

* `ip_configurations` - (Optional) One or more `ip_configuration` block defined below.

* `load_balancer_frontend_ip_configurations` - (Optional) One or more `load_balancer_frontend_ip_configuration` block defined below. It only supports Internal LoadBalancer.

* `private_endpoint_connections` - (Optional) One or more `private_endpoint_connection` block defined below.

* `visibility` - (Optional) One `visibility` block defined below.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `auto_approval` block supports the following:

* `subscriptions` - (Optional) The list of subscriptions.

---

The `ip_configuration` block supports the following:

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `private_ip_allocation_method` - (Optional) The private IP address allocation method. Defaults to `Static`.

* `subnet` - (Optional) One `subnet` block defined below.

* `private_ip_address_version` - (Optional) Available from Api-Version 2016-03-30 onwards, it represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4. Defaults to `IPv4`.

* `name` - (Optional) The name of private link service ip configuration.

---

The `subnet` block supports the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `load_balancer_frontend_ip_configuration` block supports the following:

* `id` - (Optional) Resource ID.

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `private_ip_allocation_method` - The Private IP allocation method.

* `private_ip_address_version` - It represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.

* `subnet` - (Optional) One `subnet` block defined below.

* `public_ip_address` - (Optional) One `public_ip_address` block defined below.

* `public_ip_prefix` - (Optional) One `public_ipprefix` block defined below.

* `name` - (Optional) The name of the resource that is unique within the set of frontend IP configurations used by the load balancer. This name can be used to access the resource.

* `zones` - A list of availability zones denoting the IP allocated for the resource needs to come from.

---

The `subnet` block supports the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `public_ip_address` block supports the following:

* `id` - (Optional) Resource ID.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags.

* `zones` - (Optional) A list of availability zones denoting the IP allocated for the resource needs to come from.

---

The `public_ip_prefix` block supports the following:

* `id` - (Optional) Resource ID.

---

The `private_endpoint_connection` block supports the following:

* `id` - (Optional) Resource ID.

* `private_endpoint` - (Optional) One `private_endpoint` block defined below.

* `private_link_service_connection_state` - (Optional) One `private_link_service_connection_state` block defined below.

* `name` - (Optional) The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `private_endpoint` block supports the following:

* `id` - (Optional) Resource ID.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags.

---

The `private_link_service_connection_state` block supports the following:

* `status` - (Optional) Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

* `description` - (Optional) The reason for approval/rejection of the connection.

* `action_required` - (Optional) A message indicating if changes on the service provider require any updates on the consumer.

---

The `visibility` block supports the following:

* `subscriptions` - (Optional) The list of subscriptions.

## Attributes Reference

The following attributes are exported:

* `network_interfaces` - One or more `network_interface` block defined below.

* `alias` - The alias of the private link service.

* `type` - Resource type.

---

The `network_interface` block contains the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) Resource name.

* `type` - (Optional) Resource type.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags.

* `virtual_machine` - (Optional) One `virtual_machine` block defined below.

* `network_security_group` - (Optional) One `network_security_group` block defined below.

* `private_endpoint` - (Optional) One `private_endpoint` block defined below.

* `ip_configurations` - (Optional) One `ip_configuration` block defined below.

* `tap_configurations` - (Optional) One `tap_configuration` block defined below.

* `dns_settings` - (Optional) One `dns_setting` block defined below.

* `mac_address` - (Optional) The MAC address of the network interface.

* `primary` - (Optional) Gets whether this is a primary network interface on a virtual machine.

* `enable_accelerated_networking` - (Optional) If the network interface is accelerated networking enabled.

* `enable_ip_forwarding` - (Optional) Indicates whether IP forwarding is enabled on this network interface.

* `hosted_workloads` - (Optional) A list of references to linked BareMetal resources.

* `resource_guid` - (Optional) The resource GUID property of the network interface resource.

---

The `virtual_machine` block supports the following:

* `id` - (Optional) Resource ID.

---

The `network_security_group` block supports the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) Resource name.

* `type` - (Optional) Resource type.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags.

---

The `private_endpoint` block supports the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) Resource name.

* `type` - (Optional) Resource type.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags.

---

The `ip_configuration` block supports the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `tap_configuration` block supports the following:

* `id` - (Optional) Resource ID.

* `name` - (Optional) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `type` - (Optional) Sub Resource type.

---

The `dns_setting` block supports the following:

* `dns_servers` - (Optional) List of DNS servers IP addresses. Use 'AzureProvidedDNS' to switch to azure provided DNS resolution. 'AzureProvidedDNS' value cannot be combined with other IPs, it must be the only value in dnsServers collection.

* `applied_dns_servers` - (Optional) If the VM that uses this NIC is part of an Availability Set, then this list will have the union of all DNS servers from all NICs that are part of the Availability Set. This property is what is configured on each of those VMs.

* `internal_dns_name_label` - (Optional) Relative DNS name for this NIC used for internal communications between VMs in the same virtual network.

* `internal_fqdn` - (Optional) Fully qualified DNS name supporting internal communications between VMs in the same virtual network.

* `internal_domain_name_suffix` - (Optional) Even if internalDnsNameLabel is not specified, a DNS entry is created for the primary NIC of the VM. This DNS name can be constructed by concatenating the VM name with the value of internalDomainNameSuffix.

## Import

Private Link Service can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.Network/privateLinkServices/
```
