package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/sdk/authorizations"
)

func ExpressRouteAuthorizationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := authorizations.ParseAuthorizationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
