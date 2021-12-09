package env0

import (
	"errors"
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"strings"
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

	// TODO: test deploy with variables ( needs configuration changes read )

	t.Run("Update to: revision, configuration should trigger a deployment", func(t *testing.T) {
		environment := client.Environment{
			Id:        "id0",
			Name:      "my-environment",
			ProjectId: "project-id",
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId:       "template-id",
				BlueprintRevision: "revision",
			},
		}
		updatedEnvironment := client.Environment{
			Id:        updatedEnvironment.Id,
			Name:      environment.Name,
			ProjectId: environment.ProjectId,
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision: "updated revision",
			},
		}

		varType := client.ConfigurationVariableTypeEnvironment
		varSchema := client.ConfigurationVariableSchema{
			Type: "string",
			Enum: []string{"a", "b"},
		}
		configurationVariables := client.ConfigurationVariable{
			Value:  "my env var value",
			Name:   "my env var",
			Type:   &varType,
			Schema: &varSchema,
		}

		updatedEnvironmentResource := fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					configuration {
						name = "%s"
						value = "%s"
						schema_type = "%s"
						schema_enum = ["%s"]
					}
				}`,
			resourceType, resourceName, environment.Name,
			updatedEnvironment.ProjectId, updatedEnvironment.LatestDeploymentLog.BlueprintId,
			updatedEnvironment.LatestDeploymentLog.BlueprintRevision, configurationVariables.Name,
			configurationVariables.Value, configurationVariables.Schema.Type,
			strings.Join(configurationVariables.Schema.Enum, "\",\""),
		)

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        environment.Name,
						"project_id":  environment.ProjectId,
						"template_id": environment.LatestDeploymentLog.BlueprintId,
						"revision":    environment.LatestDeploymentLog.BlueprintRevision,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
					),
				},
				{
					Config: updatedEnvironmentResource,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", updatedEnvironment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "revision", updatedEnvironment.LatestDeploymentLog.BlueprintRevision),
						resource.TestCheckResourceAttr(accessor, "configuration.0.name", configurationVariables.Name),
						resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Value),
						resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", configurationVariables.Schema.Type),
						resource.TestCheckResourceAttr(accessor, "configuration.0.schema_enum.0", configurationVariables.Schema.Enum[0]),
						resource.TestCheckResourceAttr(accessor, "configuration.0.schema_enum.1", configurationVariables.Schema.Enum[1]),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
				},
			}).Times(1).Return(environment, nil)

			mock.EXPECT().EnvironmentDeploy(environment.Id, gomock.Any()).Times(1).Return(client.EnvironmentDeployResponse{
				Id: "deployment-id",
			}, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("TTL update", func(t *testing.T) {
		environment := client.Environment{
			Id:            "id0",
			Name:          "my-environment",
			ProjectId:     "project-id",
			LifespanEndAt: "2021-12-08T11:45:11Z",
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId: "template-id",
			},
		}
		updatedEnvironment := client.Environment{
			Id:            updatedEnvironment.Id,
			Name:          environment.Name,
			ProjectId:     environment.ProjectId,
			LifespanEndAt: "2021-12-08T12:45:11Z",
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId: environment.LatestDeploymentLog.BlueprintId,
			},
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        environment.Name,
						"project_id":  environment.ProjectId,
						"template_id": environment.LatestDeploymentLog.BlueprintId,
						"ttl":         environment.LifespanEndAt,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "ttl", environment.LifespanEndAt),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedEnvironment.Name,
						"project_id":  updatedEnvironment.ProjectId,
						"template_id": updatedEnvironment.LatestDeploymentLog.BlueprintId,
						"ttl":         updatedEnvironment.LifespanEndAt,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", updatedEnvironment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "ttl", updatedEnvironment.LifespanEndAt),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdateTTL(environment.Id, client.TTL{
				Type:  client.TTLTypeDate,
				Value: updatedEnvironment.LifespanEndAt,
			}).Times(1).Return(environment, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("Deleting TTL from environment should update ttl to infinite", func(t *testing.T) {
		environment := client.Environment{
			Id:            "id0",
			Name:          "my-environment",
			ProjectId:     "project-id",
			LifespanEndAt: "2021-12-08T11:45:11Z",
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId: "template-id",
			},
		}
		updatedEnvironment := client.Environment{
			Id:        updatedEnvironment.Id,
			Name:      environment.Name,
			ProjectId: environment.ProjectId,
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId: environment.LatestDeploymentLog.BlueprintId,
			},
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        environment.Name,
						"project_id":  environment.ProjectId,
						"template_id": environment.LatestDeploymentLog.BlueprintId,
						"ttl":         environment.LifespanEndAt,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "ttl", environment.LifespanEndAt),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedEnvironment.Name,
						"project_id":  updatedEnvironment.ProjectId,
						"template_id": updatedEnvironment.LatestDeploymentLog.BlueprintId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", updatedEnvironment.LatestDeploymentLog.BlueprintId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdateTTL(environment.Id, client.TTL{
				Type:  client.TTlTypeInfinite,
				Value: "",
			}).Times(1).Return(environment, nil)

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

	t.Run("Failure in deploy", func(t *testing.T) {
		updatedEnvironment := client.Environment{
			Id:        updatedEnvironment.Id,
			Name:      environment.Name,
			ProjectId: environment.ProjectId,
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision: "updated template id",
			},
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createEnvironmentResourceConfig(environment),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedEnvironment.Name,
						"project_id":  updatedEnvironment.ProjectId,
						"template_id": updatedEnvironment.LatestDeploymentLog.BlueprintId,
						"revision":    updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
					}),
					ExpectError: regexp.MustCompile("failed deploying environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentDeploy(updatedEnvironment.Id, client.DeployRequest{
				BlueprintId:       updatedEnvironment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision: updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
			}).Times(1).Return(client.EnvironmentDeployResponse{}, errors.New("error"))
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
