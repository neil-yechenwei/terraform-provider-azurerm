package customer

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomerProperties struct {
	BillingProfileDisplayName *string            `json:"billingProfileDisplayName,omitempty"`
	BillingProfileId          *string            `json:"billingProfileId,omitempty"`
	DisplayName               *string            `json:"displayName,omitempty"`
	EnabledAzurePlans         *[]AzurePlan       `json:"enabledAzurePlans,omitempty"`
	Resellers                 *[]Reseller        `json:"resellers,omitempty"`
	Status                    *CustomerStatus    `json:"status,omitempty"`
	SystemId                  *string            `json:"systemId,omitempty"`
	Tags                      *map[string]string `json:"tags,omitempty"`
}
