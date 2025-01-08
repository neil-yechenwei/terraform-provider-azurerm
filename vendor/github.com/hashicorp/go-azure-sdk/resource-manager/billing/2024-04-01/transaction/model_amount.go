package transaction

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Amount struct {
	Currency *string  `json:"currency,omitempty"`
	Value    *float64 `json:"value,omitempty"`
}
