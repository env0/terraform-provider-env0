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

func TestUnitGitTokenResource(t *testing.T) {
	resourceType := "env0_git_token"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	gitToken := client.GitToken{
		Id:             uuid.NewString(),
		Name:           "name",
		Value:          "value",
		OrganizationId: "org",
	}

	updatedGitToken := client.GitToken{
		Id:             "id2",
		Name:           "name",
		Value:          "value2",
		OrganizationId: "org",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  gitToken.Name,
						"value": gitToken.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", gitToken.Id),
						resource.TestCheckResourceAttr(accessor, "name", gitToken.Name),
						resource.TestCheckResourceAttr(accessor, "value", gitToken.Value),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  updatedGitToken.Name,
						"value": updatedGitToken.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedGitToken.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedGitToken.Name),
						resource.TestCheckResourceAttr(accessor, "value", updatedGitToken.Value),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().GitTokenCreate(client.GitTokenCreatePayload{
					Name:  gitToken.Name,
					Value: gitToken.Value,
				}).Times(1).Return(&gitToken, nil),
				mock.EXPECT().GitToken(gitToken.Id).Times(2).Return(&gitToken, nil),
				mock.EXPECT().GitTokenDelete(gitToken.Id).Times(1),
				mock.EXPECT().GitTokenCreate(client.GitTokenCreatePayload{
					Name:  updatedGitToken.Name,
					Value: updatedGitToken.Value,
				}).Times(1).Return(&updatedGitToken, nil),
				mock.EXPECT().GitToken(updatedGitToken.Id).Times(1).Return(&updatedGitToken, nil),
				mock.EXPECT().GitTokenDelete(updatedGitToken.Id).Times(1),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  gitToken.Name,
						"value": gitToken.Value,
					}),
					ExpectError: regexp.MustCompile("could not create git token: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GitTokenCreate(client.GitTokenCreatePayload{
				Name:  gitToken.Name,
				Value: gitToken.Value,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  gitToken.Name,
						"value": gitToken.Value,
					}),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           gitToken.Name,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"value"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GitTokenCreate(client.GitTokenCreatePayload{
				Name:  gitToken.Name,
				Value: gitToken.Value,
			}).Times(1).Return(&gitToken, nil)
			mock.EXPECT().GitTokens().Times(1).Return([]client.GitToken{gitToken}, nil)
			mock.EXPECT().GitToken(gitToken.Id).Times(2).Return(&gitToken, nil)
			mock.EXPECT().GitTokenDelete(gitToken.Id).Times(1)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  gitToken.Name,
						"value": gitToken.Value,
					}),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           gitToken.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"value"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GitTokenCreate(client.GitTokenCreatePayload{
				Name:  gitToken.Name,
				Value: gitToken.Value,
			}).Times(1).Return(&gitToken, nil)
			mock.EXPECT().GitToken(gitToken.Id).Times(3).Return(&gitToken, nil)
			mock.EXPECT().GitTokenDelete(gitToken.Id).Times(1)
		})
	})

}
