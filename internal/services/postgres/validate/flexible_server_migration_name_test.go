// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestFlexibleServerMigrationName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "flexdb%",
			valid: false,
		},
		{
			input: "1flexdb",
			valid: false,
		},
		{
			input: "a",
			valid: true,
		},
		{
			input: "flexdb_",
			valid: false,
		},
		{
			input: "flexdb",
			valid: true,
		},
		{
			input: "flexdb1test",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FlexibleServerMigrationName(tt.input, "name")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.input)
			}
		})
	}
}
