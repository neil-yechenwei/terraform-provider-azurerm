package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DomainServiceOuContainerId struct {
	SubscriptionId    string
	ResourceGroup     string
	DomainServiceName string
	OuContainerName   string
}

func NewDomainServiceOuContainerID(subscriptionId, resourceGroup, domainServiceName, ouContainerName string) DomainServiceOuContainerId {
	return DomainServiceOuContainerId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		DomainServiceName: domainServiceName,
		OuContainerName:   ouContainerName,
	}
}

func (id DomainServiceOuContainerId) String() string {
	segments := []string{
		fmt.Sprintf("Ou Container Name %q", id.OuContainerName),
		fmt.Sprintf("Domain Service Name %q", id.DomainServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Domain Service Ou Container", segmentsStr)
}

func (id DomainServiceOuContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AAD/domainServices/%s/ouContainer/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, id.OuContainerName)
}

// DomainServiceOuContainerID parses a DomainServiceOuContainer ID into an DomainServiceOuContainerId struct
func DomainServiceOuContainerID(input string) (*DomainServiceOuContainerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DomainServiceOuContainerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DomainServiceName, err = id.PopSegment("domainServices"); err != nil {
		return nil, err
	}
	if resourceId.OuContainerName, err = id.PopSegment("ouContainer"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
