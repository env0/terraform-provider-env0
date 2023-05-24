package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProviderDataSource(t *testing.T) {
	provider := client.Provider{
		Id:          "id0",
		Type:        "type0",
		Description: "des0",
	}

	otherProvider := client.Provider{
		Id:          "id1",
		Type:        "type1",
		Description: "des1",
	}

	providerFieldByName := map[string]interface{}{"type": provider.Type}
	providerFieldById := map[string]interface{}{"id": provider.Id}

	resourceType := "env0_provider"
	resourceName := "test_provider"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", provider.Id),
						resource.TestCheckResourceAttr(accessor, "type", provider.Type),
						resource.TestCheckResourceAttr(accessor, "description", provider.Description),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]interface{}, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockListProvidersCall := func(returnValue []client.Provider) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Providers().AnyTimes().Return(returnValue, nil)
		}
	}

	mockProviderCall := func(returnValue *client.Provider) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Provider(provider.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(providerFieldById),
			mockProviderCall(&provider),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(providerFieldByName),
			mockListProvidersCall([]client.Provider{otherProvider, provider}),
		)
	})

	t.Run("Throw error when by name and more than one provider exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(providerFieldByName, "found multiple providers"),
			mockListProvidersCall([]client.Provider{provider, otherProvider, provider}),
		)
	})

	t.Run("Throw error when by id and no provider found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(providerFieldById, "could not read provider: id "+provider.Id+" not found"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Provider(provider.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404))
			},
		)
	})
}
