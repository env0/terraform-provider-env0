package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestGpgKeyDataSource(t *testing.T) {
	gpgKey := client.GpgKey{
		Id:      "id0",
		Name:    "name0",
		KeyId:   "ABCDABCDABCD1113",
		Content: "content1",
	}

	otherGpgKey := client.GpgKey{
		Id:      "id1",
		Name:    "name1",
		KeyId:   "ABCDABCDABCD1112",
		Content: "content2",
	}

	gpgKeyFieldsByName := map[string]interface{}{"name": gpgKey.Name}
	gpgKeyFieldsById := map[string]interface{}{"id": gpgKey.Id}

	resourceType := "env0_gpg_key"
	resourceName := "test_gpg_key"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", gpgKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", gpgKey.Name),
						resource.TestCheckResourceAttr(accessor, "key_id", gpgKey.KeyId),
						resource.TestCheckResourceAttr(accessor, "content", gpgKey.Content),
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

	mockListGpgKeysCall := func(returnValue []client.GpgKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GpgKeys().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(gpgKeyFieldsById),
			mockListGpgKeysCall([]client.GpgKey{gpgKey, otherGpgKey}),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(gpgKeyFieldsByName),
			mockListGpgKeysCall([]client.GpgKey{gpgKey, otherGpgKey}),
		)
	})

	t.Run("Throw error when by name and more than one gpg key exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(gpgKeyFieldsByName, "found multiple gpg keys"),
			mockListGpgKeysCall([]client.GpgKey{gpgKey, otherGpgKey, gpgKey}),
		)
	})

	t.Run("Throw error when by id and no gpg key found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(gpgKeyFieldsById, "could not read gpg key: not found"),
			mockListGpgKeysCall([]client.GpgKey{otherGpgKey}),
		)
	})
}
