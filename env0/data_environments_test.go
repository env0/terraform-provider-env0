package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestEnvironmentsDataSource(t *testing.T) {
	env1 := client.Environment{
		Id:         "env1",
		Name:       "Environment 1",
		ProjectId:  "proj1",
		IsArchived: new(false),
	}

	env2 := client.Environment{
		Id:         "env2",
		Name:       "Environment 2",
		ProjectId:  "proj1",
		IsArchived: new(false),
	}

	env3Archived := client.Environment{
		Id:         "env3",
		Name:       "Environment 3",
		ProjectId:  "proj1",
		IsArchived: new(true),
	}

	resourceType := "env0_environments"
	resourceName := "test_environments"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(includeArchived bool) resource.TestCase {
		fields := map[string]any{
			"project_id": "proj1",
		}
		if includeArchived {
			fields["include_archived_environments"] = "true"
		}

		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, fields),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environments.0.id", env1.Id),
						resource.TestCheckResourceAttr(accessor, "environments.0.name", env1.Name),
						resource.TestCheckResourceAttr(accessor, "environments.0.project_id", env1.ProjectId),
						resource.TestCheckResourceAttr(accessor, "environments.1.id", env2.Id),
						resource.TestCheckResourceAttr(accessor, "environments.1.name", env2.Name),
						resource.TestCheckResourceAttr(accessor, "environments.1.project_id", env2.ProjectId),
						resource.TestCheckResourceAttr(accessor, "environments.#", func() string {
							if includeArchived {
								return "3"
							} else {
								return "2"
							}
						}()),
					),
				},
			},
		}
	}

	getErrorTestCase := func(expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]any{"project_id": "proj1"}),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockProjectEnvsCall := func(returnValue []client.Environment) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectEnvironments("proj1").AnyTimes().Return(returnValue, nil)
		}
	}

	mockProjectEnvsCallFailed := func(err string) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectEnvironments("proj1").AnyTimes().Return([]client.Environment{}, errors.New(err))
		}
	}

	t.Run("get project environments excluding archived", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(false),
			mockProjectEnvsCall([]client.Environment{env1, env2, env3Archived}),
		)
	})

	t.Run("get project environments including archived", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(true),
			mockProjectEnvsCall([]client.Environment{env1, env2, env3Archived}),
		)
	})

	t.Run("Error when API call fails", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase("failed to get list of environments: error"),
			mockProjectEnvsCallFailed("error"),
		)
	})
}
