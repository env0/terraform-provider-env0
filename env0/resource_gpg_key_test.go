package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitGpgKeyyResource(t *testing.T) {
	resourceType := "env0_gpg_key"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	gpgKey := client.GpgKey{
		Id:      uuid.NewString(),
		Name:    "name",
		KeyId:   "ABCD0123ABCD0123",
		Content: "content",
	}

	updatedGpgKey := client.GpgKey{
		Id:      uuid.NewString(),
		Name:    "name2",
		KeyId:   "11110123ABCD4567",
		Content: "content2",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":    gpgKey.Name,
						"key_id":  gpgKey.KeyId,
						"content": gpgKey.Content,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", gpgKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", gpgKey.Name),
						resource.TestCheckResourceAttr(accessor, "key_id", gpgKey.KeyId),
						resource.TestCheckResourceAttr(accessor, "content", gpgKey.Content),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":    updatedGpgKey.Name,
						"key_id":  updatedGpgKey.KeyId,
						"content": updatedGpgKey.Content,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedGpgKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedGpgKey.Name),
						resource.TestCheckResourceAttr(accessor, "key_id", updatedGpgKey.KeyId),
						resource.TestCheckResourceAttr(accessor, "content", updatedGpgKey.Content),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().GpgKeyCreate(&client.GpgKeyCreatePayload{
					Name:    gpgKey.Name,
					KeyId:   gpgKey.KeyId,
					Content: gpgKey.Content,
				}).Times(1).Return(&gpgKey, nil),
				mock.EXPECT().GpgKeys().Times(2).Return([]client.GpgKey{gpgKey}, nil),
				mock.EXPECT().GpgKeyDelete(gpgKey.Id).Times(1),
				mock.EXPECT().GpgKeyCreate(&client.GpgKeyCreatePayload{
					Name:    updatedGpgKey.Name,
					KeyId:   updatedGpgKey.KeyId,
					Content: updatedGpgKey.Content,
				}).Times(1).Return(&updatedGpgKey, nil),
				mock.EXPECT().GpgKeys().Times(1).Return([]client.GpgKey{updatedGpgKey}, nil),
				mock.EXPECT().GpgKeyDelete(updatedGpgKey.Id).Times(1),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":    gpgKey.Name,
						"key_id":  gpgKey.KeyId,
						"content": gpgKey.Content,
					}),
					ExpectError: regexp.MustCompile("could not create gpg key: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GpgKeyCreate(&client.GpgKeyCreatePayload{
				Name:    gpgKey.Name,
				KeyId:   gpgKey.KeyId,
				Content: gpgKey.Content,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":    gpgKey.Name,
						"key_id":  gpgKey.KeyId,
						"content": gpgKey.Content,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     gpgKey.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GpgKeyCreate(&client.GpgKeyCreatePayload{
				Name:    gpgKey.Name,
				KeyId:   gpgKey.KeyId,
				Content: gpgKey.Content,
			}).Times(1).Return(&gpgKey, nil)
			mock.EXPECT().GpgKeys().Times(3).Return([]client.GpgKey{gpgKey}, nil)
			mock.EXPECT().GpgKeyDelete(gpgKey.Id).Times(1)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":    gpgKey.Name,
						"key_id":  gpgKey.KeyId,
						"content": gpgKey.Content,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     gpgKey.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GpgKeyCreate(&client.GpgKeyCreatePayload{
				Name:    gpgKey.Name,
				KeyId:   gpgKey.KeyId,
				Content: gpgKey.Content,
			}).Times(1).Return(&gpgKey, nil)
			mock.EXPECT().GpgKeys().Times(3).Return([]client.GpgKey{gpgKey}, nil)
			mock.EXPECT().GpgKeyDelete(gpgKey.Id).Times(1)
		})
	})

	t.Run("Import By Id Not Found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":    updatedGpgKey.Name,
						"key_id":  updatedGpgKey.KeyId,
						"content": updatedGpgKey.Content,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     gpgKey.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GpgKeyCreate(&client.GpgKeyCreatePayload{
				Name:    updatedGpgKey.Name,
				KeyId:   updatedGpgKey.KeyId,
				Content: updatedGpgKey.Content,
			}).Times(1).Return(&updatedGpgKey, nil)
			mock.EXPECT().GpgKeys().Times(2).Return([]client.GpgKey{updatedGpgKey}, nil)
			mock.EXPECT().GpgKeyDelete(updatedGpgKey.Id).Times(1)
		})
	})
}
