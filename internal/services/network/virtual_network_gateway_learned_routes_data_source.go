package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualNetworkGatewayLearnedRoutesDataSource struct{}

var _ sdk.DataSource = VirtualNetworkGatewayLearnedRoutesDataSource{}

type VirtualNetworkGatewayLearnedRoutesModel struct {
	VirtualNetworkGatewayId string         `tfschema:"virtual_network_gateway_id"`
	GatewayRoutes           []GatewayRoute `tfschema:"gateway_routes"`
}

type GatewayRoute struct {
	AsPath       string `tfschema:"as_path"`
	LocalAddress string `tfschema:"local_address"`
	Network      string `tfschema:"network"`
	NextHop      string `tfschema:"next_hop"`
	Origin       string `tfschema:"origin"`
	SourcePeer   string `tfschema:"source_peer"`
	Weight       int64  `tfschema:"weight"`
}

func (VirtualNetworkGatewayLearnedRoutesDataSource) ResourceType() string {
	return "azurerm_virtual_network_gateway_learned_routes"
}

func (VirtualNetworkGatewayLearnedRoutesDataSource) ModelObject() interface{} {
	return &VirtualNetworkGatewayLearnedRoutesModel{}
}

func (VirtualNetworkGatewayLearnedRoutesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"virtual_network_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
		},
	}
}

func (VirtualNetworkGatewayLearnedRoutesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"gateway_routes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"as_path": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"local_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"network": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"next_hop": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"origin": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"source_peer": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"weight": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (VirtualNetworkGatewayLearnedRoutesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkGateways

			var state VirtualNetworkGatewayLearnedRoutesModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			vnetGatewayId, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(state.VirtualNetworkGatewayId)
			if err != nil {
				return err
			}

			resp, err := client.GetLearnedRoutes(ctx, *vnetGatewayId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", vnetGatewayId)
				}
				return fmt.Errorf("reading %s: %+v", vnetGatewayId, err)
			}

			if model := resp.Model; model != nil {
				state.GatewayRoutes = flattenVirtualNetworkGatewayLearnedRoutesGatewayRoutes(model.Value)
			}

			metadata.ResourceData.SetId(vnetGatewayId.ID())

			return metadata.Encode(&state)
		},
	}
}

func flattenVirtualNetworkGatewayLearnedRoutesGatewayRoutes(input *[]virtualnetworkgateways.GatewayRoute) []GatewayRoute {
	result := make([]GatewayRoute, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, GatewayRoute{
			AsPath:       pointer.From(v.AsPath),
			LocalAddress: pointer.From(v.LocalAddress),
			Network:      pointer.From(v.Network),
			NextHop:      pointer.From(v.NextHop),
			Origin:       pointer.From(v.Origin),
			SourcePeer:   pointer.From(v.SourcePeer),
			Weight:       pointer.From(v.Weight),
		})
	}

	return result
}
