package agreement

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Participant struct {
	Email      *string `json:"email,omitempty"`
	Status     *string `json:"status,omitempty"`
	StatusDate *string `json:"statusDate,omitempty"`
}

func (o *Participant) GetStatusDateAsTime() (*time.Time, error) {
	if o.StatusDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StatusDate, "2006-01-02T15:04:05Z07:00")
}

func (o *Participant) SetStatusDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StatusDate = &formatted
}
