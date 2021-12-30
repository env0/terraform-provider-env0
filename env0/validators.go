package env0

import (
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
)

func ValidateConfigurationPropertySchema(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)
	if value != string(client.HCL) && value != string(client.Text) && value != string(client.JSON) {
		errs = append(errs, fmt.Errorf("%q can be either \"HCL\", \"JSON\" or empty, got: %q", key, value))
	}
	return
}
