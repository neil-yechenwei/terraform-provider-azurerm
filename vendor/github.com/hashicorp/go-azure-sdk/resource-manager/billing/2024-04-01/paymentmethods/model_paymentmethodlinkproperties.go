package paymentmethods

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PaymentMethodLinkProperties struct {
	AccountHolderName *string                  `json:"accountHolderName,omitempty"`
	DisplayName       *string                  `json:"displayName,omitempty"`
	Expiration        *string                  `json:"expiration,omitempty"`
	Family            *PaymentMethodFamily     `json:"family,omitempty"`
	LastFourDigits    *string                  `json:"lastFourDigits,omitempty"`
	Logos             *[]PaymentMethodLogo     `json:"logos,omitempty"`
	PaymentMethod     *PaymentMethodProperties `json:"paymentMethod,omitempty"`
	PaymentMethodId   *string                  `json:"paymentMethodId,omitempty"`
	PaymentMethodType *string                  `json:"paymentMethodType,omitempty"`
	Status            *PaymentMethodStatus     `json:"status,omitempty"`
}
