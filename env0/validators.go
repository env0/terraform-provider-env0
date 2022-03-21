package env0

import (
	"fmt"
	"regexp"

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
	s := i.(string)
	if len(s) == 0 {
		return diag.Errorf("may not be empty")
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
