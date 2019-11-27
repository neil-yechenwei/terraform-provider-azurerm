package hanaonazure

import (
	"fmt"
	"regexp"
)

func ValidateHanaOnAzureSapMonitorName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z]{6,30}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 6 and 30 characters in length and and may contain only letters and numbers.", k))
	}

	return warnings, errors
}
