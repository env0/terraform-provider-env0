package env0

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/adhocore/gronx"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ValidateConfigurationPropertySchema(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)
	if value != string(client.HCL) && value != string(client.Text) && value != string(client.JSON) {
		errs = append(errs, fmt.Errorf("%q can be either \"HCL\", \"JSON\" or empty, got: %q", key, value))
	}
	return
}

func ValidateCronExpression(i interface{}, path cty.Path) diag.Diagnostics {
	expr := i.(string)
	parser := gronx.New()
	isValid := parser.IsValid(expr)

	if !isValid {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid cron expression",
				AttributePath: path,
			}}
	}

	return nil
}

func ValidateNotEmptyString(i interface{}, path cty.Path) diag.Diagnostics {
	s := strings.TrimSpace(i.(string))
	if len(s) == 0 {
		return diag.Errorf("may not be empty")
	}

	return nil
}

func ValidateRetries(i interface{}, path cty.Path) diag.Diagnostics {
	retries := i.(int)
	if retries < 1 || retries > 3 {
		return diag.Errorf("retries amount must be between 1 and 3")
	}

	return nil
}

func ValidateRole(i interface{}, path cty.Path) diag.Diagnostics {
	role := client.ProjectRole(i.(string))
	if role == "" ||
		role != client.Admin &&
			role != client.Deployer &&
			role != client.Viewer &&
			role != client.Planner {
		return diag.Errorf("must be one of [Admin, Deployer, Viewer, Planner], got: %v", role)
	}
	return nil
}

func NewRegexValidator(r string) schema.SchemaValidateDiagFunc {
	cr := regexp.MustCompile(r)

	return func(i interface{}, p cty.Path) diag.Diagnostics {
		if !cr.MatchString(i.(string)) {
			return diag.Errorf("must match pattern %v", r)
		}
		return nil
	}
}

func NewStringInValidator(allowedValues []string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		value := i.(string)
		for _, allowedValue := range allowedValues {
			if value == allowedValue {
				return nil
			}
		}

		return diag.Errorf("'%s' must be one of: %s", value, strings.Join(allowedValues, ", "))
	}
}

func NewGreaterThanValidator(greaterThan int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		value := i.(int)
		if value <= greaterThan {
			return diag.Errorf("%d must be greater than %d", value, greaterThan)
		}

		return nil
	}
}
