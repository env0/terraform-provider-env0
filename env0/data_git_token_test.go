package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestGitTokenDataSource(t *testing.T) {
	gitToken := client.GitToken{
		Id:   "id0",
		Name: "name0",
	}

	otherGitToken := client.GitToken{
		Id:   "id1",
		Name: "name1",
	}

	gitTokenFieldsByName := map[string]interface{}{"name": gitToken.Name}
	gitTokenFieldsById := map[string]interface{}{"id": gitToken.Id}

	resourceType := "env0_git_token"
	resourceName := "test_git_token"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", gitToken.Id),
						resource.TestCheckResourceAttr(accessor, "name", gitToken.Name),
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

	mockListGitTokensCall := func(returnValue []client.GitToken) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GitTokens().AnyTimes().Return(returnValue, nil)
		}
	}

	mockGitTokenCall := func(returnValue *client.GitToken) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().GitToken(gitToken.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(gitTokenFieldsById),
			mockGitTokenCall(&gitToken),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(gitTokenFieldsByName),
			mockListGitTokensCall([]client.GitToken{gitToken, otherGitToken}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one git token exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(gitTokenFieldsByName, "found multiple git tokens"),
			mockListGitTokensCall([]client.GitToken{gitToken, otherGitToken, gitToken}),
		)
	})

	t.Run("Throw error when by name and no git token found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(gitTokenFieldsByName, "not found"),
			mockListGitTokensCall([]client.GitToken{otherGitToken}),
		)
	})

	t.Run("Throw error when by id and no git token found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(gitTokenFieldsById, fmt.Sprintf("id %s not found", gitToken.Id)),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().GitToken(gitToken.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404))
			},
		)
	})
}
