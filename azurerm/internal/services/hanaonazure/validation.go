package hanaonazure

import (
	"fmt"
	"regexp"
)

func ValidateHanaSapMonitorName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z]{6,30}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 6 and 30 characters in length and contain only letters and numbers.", k))
	}

	return warnings, errors
}

func ValidateHanaDBName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\dA-Z]{2,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 64 characters in length and contain only uppercase letters and numbers.", k))
	}

	return warnings, errors
}

func ValidateHanaDBUserName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z]{1,32}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 32 characters in length and contain only letters and numbers.", k))
	}

	return warnings, errors
}

func ValidateHanaDBPassword(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{1,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 64 characters in length.", k))
	}

	return warnings, errors
}
