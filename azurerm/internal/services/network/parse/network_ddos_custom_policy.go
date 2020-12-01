package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkDDoSCustomPolicyId struct {
	ResourceGroup string
	Name          string
}

func NetworkDDoSCustomPolicyID(input string) (*NetworkDDoSCustomPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing networkDDoSCustomPolicy ID %q: %+v", input, err)
	}

	networkDDoSCustomPolicy := NetworkDDoSCustomPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if networkDDoSCustomPolicy.Name, err = id.PopSegment("ddosCustomPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &networkDDoSCustomPolicy, nil
}
