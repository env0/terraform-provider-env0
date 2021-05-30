package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var (
	apiClientMock *client.MockApiClientInterface
	ctrl          *gomock.Controller
)

func mockApiClient(t *testing.T) (*gomock.Controller, *client.MockApiClientInterface) {
	if ctrl == nil {
		ctrl = gomock.NewController(t)
	}

	if apiClientMock == nil {
		apiClientMock = client.NewMockApiClientInterface(ctrl)
	}

	return ctrl, apiClientMock
}

func testUnitProviders(ctrl *gomock.Controller) map[string]func() (*schema.Provider, error) {
	if apiClientMock == nil {
		apiClientMock = client.NewMockApiClientInterface(ctrl)
	}

	return map[string]func() (*schema.Provider, error){
		"env0": func() (*schema.Provider, error) {
			provider := Provider()
			provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
				return apiClientMock, nil
			}
			return provider, nil
		},
	}
}

func runUnitTest(t *testing.T, c resource.TestCase) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c.ProviderFactories = testUnitProviders(ctrl)
	resource.Test(t, c)
}
