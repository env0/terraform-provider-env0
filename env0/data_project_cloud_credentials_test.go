package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProjectCloudCredentialsDataSource(t *testing.T) {
	resourceType := "env0_project_cloud_credentials"
	resourceName := "test_project_cloud_credentials"
	accessor := dataSourceAccessor(resourceType, resourceName)

	credentialIds := []string{"id1", "id2"}
	projectId := "project_id"

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{
							"project_id": projectId,
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "ids.0", credentialIds[0]),
							resource.TestCheckResourceAttr(accessor, "ids.1", credentialIds[1]),
						),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().CloudCredentialIdsInProject(projectId).AnyTimes().Return(credentialIds, nil)
			},
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{
							"project_id": projectId,
						}),
						ExpectError: regexp.MustCompile("could not get cloud credentials associated with project.*error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().CloudCredentialIdsInProject(projectId).AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
