package enrollmentaccount

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnrollmentAccountProperties struct {
	AccountOwner          *string `json:"accountOwner,omitempty"`
	AuthType              *string `json:"authType,omitempty"`
	CostCenter            *string `json:"costCenter,omitempty"`
	DepartmentDisplayName *string `json:"departmentDisplayName,omitempty"`
	DepartmentId          *string `json:"departmentId,omitempty"`
	DisplayName           *string `json:"displayName,omitempty"`
	EndDate               *string `json:"endDate,omitempty"`
	IsDevTestEnabled      *bool   `json:"isDevTestEnabled,omitempty"`
	StartDate             *string `json:"startDate,omitempty"`
	Status                *string `json:"status,omitempty"`
}

func (o *EnrollmentAccountProperties) GetEndDateAsTime() (*time.Time, error) {
	if o.EndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *EnrollmentAccountProperties) SetEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDate = &formatted
}

func (o *EnrollmentAccountProperties) GetStartDateAsTime() (*time.Time, error) {
	if o.StartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *EnrollmentAccountProperties) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = &formatted
}
