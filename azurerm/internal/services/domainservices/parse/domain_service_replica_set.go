package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DomainServiceReplicaSetId struct {
	SubscriptionId    string
	ResourceGroup     string
	DomainServiceName string
	ReplicaSetName    string
}

func NewDomainServiceReplicaSetID(subscriptionId, resourceGroup, domainServiceName, replicaSetName string) DomainServiceReplicaSetId {
	return DomainServiceReplicaSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		DomainServiceName: domainServiceName,
		ReplicaSetName:    replicaSetName,
	}
}

func (id DomainServiceReplicaSetId) String() string {
	segments := []string{
		fmt.Sprintf("Replica Set Name %q", id.ReplicaSetName),
		fmt.Sprintf("Domain Service Name %q", id.DomainServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Domain Service Replica Set", segmentsStr)
}

func (id DomainServiceReplicaSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AAD/domainServices/%s/replicaSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, id.ReplicaSetName)
}

// DomainServiceReplicaSetID parses a DomainServiceReplicaSet ID into an DomainServiceReplicaSetId struct
func DomainServiceReplicaSetID(input string) (*DomainServiceReplicaSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DomainServiceReplicaSetId{
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
	if resourceId.ReplicaSetName, err = id.PopSegment("replicaSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
