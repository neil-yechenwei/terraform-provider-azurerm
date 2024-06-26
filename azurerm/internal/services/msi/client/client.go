package client

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/sdk/managedidentity"
)

type Client struct {
	UserAssignedIdentitiesClient *managedidentity.ManagedIdentityClient
}

func NewClient(o *common.ClientOptions) *Client {
	UserAssignedIdentitiesClient := managedidentity.NewManagedIdentityClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&UserAssignedIdentitiesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		UserAssignedIdentitiesClient: &UserAssignedIdentitiesClient,
	}
}
