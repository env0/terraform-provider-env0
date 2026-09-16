package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitEnvironmentResourceWaitForDeployment(t *testing.T) {
	t.Parallel()

	resourceType := "env0_environment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	templateId := "template-id"
	deploymentLogId := "deployment-log-id"

	environment := client.Environment{
		Id:                    uuid.New().String(),
		Name:                  "my-environment",
		ProjectId:             "project-id",
		LatestDeploymentLogId: deploymentLogId,
		LatestDeploymentLog: client.DeploymentLog{
			Id:          deploymentLogId,
			BlueprintId: templateId,
		},
	}

	deployedEnvironment := environment
	deployedEnvironment.LatestDeploymentLog.Output = []byte(`{"a":"b"}`)

	template := client.Template{ProjectId: environment.ProjectId}

	environmentCreate := client.EnvironmentCreate{
		Name:          environment.Name,
		ProjectId:     environment.ProjectId,
		DeployRequest: &client.DeployRequest{BlueprintId: templateId},
	}

	deploymentWithStatus := func(status string) *client.DeploymentLog {
		deployment := environment.LatestDeploymentLog
		deployment.Status = status

		return &deployment
	}

	config := func(fields map[string]any) string {
		base := map[string]any{
			"name":                environment.Name,
			"project_id":          environment.ProjectId,
			"template_id":         templateId,
			"force_destroy":       true,
			"wait_for_deployment": true,
		}

		for key, value := range fields {
			base[key] = value
		}

		return resourceConfigCreate(resourceType, resourceName, base)
	}

	// The teardown destroy only reaches EnvironmentDestroy when "force_destroy" is true in the state, so
	// expecting the call is how a lifted safeguard is asserted after a failed create.
	expectTeardownDestroy := func(mock *client.MockApiClientInterface) any {
		return mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1).Return(&client.EnvironmentDestroyResponse{Id: deploymentLogId}, nil)
	}

	expectRead := func(mock *client.MockApiClientInterface, env client.Environment) []any {
		return []any{
			mock.EXPECT().Environment(environment.Id).Times(1).Return(env, nil),
			mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
			mock.EXPECT().ConfigurationSetsAssignments("ENVIRONMENT", environment.Id).Times(1).Return(nil, nil),
		}
	}

	t.Run("create waits and populates output", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config(nil),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "deployment_id", deploymentLogId),
						resource.TestCheckResourceAttr(accessor, "output", `{"a":"b"}`),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			calls := []any{
				mock.EXPECT().Template(templateId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).Times(1).Return(deploymentWithStatus("IN_PROGRESS"), nil),
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).Times(1).Return(deploymentWithStatus("SUCCESS"), nil),
				mock.EXPECT().Environment(environment.Id).Times(1).Return(deployedEnvironment, nil),
			}
			calls = append(calls, expectRead(mock, deployedEnvironment)...)
			calls = append(calls, expectTeardownDestroy(mock))

			gomock.InOrder(calls...)
		})
	})

	t.Run("a deployment waiting for approval does not block", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config(nil),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "output", "null"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			calls := []any{
				mock.EXPECT().Template(templateId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).Times(1).Return(deploymentWithStatus("WAITING_FOR_USER"), nil),
			}
			calls = append(calls, expectRead(mock, environment)...)
			calls = append(calls, expectTeardownDestroy(mock))

			gomock.InOrder(calls...)
		})
	})

	t.Run("a failed deployment fails the apply and lifts force_destroy", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      config(map[string]any{"force_destroy": false}),
					ExpectError: regexp.MustCompile("did not succeed: failed to wait for environment deploy to complete, deployment status is: FAILURE"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(templateId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).Times(1).Return(deploymentWithStatus("FAILURE"), nil),
				expectTeardownDestroy(mock),
			)
		})
	})

	t.Run("a timed out deployment fails the apply and lifts force_destroy", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      config(map[string]any{"force_destroy": false}),
					ExpectError: regexp.MustCompile("did not succeed: timeout! last 'deploy' deployment status was 'IN_PROGRESS'"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(templateId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).AnyTimes().Return(deploymentWithStatus("IN_PROGRESS"), nil),
				expectTeardownDestroy(mock),
			)
		})
	})

	t.Run("update waits for the redeploy", func(t *testing.T) {
		updatedDeploymentLogId := "updated-deployment-log-id"

		redeployedEnvironment := environment
		redeployedEnvironment.LatestDeploymentLogId = updatedDeploymentLogId
		redeployedEnvironment.LatestDeploymentLog.BlueprintRevision = "v2"
		redeployedEnvironment.LatestDeploymentLog.Output = []byte(`{"a":"c"}`)

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config(nil),
					Check:  resource.TestCheckResourceAttr(accessor, "deployment_id", deploymentLogId),
				},
				{
					Config: config(map[string]any{"revision": "v2"}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "deployment_id", updatedDeploymentLogId),
						resource.TestCheckResourceAttr(accessor, "output", `{"a":"c"}`),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(templateId).AnyTimes().Return(template, nil)
			mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil)
			mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).AnyTimes().Return(client.ConfigurationChanges{}, nil)
			mock.EXPECT().ConfigurationSetsAssignments("ENVIRONMENT", environment.Id).AnyTimes().Return(nil, nil)

			// gomock hands out matching expectations in declaration order, so every read before the
			// redeploy sees the first deployment and every read after it sees the second.
			gomock.InOrder(
				mock.EXPECT().Environment(environment.Id).Times(3).Return(deployedEnvironment, nil),
				mock.EXPECT().Environment(environment.Id).AnyTimes().Return(redeployedEnvironment, nil),
			)

			gomock.InOrder(
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).Times(1).Return(deploymentWithStatus("SUCCESS"), nil),
				mock.EXPECT().EnvironmentDeploy(environment.Id, gomock.Any()).Times(1).Return(client.EnvironmentDeployResponse{Id: updatedDeploymentLogId}, nil),
				mock.EXPECT().EnvironmentDeploymentLog(updatedDeploymentLogId).Times(1).Return(deploymentWithStatus("IN_PROGRESS"), nil),
				mock.EXPECT().EnvironmentDeploymentLog(updatedDeploymentLogId).Times(1).Return(deploymentWithStatus("SUCCESS"), nil),
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1).Return(&client.EnvironmentDestroyResponse{Id: updatedDeploymentLogId}, nil)
		})
	})

	t.Run("a deployment superseded while waiting is reported", func(t *testing.T) {
		supersededEnvironment := environment
		supersededEnvironment.LatestDeploymentLogId = "another-deployment-log-id"

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config(nil),
					Check:  resource.TestCheckResourceAttr(accessor, "deployment_id", "another-deployment-log-id"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			calls := []any{
				mock.EXPECT().Template(templateId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
				mock.EXPECT().EnvironmentDeploymentLog(deploymentLogId).Times(1).Return(deploymentWithStatus("SUCCESS"), nil),
				mock.EXPECT().Environment(environment.Id).Times(1).Return(supersededEnvironment, nil),
			}
			calls = append(calls, expectRead(mock, supersededEnvironment)...)
			calls = append(calls, expectTeardownDestroy(mock))

			gomock.InOrder(calls...)
		})
	})

	t.Run("without the flag nothing is polled", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config(map[string]any{"wait_for_deployment": false}),
					Check:  resource.TestCheckResourceAttr(accessor, "output", "null"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			// No EnvironmentDeploymentLog expectation: gomock fails the test if the poller runs at all.
			calls := []any{
				mock.EXPECT().Template(templateId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
			}
			calls = append(calls, expectRead(mock, environment)...)
			calls = append(calls, expectTeardownDestroy(mock))

			gomock.InOrder(calls...)
		})
	})
}
