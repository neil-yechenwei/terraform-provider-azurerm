package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventHubNamespacePrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	NamespaceName                 string
	PrivateEndpointConnectionName string
}

func NewEventHubNamespacePrivateEndpointConnectionID(subscriptionId, resourceGroup, namespaceName, privateEndpointConnectionName string) EventHubNamespacePrivateEndpointConnectionId {
	return EventHubNamespacePrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		NamespaceName:                 namespaceName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

func (id EventHubNamespacePrivateEndpointConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Private Endpoint Connection Name %q", id.PrivateEndpointConnectionName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Event Hub Namespace Private Endpoint Connection", segmentsStr)
}

func (id EventHubNamespacePrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.PrivateEndpointConnectionName)
}

// EventHubNamespacePrivateEndpointConnectionID parses a EventHubNamespacePrivateEndpointConnection ID into an EventHubNamespacePrivateEndpointConnectionId struct
func EventHubNamespacePrivateEndpointConnectionID(input string) (*EventHubNamespacePrivateEndpointConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EventHubNamespacePrivateEndpointConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.PrivateEndpointConnectionName, err = id.PopSegment("privateEndpointConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
