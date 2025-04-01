// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func FlexibleServerMigrationName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if !regexp.MustCompile(`^[a-z][a-z0-9]*$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q must start with lower characters and may contain numbers, lower characters, got %v", k, v))
		return
	}
	return
}
