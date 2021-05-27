package env0

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var testAccProvider *schema.Provider
var testAccProviders = map[string]func() (*schema.Provider, error){
	"env0": func() (*schema.Provider, error) {
		if testAccProvider == nil {
			testAccProvider = Provider()
		}
		return testAccProvider, nil
	},
}

func TestMain(m *testing.M) {

}
