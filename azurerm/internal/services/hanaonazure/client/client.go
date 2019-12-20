package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/hanaonazure/mgmt/2017-11-03-preview/hanaonazure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	SapMonitorClient *hanaonazure.SapMonitorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	sapMonitorClient := hanaonazure.NewSapMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sapMonitorClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SapMonitorClient: &sapMonitorClient,
	}
}
