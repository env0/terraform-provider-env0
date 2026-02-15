package env0

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestEnvironmentDataSource(t *testing.T) {
	template := client.Template{
		Id:                   "template-id",
		TokenId:              "tokenId",
		GithubInstallationId: 100,
		BitbucketClientKey:   "bitbucket",
	}

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
			BlueprintId:       template.Id,
			BlueprintRevision: "revision",
			Output:            []byte(`{"a": "b"}`),
			WorkflowFile: &client.WorkflowFile{
				Environments: map[string]client.WorkflowSubEnvironment{
					"db": {
						EnvironmentId: "id_db",
					},
					"compute": {
						EnvironmentId: "id_compute",
					},
				},
			},
		},
	}

	environmentWithSameName := client.Environment{
		Id:        "other-id",
		Name:      environment.Name,
		ProjectId: "other-project-id",
	}

	archivedEnvironment := client.Environment{
		Id:         "id-archived",
		Name:       environment.Name,
		IsArchived: new(true),
	}

	environmentFieldsByName := map[string]any{"name": environment.Name}
	environmentFieldsByNameWithExclude := map[string]any{"name": environment.Name, "exclude_archived": "true"}
	environmentFieldByNameWithProjectId := map[string]any{"name": environment.Name, "project_id": environment.ProjectId}
	environmentFieldsById := map[string]any{"id": environment.Id}

	resourceType := "env0_environment"
	resourceName := "test_environment"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", strconv.FormatBool(!*environment.RequiresApproval)),
						resource.TestCheckResourceAttr(accessor, "run_plan_on_pull_requests", strconv.FormatBool(*environment.PullRequestPlanDeployments)),
						resource.TestCheckResourceAttr(accessor, "auto_deploy_on_path_changes_only", strconv.FormatBool(*environment.AutoDeployOnPathChangesOnly)),
						resource.TestCheckResourceAttr(accessor, "deploy_on_push", strconv.FormatBool(*environment.ContinuousDeployment)),
						resource.TestCheckResourceAttr(accessor, "status", environment.Status),
						resource.TestCheckResourceAttr(accessor, "deployment_id", environment.LatestDeploymentLogId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
						resource.TestCheckResourceAttr(accessor, "output", string(environment.LatestDeploymentLog.Output)),
						resource.TestCheckResourceAttr(accessor, "token_id", template.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(template.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "bitbucket_client_key", template.BitbucketClientKey),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.alias", "compute"),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.id", "id_compute"),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.1.alias", "db"),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.1.id", "id_db"),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockGetEnvironmentCall := func(env client.Environment, tem client.Template) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Environment(environment.Id).AnyTimes().Return(env, nil)
			mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).AnyTimes().Return(tem, nil)
		}
	}

	mockListEnvironmentsCall := func(returnValue []client.Environment, tem *client.Template) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentsByName(environment.Name).AnyTimes().Return(returnValue, nil)

			if tem != nil {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).AnyTimes().Return(*tem, nil)
			}
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(environmentFieldsById),
			mockGetEnvironmentCall(environment, template),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(environmentFieldsByName),
			mockListEnvironmentsCall([]client.Environment{environment}, &template),
		)
	})

	t.Run("By Name with Archived", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(environmentFieldsByNameWithExclude),
			mockListEnvironmentsCall([]client.Environment{environment, archivedEnvironment}, &template),
		)
	})

	t.Run("By Name with Different Project Id", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(environmentFieldByNameWithProjectId),
			mockListEnvironmentsCall([]client.Environment{environment, environmentWithSameName}, &template),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]any{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one environment exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(environmentFieldsByName, "Found multiple environments for name"),
			mockListEnvironmentsCall([]client.Environment{environment, environment}, nil),
		)
	})

	t.Run("Throw error when by name and more than one environment exists (archived use-case)", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(environmentFieldsByName, "Found multiple environments for name"),
			mockListEnvironmentsCall([]client.Environment{environment, archivedEnvironment}, nil),
		)
	})

	t.Run("Throw error when by name and no environments found at all", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(environmentFieldsByName, "Could not find an env0 environment with name"),
			mockListEnvironmentsCall([]client.Environment{}, nil),
		)
	})
}
