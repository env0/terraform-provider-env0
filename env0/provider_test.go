package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/utils"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func runUnitTest(t *testing.T, testCase resource.TestCase, mockFunc func(mockFunc *client.MockApiClientInterface)) {
	testReporter := utils.TestReporter{T: t}

	ctrl = gomock.NewController(&testReporter)

	apiClientMock = client.NewMockApiClientInterface(ctrl)
	mockFunc(apiClientMock)

	testCase.ProviderFactories = testUnitProviders
	testCase.PreventPostDestroyRefresh = true
	resource.UnitTest(&testReporter, testCase)

	ctrl.Finish()
}
