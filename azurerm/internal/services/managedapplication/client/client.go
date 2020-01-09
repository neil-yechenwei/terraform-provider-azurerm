package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-06-01/managedapplications"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationDefinitionClient *managedapplications.ApplicationDefinitionsClient
	ApplicationClient           *managedapplications.ApplicationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	applicationDefinitionClient := managedapplications.NewApplicationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationDefinitionClient.Client, o.ResourceManagerAuthorizer)

	applicationClient := managedapplications.NewApplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationDefinitionClient: &applicationDefinitionClient,
		ApplicationClient:           &applicationClient,
	}
}
