package env0

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestGithubInstallationIdDataSource(t *testing.T) {
	mockToken := client.VscToken{
		Token: 12345,
	}

	mockRepositroy := "http://myrepo.com"

	resourceType := "env0_github_installation_id"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(repository string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"repository": repository,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(mockToken.Token)),
					),
				},
			},
		}
	}

	getErrorTestCase := func(repository string, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"repository": repository,
					}),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	t.Run("get by repository", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(mockRepositroy),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().VcsToken("github", mockRepositroy).Return(&mockToken, nil).AnyTimes()
			},
		)
	})

	t.Run("get by repository - failed", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(mockRepositroy, "failed to get github installation id: error"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().VcsToken("github", mockRepositroy).Return(nil, errors.New("error")).AnyTimes()
			},
		)
	})
}
