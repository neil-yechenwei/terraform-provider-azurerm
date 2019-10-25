package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func ValidatePrivateLinkServiceName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[\da-zA-Z]){1,}([\d\._\-a-zA-Z]{0,77})([\da-zA-Z_]$)`); !m {
		errors = append(regexErrs, fmt.Errorf(`%q must be between 1 and 80 characters, begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.`, k))
	}

	return nil, errors
}
