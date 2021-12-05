package env0

import (
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestUnitEnvironmentResource(t *testing.T) {
	resourceType := "env0_environment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	templateId := "template-id"

	environment := client.Environment{
		Id:        "id0",
		Name:      "my-environment",
		ProjectId: "project-id",
		LatestDeploymentLog: client.DeploymentLog{
			BlueprintId: templateId,
		},
	}

	updatedEnvironment := client.Environment{
		Id:        environment.Id,
		Name:      "my-updated-environment-name",
		ProjectId: "project-id",
		LatestDeploymentLog: client.DeploymentLog{
			BlueprintId: templateId,
		},
	}

	createEnvironmentResourceConfig := func(environment client.Environment) string {
		return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"name":        environment.Name,
			"project_id":  environment.ProjectId,
			"template_id": environment.LatestDeploymentLog.BlueprintId,
		})
	}

	t.Run("Success in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createEnvironmentResourceConfig(environment),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", templateId),
					),
				},
				{
					Config: createEnvironmentResourceConfig(updatedEnvironment),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", templateId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
				Name: updatedEnvironment.Name,
			}).Times(1).Return(updatedEnvironment, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	// TODO: test deploy with variables

	t.Run("Update to: template id, revision, repository, configuration should trigger a deployment", func(t *testing.T) {
		environment := client.Environment{
			Id:        "id0",
			Name:      "my-environment",
			ProjectId: "project-id",
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId:         "template-id",
				BlueprintRepository: "repository",
				BlueprintRevision:   "revision",
			},
		}
		updatedEnvironment := client.Environment{
			Id:        updatedEnvironment.Id,
			Name:      environment.Name,
			ProjectId: environment.ProjectId,
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId:         environment.LatestDeploymentLog.BlueprintId,
				BlueprintRepository: "updated repository",
				BlueprintRevision:   "updated revision",
			},
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        environment.Name,
						"project_id":  environment.ProjectId,
						"template_id": environment.LatestDeploymentLog.BlueprintId,
						"repository":  environment.LatestDeploymentLog.BlueprintRepository,
						"revision":    environment.LatestDeploymentLog.BlueprintRevision,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "repository", environment.LatestDeploymentLog.BlueprintRepository),
						resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedEnvironment.Name,
						"project_id":  updatedEnvironment.ProjectId,
						"template_id": updatedEnvironment.LatestDeploymentLog.BlueprintId,
						"repository":  updatedEnvironment.LatestDeploymentLog.BlueprintRepository,
						"revision":    updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", updatedEnvironment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "repository", updatedEnvironment.LatestDeploymentLog.BlueprintRepository),
						resource.TestCheckResourceAttr(accessor, "revision", updatedEnvironment.LatestDeploymentLog.BlueprintRevision),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId:         environment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision:   environment.LatestDeploymentLog.BlueprintRevision,
					BlueprintRepository: environment.LatestDeploymentLog.BlueprintRepository,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentDeploy(environment.Id, client.DeployRequest{
				BlueprintId:         updatedEnvironment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision:   updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
				BlueprintRepository: updatedEnvironment.LatestDeploymentLog.BlueprintRepository,
			}).Times(1).Return(client.EnvironmentDeployResponse{
				Id: "deployment-id",
			}, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("Failure in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createEnvironmentResourceConfig(environment),
					ExpectError: regexp.MustCompile("could not create environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(client.Environment{}, errors.New("error"))
		})

	})

	t.Run("Failure in update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createEnvironmentResourceConfig(environment),
				},
				{
					Config:      createEnvironmentResourceConfig(updatedEnvironment),
					ExpectError: regexp.MustCompile("could not update environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
				Name: updatedEnvironment.Name,
			}).Times(1).Return(client.Environment{}, errors.New("error"))
			mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil) // 1 after create, 1 before update
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})

	})

	t.Run("Failure in read", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createEnvironmentResourceConfig(environment),
					ExpectError: regexp.MustCompile("could not get environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().Environment(gomock.Any()).Return(client.Environment{}, errors.New("error"))
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})

	})
}
