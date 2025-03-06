// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestBillingAccountName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "a",
			Expected: false,
		},
		{
			Input:    "123",
			Expected: true,
		},
		{
			Input:    "PCN.adf",
			Expected: true,
		},
		{
			Input:    "447ef28d-ce44-4db1-b814-343395118e95",
			Expected: true,
		},
	}

	for _, v := range testCases {
		_, errors := BillingAccountName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
