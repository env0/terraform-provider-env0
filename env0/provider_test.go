package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	. "github.com/onsi/ginkgo"
	"testing"
)

var (
	apiClientMock *client.MockApiClientInterface
	ctrl          *gomock.Controller
)

var testUnitProviders = map[string]func() (*schema.Provider, error){
	"env0": func() (*schema.Provider, error) {
		provider := Provider()
		provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
			return apiClientMock, nil
		}
		return provider, nil
	},
}

var _ = BeforeEach(func() {
	ctrl = gomock.NewController(utils.RecoveringGinkgoT())
	apiClientMock = client.NewMockApiClientInterface(ctrl)
})

var _ = AfterEach(func() {
	ctrl.Finish()
})

func TestProvider(t *testing.T) {
	RunSpecs(t, "Provider Tests")
}
