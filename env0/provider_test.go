package env0

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testUnitProvider *schema.Provider

var testAccProviders = map[string]func() (*schema.Provider, error){
	"env0": func() (*schema.Provider, error) {
		if testAccProvider == nil {
			testAccProvider = Provider()
		}
		return testAccProvider, nil
	},
}

var testUnitProviders = map[string]func() (*schema.Provider, error){
	"env0": func() (*schema.Provider, error) {
		if testUnitProvider == nil {
			provider := Provider()
			provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
				return nil, nil
			}
			testUnitProvider = Provider()
		}
		return testUnitProvider, nil
	},
}
