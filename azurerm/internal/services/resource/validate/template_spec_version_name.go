package validate

import (
	"fmt"
	"regexp"
)

func TemplateSpecVersionName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 {
		errors = append(errors, fmt.Errorf("length should be greater than %d", 1))
		return
	}

	if len(v) > 90 {
		errors = append(errors, fmt.Errorf("length should be less than %d", 90))
		return
	}

	if !regexp.MustCompile(`^[-\w\._\(\)]+$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%q only contains alpha-numeric characters, parenthesis, underscores, dashes and periods", k))
		return
	}

	return
}
