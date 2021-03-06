package policyinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// PolicyStatesResource enumerates the values for policy states resource.
type PolicyStatesResource string

const (
	// Default ...
	Default PolicyStatesResource = "default"
	// Latest ...
	Latest PolicyStatesResource = "latest"
)

// PossiblePolicyStatesResourceValues returns an array of possible values for the PolicyStatesResource const type.
func PossiblePolicyStatesResourceValues() []PolicyStatesResource {
	return []PolicyStatesResource{Default, Latest}
}

// ResourceDiscoveryMode enumerates the values for resource discovery mode.
type ResourceDiscoveryMode string

const (
	// ExistingNonCompliant Remediate resources that are already known to be non-compliant.
	ExistingNonCompliant ResourceDiscoveryMode = "ExistingNonCompliant"
	// ReEvaluateCompliance Re-evaluate the compliance state of resources and then remediate the resources
	// found to be non-compliant.
	ReEvaluateCompliance ResourceDiscoveryMode = "ReEvaluateCompliance"
)

// PossibleResourceDiscoveryModeValues returns an array of possible values for the ResourceDiscoveryMode const type.
func PossibleResourceDiscoveryModeValues() []ResourceDiscoveryMode {
	return []ResourceDiscoveryMode{ExistingNonCompliant, ReEvaluateCompliance}
}
