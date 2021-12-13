package env0

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestEnvironmentDataSource(t *testing.T) {
	boolean := true
	environment := client.Environment{
		Id:                          "id-0",
		Name:                        "my-environment-1",
		ProjectId:                   "project-id",
		Status:                      "status",
		LatestDeploymentLogId:       "latest-deployment-log-id",
		RequiresApproval:            &boolean,
		PullRequestPlanDeployments:  &boolean,
		AutoDeployOnPathChangesOnly: &boolean,
		ContinuousDeployment:        &boolean,
		LatestDeploymentLog: client.DeploymentLog{
			BlueprintId:       "blueprint-id",
			BlueprintRevision: "revision",
		},
	}

	otherEnvironment := client.Environment{
		Id:   "other-id",
		Name: "other-name",
	}

	environmentFieldsByName := map[string]interface{}{"name": environment.Name}
	environmentFieldsById := map[string]interface{}{"id": environment.Id}

	resourceType := "env0_environment"
	resourceName := "test_environment"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", strconv.FormatBool(!*environment.RequiresApproval)),
						resource.TestCheckResourceAttr(accessor, "run_plan_on_pull_requests", strconv.FormatBool(*environment.PullRequestPlanDeployments)),
						resource.TestCheckResourceAttr(accessor, "auto_deploy_on_path_changes_only", strconv.FormatBool(*environment.AutoDeployOnPathChangesOnly)),
						resource.TestCheckResourceAttr(accessor, "deploy_on_push", strconv.FormatBool(*environment.ContinuousDeployment)),
						resource.TestCheckResourceAttr(accessor, "status", environment.Status),
						resource.TestCheckResourceAttr(accessor, "latest_deployment_log_id", environment.LatestDeploymentLogId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
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

	mockGetEnvironmentCall := func(returnValue client.Environment) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Environment(environment.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListEnvironmentsCall := func(returnValue []client.Environment) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Environments().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(environmentFieldsById),
			mockGetEnvironmentCall(environment),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(environmentFieldsByName),
			mockListEnvironmentsCall([]client.Environment{environment, otherEnvironment}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one environment exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(environmentFieldsByName, "Found multiple environments for name"),
			mockListEnvironmentsCall([]client.Environment{environment, environment, otherEnvironment}),
		)
	})

	t.Run("Throw error when by name and no environments found at all", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(environmentFieldsByName, "Could not find an env0 environment with name"),
			mockListEnvironmentsCall([]client.Environment{}),
		)
	})

	t.Run("Throw error when by name and no environments found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(environmentFieldsByName, "Could not find an env0 environment with name"),
			mockListEnvironmentsCall([]client.Environment{otherEnvironment}),
		)
	})
}
