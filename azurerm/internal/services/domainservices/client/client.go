package client

import (
	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2020-01-01/aad"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DomainServicesClient *aad.DomainServicesClient
	OuContainerClient    *aad.OuContainerClient
}

func NewClient(o *common.ClientOptions) *Client {
	domainServicesClient := aad.NewDomainServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&domainServicesClient.Client, o.ResourceManagerAuthorizer)

	ouContainerClient := aad.NewOuContainerClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ouContainerClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DomainServicesClient: &domainServicesClient,
		OuContainerClient:    &ouContainerClient,
	}
}
