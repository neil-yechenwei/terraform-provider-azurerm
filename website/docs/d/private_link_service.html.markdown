---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-datasource-private-link-service"
description: |-
  Gets information about an existing Private Link Service
---

# Data Source: azurerm_private_link_service

Use this data source to access information about an existing Private Link Service.


## Private Link Service Usage

```hcl
data "azurerm_private_link_service" "example" {
  resource_group = "acctestRG"
  name           = "acctestpls"
}

output "private_link_service_id" {
  value = "${data.azurerm_private_link_service.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service.

* `resource_group` - (Required) The name of the resource group.


## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `alias` - The alias of the private link service.

* `auto_approval` - One `auto_approval` block defined below.

* `fqdns` - The list of Fqdn.

* `ip_configurations` - One or more `ip_configuration` block defined below.

* `load_balancer_frontend_ip_configurations` - One or more `load_balancer_frontend_ip_configuration` block defined below. It only supports Internal LoadBalancer.

* `network_interfaces` - One or more `network_interface` block defined below.

* `private_endpoint_connections` - One or more `private_endpoint_connection` block defined below.

* `type` - Resource type.

* `visibility` - One `visibility` block defined below.

* `tags` - Resource tags.


---

The `auto_approval` block contains the following:

* `subscriptions` - The list of subscriptions.

---

The `ip_configuration` block contains the following:

* `private_ip_address` - The private IP address of the IP configuration.

* `private_ip_allocation_method` - The private IP address allocation method.

* `subnet` - One `subnet` block defined below.

* `private_ip_address_version` - Available from Api-Version 2016-03-30 onwards, it represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.

* `name` - The name of private link service ip configuration.


---

The `subnet` block contains the following:

* `id` - Resource ID.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `load_balancer_frontend_ip_configuration` block contains the following:

* `id` - Resource ID.

* `private_ip_address` - The private IP address of the IP configuration.

* `private_ip_allocation_method` - The Private IP allocation method.

* `private_ip_address_version` - It represents whether the specific ipconfiguration is IPv4 or IPv6. Default is taken as IPv4.

* `subnet` - One `subnet` block defined below.

* `public_ip_address` - One `public_ip_address` block defined below.

* `public_ip_prefix` - One `public_ipprefix` block defined below.

* `name` - The name of the resource that is unique within the set of frontend IP configurations used by the load balancer. This name can be used to access the resource.

* `zones` - A list of availability zones denoting the IP allocated for the resource needs to come from.


---

The `subnet` block contains the following:

* `id` - Resource ID.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `public_ip_address` block contains the following:

* `id` - Resource ID.

* `location` - Resource location.

* `tags` - Resource tags.

* `zones` - A list of availability zones denoting the IP allocated for the resource needs to come from.

---

The `public_ip_prefix` block contains the following:

* `id` - Resource ID.

---

The `network_interface` block contains the following:

* `id` - Resource ID.

* `name` - Resource name.

* `type` - Resource type.

* `location` - Resource location.

* `tags` - Resource tags.

* `virtual_machine` - One `virtual_machine` block defined below.

* `network_security_group` - One `network_security_group` block defined below.

* `private_endpoint` - One `private_endpoint` block defined below.

* `ip_configurations` - One `ip_configuration` block defined below.

* `tap_configurations` - One `tap_configuration` block defined below.

* `dns_settings` - One `dns_setting` block defined below.

* `mac_address` - The MAC address of the network interface.

* `primary` - Gets whether this is a primary network interface on a virtual machine.

* `enable_accelerated_networking` - If the network interface is accelerated networking enabled.

* `enable_ip_forwarding` - Indicates whether IP forwarding is enabled on this network interface.

* `hosted_workloads` - A list of references to linked BareMetal resources.

* `resource_guid` - The resource GUID property of the network interface resource.

---

The `virtual_machine` block contains the following:

* `id` - Resource ID.

---

The `network_security_group` block contains the following:

* `id` - Resource ID.

* `name` - Resource name.

* `type` - Resource type.

* `location` - Resource location.

* `tags` - Resource tags.

---

The `private_endpoint` block contains the following:

* `id` - Resource ID.

* `name` - Resource name.

* `type` - Resource type.

* `location` - Resource location.

* `tags` - Resource tags.

---

The `ip_configuration` block contains the following:

* `id` - Resource ID.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

The `tap_configuration` block contains the following:

* `id` - Resource ID.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `type` - Sub Resource type.

---

The `dns_setting` block contains the following:

* `dns_servers` - List of DNS servers IP addresses. Use 'AzureProvidedDNS' to switch to azure provided DNS resolution. 'AzureProvidedDNS' value cannot be combined with other IPs, it must be the only value in dnsServers collection.

* `applied_dns_servers` - If the VM that uses this NIC is part of an Availability Set, then this list will have the union of all DNS servers from all NICs that are part of the Availability Set. This property is what is configured on each of those VMs.

* `internal_dns_name_label` - Relative DNS name for this NIC used for internal communications between VMs in the same virtual network.

* `internal_fqdn` - Fully qualified DNS name supporting internal communications between VMs in the same virtual network.

* `internal_domain_name_suffix` - Even if internalDnsNameLabel is not specified, a DNS entry is created for the primary NIC of the VM. This DNS name can be constructed by concatenating the VM name with the value of internalDomainNameSuffix.

---

The `private_endpoint_connection` block contains the following:

* `id` - Resource ID.

* `private_endpoint` - One `private_endpoint` block defined below.

* `private_link_service_connection_state` - One `private_link_service_connection_state` block defined below.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.


---

The `private_endpoint` block contains the following:

* `id` - Resource ID.

* `location` - Resource location.

* `tags` - Resource tags.

---

The `private_link_service_connection_state` block contains the following:

* `status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

* `description` - The reason for approval/rejection of the connection.

* `action_required` - A message indicating if changes on the service provider require any updates on the consumer.

---

The `visibility` block contains the following:

* `subscriptions` - The list of subscriptions.
