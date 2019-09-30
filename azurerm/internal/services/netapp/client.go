package netapp

import (
	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient  *netapp.AccountsClient
	PoolClient     *netapp.PoolsClient
	VolumeClient   *netapp.VolumesClient
	SnapshotClient *netapp.SnapshotsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	AccountClient := netapp.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AccountClient.Client, o.ResourceManagerAuthorizer)

	PoolClient := netapp.NewPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PoolClient.Client, o.ResourceManagerAuthorizer)

	VolumeClient := netapp.NewVolumesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VolumeClient.Client, o.ResourceManagerAuthorizer)

	SnapshotClient := netapp.NewSnapshotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SnapshotClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:  &AccountClient,
		PoolClient:     &PoolClient,
		VolumeClient:   &VolumeClient,
		SnapshotClient: &SnapshotClient,
	}
}
