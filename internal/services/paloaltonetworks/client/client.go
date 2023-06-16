package client

import (
	paloaltonetworks_v2022_08_29 "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*paloaltonetworks_v2022_08_29.Client, error) {
	client, err := paloaltonetworks_v2022_08_29.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
