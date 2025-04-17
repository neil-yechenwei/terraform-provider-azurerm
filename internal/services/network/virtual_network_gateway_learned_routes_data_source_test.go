// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type VirtualNetworkGatewayLearnedRoutesDataSource struct{}

func TestAccAzureRMDataSourceVirtualNetworkGatewayLearnedRoutes_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_learned_routes", "test")
	r := VirtualNetworkGatewayLearnedRoutesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (VirtualNetworkGatewayLearnedRoutesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_virtual_network_gateway_learned_routes" "test" {
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
}
`, VirtualNetworkGatewayResource{}.activeActiveEnableBgpWithAPIPA(data))
}
