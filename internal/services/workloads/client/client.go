package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorsClient *monitors.MonitorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorsClient := monitors.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(monitorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorsClient: &monitorsClient,
	}
}
