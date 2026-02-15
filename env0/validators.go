package env0

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/adhocore/gronx"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ValidateConfigurationPropertySchema(val any, key string) (warns []string, errs []error) {
	value := val.(string)
	if value != string(client.HCL) && value != string(client.Text) && value != string(client.JSON) && value != string(client.ENVIRONMENT_OUTPUT) {
		errs = append(errs, fmt.Errorf("%q can be either \"HCL\", \"JSON\", \"ENVIRONMENT_OUTPUT\" or empty, got: %q", key, value))
	}

	return
}

func ValidateCronExpression(i any, path cty.Path) diag.Diagnostics {
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

func ValidateNotEmptyString(i any, path cty.Path) diag.Diagnostics {
	s := strings.TrimSpace(i.(string))
	if len(s) == 0 {
		return diag.Errorf("may not be empty")
	}

	return nil
}

func ValidateRetries(i any, path cty.Path) diag.Diagnostics {
	retries := i.(int)
	if retries < 1 || retries > 3 {
		return diag.Errorf("retries amount must be between 1 and 3")
	}

	return nil
}

func NewRegexValidator(r string) schema.SchemaValidateDiagFunc {
	cr := regexp.MustCompile(r)

	return func(i any, p cty.Path) diag.Diagnostics {
		if !cr.MatchString(i.(string)) {
			return diag.Errorf("must match pattern %v", r)
		}

		return nil
	}
}

func NewStringInValidator(allowedValues []string) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		value := i.(string)
		for _, allowedValue := range allowedValues {
			if value == allowedValue {
				return nil
			}
		}

		return diag.Errorf("'%s' must be one of: %s", value, strings.Join(allowedValues, ", "))
	}
}

func NewIntInValidator(allowedValues []int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		value := i.(int)
		for _, allowedValue := range allowedValues {
			if value == allowedValue {
				return nil
			}
		}

		return diag.Errorf("must be one of: %s", fmt.Sprint(allowedValues))
	}
}

func NewGreaterThanValidator(greaterThan int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		value := i.(int)
		if value <= greaterThan {
			return diag.Errorf("%d must be greater than %d", value, greaterThan)
		}

		return nil
	}
}

func NewOpenTofuVersionValidator() schema.SchemaValidateDiagFunc {
	return NewRegexValidator(`(?:^[0-9]\.[0-9]{1,2}\.[0-9]{1,2}(?:-.+)?$)|^RESOLVE_FROM_CODE$|^latest$`)
}

func ValidateTtl(i any, path cty.Path) diag.Diagnostics {
	ttl := i.(string)

	_, err := ttlToDuration(&ttl)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func NewRoleValidator(supportedBuiltInRoles []string) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		role := i.(string)

		if role == "" {
			return diag.Errorf("may not be empty")
		}

		if client.IsCustomRole(role) {
			// Custom role.
			return nil
		}

		// Built-in role. Verify it's in the supported list.
		for _, supportedRole := range supportedBuiltInRoles {
			if role == supportedRole {
				// supported.
				return nil
			}
		}

		// not supported.
		return diag.Errorf("the following built-in role '%s' is not supported for this resource, must be one of %s", role, "["+strings.Join(supportedBuiltInRoles, ",")+"]")
	}
}

func ValidateUrl(i any, path cty.Path) diag.Diagnostics {
	v := i.(string)

	_, err := url.ParseRequestURI(v)
	if err != nil {
		return diag.Errorf("must be a valid URL: %v", err)
	}

	return nil
}
