---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_connected_kubernetes_cluster"
description: |-
Manages a Connected Kubernetes Cluster.
---

# azurerm_connected_kubernetes_cluster

Manages a Connected Kubernetes Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_connected_kubernetes_cluster" "test" {
  name                         = "example-connectedkubernetescluster"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  agent_public_key_certificate = "DIICYzCCAcygAwIBAgIBADANBgkqhkiG9w0BAQUFADAuMQswCQYDVQQGEwJVUzEMMAoGA1UEChMDSUJNMREwDwYDVQQLEwhMb2NhbCBDQTAeFw05OTEyMjIwNTAwMDBaFw0wMDEyMjMwNDU5NTlaMC4xCzAJBgNVBAYTAlVTMQwwCgYDVQQKEwNJQk0xETAPBgNVBAsTCExvY2FsIENBMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQD2bZEo7xGaX2/0GHkrNFZvlxBou9v1Jmt/PDiTMPve8r9FeJAQ0QdvFST/0JPQYD20rH0bimdDLgNdNynmyRoS2S/IInfpmf69iyc2G0TPyRvmHIiOZbdCd+YBHQi1adkj17NDcWj6S14tVurFX73zx0sNoMS79q3tuXKrDsxeuwIDAQABo4GQMIGNMEsGCVUdDwGG+EIBDQQ+EzxHZW5lcmF0ZWQgYnkgdGhlIFNlY3VyZVdheSBTZWN1cml0eSBTZXJ2ZXIgZm9yIE9TLzM5MCAoUkFDRikwDgYDVR0PAQH/BAQDAgAGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFJ3+ocRyCTJw067dLSwr/nalx6YMMA0GCSqGSIb3DQEBBQUAA4GBAMaQzt+zaj1GU77yzlr8iiMBXgdQrwsZZWJo5exnAucJAEYQZmOfyLiM D6oYq+ZnfvM0n8G/Y79q8nhwvuxpYOnRSAXFp6xSkrIOeZtJMY1h00LKp/JX3Ng1svZ2agE126JHsQ0bhzN5TKsYfbwfTwfjdWAGy6Vf1nYi/rO+ryMO"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Connected Kubernetes Cluster resource. Changing this forces a new Connected Kubernetes Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Connected Kubernetes Cluster should exist. Changing this forces a new Connected Kubernetes Cluster to be created.

* `location` - (Required) The Azure Region where the Connected Kubernetes Cluster should exist. Changing this forces a new Connected Kubernetes Cluster to be created.

* `agent_public_key_certificate` - (Required) The base64 encoded Public Certificate used by the agent to do the initial handshake to the backend services in Azure.

---

* `identity_type` - (Optional) The type of identity used for the Connected Kubernetes Cluster. Possible values are `None` and `SystemAssigned`. Defaults to `SystemAssigned`.

* `distribution` - (Optional) The Kubernetes distribution which will be running on this Connected Kubernetes Cluster. Possible values are `aks`, `aks_engine`, `aks_hci`, `auto`, `capz`, `eks`, `generic`, `gke`, `k3s`, `kind`, `minikube`, `openshift`, `rancher_rke` and `tkg`.

* `infrastructure` - (Optional) The infrastructure on which the Kubernetes cluster represented by this Connected Kubernetes Cluster will be running on. Possible values are `auto`, `aws`, `azure`, `azure_stack_edge`, `azure_stack_hci`, `azure_stack_hub`, `gcp`, `generic` and `vsphere`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Communication Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Connected Kubernetes Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Connected Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Connected Kubernetes Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Connected Kubernetes Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Connected Kubernetes Cluster.

## Import

Connected Kubernetes Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_connected_kubernetes_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kubernetes/connectedClusters/connectedCluster1
```
