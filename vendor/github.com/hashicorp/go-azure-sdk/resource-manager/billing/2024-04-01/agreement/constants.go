package agreement

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcceptanceMode string

const (
	AcceptanceModeClickToAccept AcceptanceMode = "ClickToAccept"
	AcceptanceModeESignEmbedded AcceptanceMode = "ESignEmbedded"
	AcceptanceModeESignOffline  AcceptanceMode = "ESignOffline"
	AcceptanceModeImplicit      AcceptanceMode = "Implicit"
	AcceptanceModeOffline       AcceptanceMode = "Offline"
	AcceptanceModeOther         AcceptanceMode = "Other"
	AcceptanceModePhysicalSign  AcceptanceMode = "PhysicalSign"
)

func PossibleValuesForAcceptanceMode() []string {
	return []string{
		string(AcceptanceModeClickToAccept),
		string(AcceptanceModeESignEmbedded),
		string(AcceptanceModeESignOffline),
		string(AcceptanceModeImplicit),
		string(AcceptanceModeOffline),
		string(AcceptanceModeOther),
		string(AcceptanceModePhysicalSign),
	}
}

func (s *AcceptanceMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAcceptanceMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAcceptanceMode(input string) (*AcceptanceMode, error) {
	vals := map[string]AcceptanceMode{
		"clicktoaccept": AcceptanceModeClickToAccept,
		"esignembedded": AcceptanceModeESignEmbedded,
		"esignoffline":  AcceptanceModeESignOffline,
		"implicit":      AcceptanceModeImplicit,
		"offline":       AcceptanceModeOffline,
		"other":         AcceptanceModeOther,
		"physicalsign":  AcceptanceModePhysicalSign,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AcceptanceMode(input)
	return &out, nil
}

type Category string

const (
	CategoryAffiliatePurchaseTerms         Category = "AffiliatePurchaseTerms"
	CategoryIndirectForGovernmentAgreement Category = "IndirectForGovernmentAgreement"
	CategoryMicrosoftCustomerAgreement     Category = "MicrosoftCustomerAgreement"
	CategoryMicrosoftPartnerAgreement      Category = "MicrosoftPartnerAgreement"
	CategoryOther                          Category = "Other"
	CategoryUKCloudComputeFramework        Category = "UKCloudComputeFramework"
)

func PossibleValuesForCategory() []string {
	return []string{
		string(CategoryAffiliatePurchaseTerms),
		string(CategoryIndirectForGovernmentAgreement),
		string(CategoryMicrosoftCustomerAgreement),
		string(CategoryMicrosoftPartnerAgreement),
		string(CategoryOther),
		string(CategoryUKCloudComputeFramework),
	}
}

func (s *Category) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCategory(input string) (*Category, error) {
	vals := map[string]Category{
		"affiliatepurchaseterms":         CategoryAffiliatePurchaseTerms,
		"indirectforgovernmentagreement": CategoryIndirectForGovernmentAgreement,
		"microsoftcustomeragreement":     CategoryMicrosoftCustomerAgreement,
		"microsoftpartneragreement":      CategoryMicrosoftPartnerAgreement,
		"other":                          CategoryOther,
		"ukcloudcomputeframework":        CategoryUKCloudComputeFramework,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Category(input)
	return &out, nil
}
