package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type IntegrationServiceEnvironmentManagedApiId struct {
	SubscriptionId                    string
	ResourceGroup                     string
	IntegrationServiceEnvironmentName string
	ManagedApiName                    string
}

func NewIntegrationServiceEnvironmentManagedApiID(subscriptionId, resourceGroup, integrationServiceEnvironmentName, managedApiName string) IntegrationServiceEnvironmentManagedApiId {
	return IntegrationServiceEnvironmentManagedApiId{
		SubscriptionId:                    subscriptionId,
		ResourceGroup:                     resourceGroup,
		IntegrationServiceEnvironmentName: integrationServiceEnvironmentName,
		ManagedApiName:                    managedApiName,
	}
}

func (id IntegrationServiceEnvironmentManagedApiId) String() string {
	segments := []string{
		fmt.Sprintf("Managed Api Name %q", id.ManagedApiName),
		fmt.Sprintf("Integration Service Environment Name %q", id.IntegrationServiceEnvironmentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Service Environment Managed Api", segmentsStr)
}

func (id IntegrationServiceEnvironmentManagedApiId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationServiceEnvironments/%s/managedApis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationServiceEnvironmentName, id.ManagedApiName)
}

// IntegrationServiceEnvironmentManagedApiID parses a IntegrationServiceEnvironmentManagedApi ID into an IntegrationServiceEnvironmentManagedApiId struct
func IntegrationServiceEnvironmentManagedApiID(input string) (*IntegrationServiceEnvironmentManagedApiId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationServiceEnvironmentManagedApiId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IntegrationServiceEnvironmentName, err = id.PopSegment("integrationServiceEnvironments"); err != nil {
		return nil, err
	}
	if resourceId.ManagedApiName, err = id.PopSegment("managedApis"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
