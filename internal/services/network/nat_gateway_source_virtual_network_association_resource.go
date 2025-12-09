// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NATGatewaySourceVirtualNetworkAssociationModel struct {
	NATGatewayId     string `tfschema:"nat_gateway_id"`
	VirtualNetworkId string `tfschema:"virtual_network_id"`
}

type NATGatewaySourceVirtualNetworkAssociationResource struct{}

var _ sdk.Resource = NATGatewaySourceVirtualNetworkAssociationResource{}

func (r NATGatewaySourceVirtualNetworkAssociationResource) ModelObject() interface{} {
	return &NATGatewaySourceVirtualNetworkAssociationModel{}
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) ResourceType() string {
	return "azurerm_nat_gateway_source_virtual_network_association"
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": commonschema.ResourceIDReferenceRequiredForceNew(&natgateways.NatGatewayId{}),

		"virtual_network_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.VirtualNetworkId{}),
	}
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			var state NATGatewaySourceVirtualNetworkAssociationModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			virtualNetworkId, err := commonids.ParseVirtualNetworkID(state.VirtualNetworkId)
			if err != nil {
				return err
			}

			natGatewayId, err := natgateways.ParseNatGatewayID(state.NATGatewayId)
			if err != nil {
				return err
			}

			locks.ByName(natGatewayId.NatGatewayName, natGatewayResourceName)
			defer locks.UnlockByName(natGatewayId.NatGatewayName, natGatewayResourceName)

			natGateway, err := client.Get(ctx, *natGatewayId, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return fmt.Errorf("%s was not found", *natGatewayId)
				}
				return fmt.Errorf("retrieving %s: %+v", *natGatewayId, err)
			}

			id := commonids.NewCompositeResourceID(natGatewayId, virtualNetworkId)

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					if sourceVirtualNetwork := props.SourceVirtualNetwork; sourceVirtualNetwork != nil {
						if existingVnetId := sourceVirtualNetwork.Id; existingVnetId != nil {
							if strings.EqualFold(*existingVnetId, virtualNetworkId.ID()) {
								return metadata.ResourceRequiresImport(r.ResourceType(), id)
							}
							return fmt.Errorf("%s already has a source virtual network association with %s", *natGatewayId, *existingVnetId)
						}
					}

					natGateway.Model.Properties.SourceVirtualNetwork = &natgateways.SubResource{
						Id: pointer.To(virtualNetworkId.ID()),
					}
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *natGatewayId, *natGateway.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *natGatewayId, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &commonids.VirtualNetworkId{})
			if err != nil {
				return err
			}

			natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id.First, err)
			}

			var state NATGatewaySourceVirtualNetworkAssociationModel

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.SourceVirtualNetwork == nil || props.SourceVirtualNetwork.Id == nil {
						return metadata.MarkAsGone(id)
					}

					if !strings.EqualFold(*props.SourceVirtualNetwork.Id, id.Second.ID()) {
						return metadata.MarkAsGone(id)
					}

					state.NATGatewayId = id.First.ID()
					state.VirtualNetworkId = id.Second.ID()
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.NatGateways

			id, err := commonids.ParseCompositeResourceID(metadata.ResourceData.Id(), &natgateways.NatGatewayId{}, &commonids.VirtualNetworkId{})
			if err != nil {
				return err
			}

			locks.ByName(id.First.NatGatewayName, natGatewayResourceName)
			defer locks.UnlockByName(id.First.NatGatewayName, natGatewayResourceName)

			natGateway, err := client.Get(ctx, *id.First, natgateways.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(natGateway.HttpResponse) {
					return fmt.Errorf("%s was not found", *id.First)
				}
				return fmt.Errorf("retrieving %s: %+v", *id.First, err)
			}

			if model := natGateway.Model; model != nil {
				if props := model.Properties; props != nil {
					natGateway.Model.Properties.SourceVirtualNetwork = nil
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id.First, *natGateway.Model); err != nil {
				return fmt.Errorf("removing association between %s and %s: %+v", *id.First, *id.Second, err)
			}

			return nil
		},
	}
}

func (r NATGatewaySourceVirtualNetworkAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		_, err := commonids.ParseCompositeResourceID(input.(string), &natgateways.NatGatewayId{}, &commonids.VirtualNetworkId{})
		if err != nil {
			errors = append(errors, err)
		}
		return warnings, errors
	}
}
