package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestApiKeyDataSource(t *testing.T) {
	apiKey := client.ApiKey{
		Id:           "id0",
		Name:         "name0",
		ApiKeyId:     "keyid0",
		ApiKeySecret: "secret0",
	}

	otherApiKey := client.ApiKey{
		Id:           "id1",
		Name:         "name1",
		ApiKeyId:     "keyid1",
		ApiKeySecret: "secret1",
	}

	apiKeyFieldsByName := map[string]interface{}{"name": apiKey.Name}
	apiKeyFieldsById := map[string]interface{}{"id": apiKey.Id}

	resourceType := "env0_api_key"
	resourceName := "test_api_key"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", apiKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", apiKey.Name),
						resource.TestCheckResourceAttr(accessor, "api_key_id", apiKey.ApiKeyId),
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

	mockListApiKeysCall := func(returnValue []client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApiKeys().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(apiKeyFieldsById),
			mockListApiKeysCall([]client.ApiKey{apiKey, otherApiKey}),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(apiKeyFieldsByName),
			mockListApiKeysCall([]client.ApiKey{apiKey, otherApiKey}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one api key exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(apiKeyFieldsByName, "found multiple api keys"),
			mockListApiKeysCall([]client.ApiKey{apiKey, otherApiKey, apiKey}),
		)
	})

	t.Run("Throw error when by name and no api key found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(apiKeyFieldsByName, "not found"),
			mockListApiKeysCall([]client.ApiKey{otherApiKey}),
		)
	})

	t.Run("Throw error when by id and no api key found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(apiKeyFieldsById, fmt.Sprintf("id %s not found", apiKey.Id)),
			mockListApiKeysCall([]client.ApiKey{otherApiKey}),
		)
	})
}
