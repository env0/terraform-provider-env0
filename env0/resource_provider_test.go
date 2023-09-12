package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitProviderResource(t *testing.T) {
	resourceType := "env0_provider"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	provider := client.Provider{
		Id:          uuid.NewString(),
		Type:        "aws",
		Description: "des",
	}

	updatedProvider := client.Provider{
		Id:          provider.Id,
		Type:        provider.Type,
		Description: "des-updated",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        provider.Type,
						"description": provider.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", provider.Id),
						resource.TestCheckResourceAttr(accessor, "type", provider.Type),
						resource.TestCheckResourceAttr(accessor, "description", provider.Description),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        updatedProvider.Type,
						"description": updatedProvider.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedProvider.Id),
						resource.TestCheckResourceAttr(accessor, "type", updatedProvider.Type),
						resource.TestCheckResourceAttr(accessor, "description", updatedProvider.Description),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProviderCreate(client.ProviderCreatePayload{
					Type:        provider.Type,
					Description: provider.Description,
				}).Times(1).Return(&provider, nil),
				mock.EXPECT().Provider(provider.Id).Times(2).Return(&provider, nil),
				mock.EXPECT().ProviderUpdate(updatedProvider.Id, client.ProviderUpdatePayload{
					Description: updatedProvider.Description,
				}).Times(1).Return(&updatedProvider, nil),
				mock.EXPECT().Provider(updatedProvider.Id).Times(1).Return(&updatedProvider, nil),
				mock.EXPECT().ProviderDelete(updatedProvider.Id).Times(1),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        provider.Type,
						"description": provider.Description,
					}),
					ExpectError: regexp.MustCompile("could not create provider: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProviderCreate(client.ProviderCreatePayload{
				Type:        provider.Type,
				Description: provider.Description,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Create Failure - Invalid Type", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        "-" + provider.Type,
						"description": provider.Description,
					}),
					ExpectError: regexp.MustCompile("must match pattern"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        provider.Type,
						"description": provider.Description,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        updatedProvider.Type,
						"description": updatedProvider.Description,
					}),
					ExpectError: regexp.MustCompile("could not update provider: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProviderCreate(client.ProviderCreatePayload{
					Type:        provider.Type,
					Description: provider.Description,
				}).Times(1).Return(&provider, nil),
				mock.EXPECT().Provider(provider.Id).Times(2).Return(&provider, nil),
				mock.EXPECT().ProviderUpdate(updatedProvider.Id, client.ProviderUpdatePayload{
					Description: updatedProvider.Description,
				}).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().ProviderDelete(updatedProvider.Id).Times(1),
			)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        provider.Type,
						"description": provider.Description,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     provider.Type,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProviderCreate(client.ProviderCreatePayload{
					Type:        provider.Type,
					Description: provider.Description,
				}).Times(1).Return(&provider, nil),
				mock.EXPECT().Provider(provider.Id).Times(1).Return(&provider, nil),
				mock.EXPECT().Providers().Times(1).Return([]client.Provider{provider}, nil),
				mock.EXPECT().Provider(provider.Id).Times(1).Return(&provider, nil),
				mock.EXPECT().ProviderDelete(provider.Id).Times(1),
			)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        provider.Type,
						"description": provider.Description,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     provider.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProviderCreate(client.ProviderCreatePayload{
					Type:        provider.Type,
					Description: provider.Description,
				}).Times(1).Return(&provider, nil),
				mock.EXPECT().Provider(provider.Id).Times(3).Return(&provider, nil),
				mock.EXPECT().ProviderDelete(provider.Id).Times(1),
			)
		})
	})

	t.Run("Import By Id Not Found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"type":        provider.Type,
						"description": provider.Description,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     provider.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProviderCreate(client.ProviderCreatePayload{
					Type:        provider.Type,
					Description: provider.Description,
				}).Times(1).Return(&provider, nil),
				mock.EXPECT().Provider(provider.Id).Times(1).Return(&provider, nil),
				mock.EXPECT().Provider(provider.Id).Times(1).Return(nil, &client.NotFoundError{}),
				mock.EXPECT().ProviderDelete(provider.Id).Times(1),
			)
		})
	})
}
