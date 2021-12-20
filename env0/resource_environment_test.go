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
			"name":          environment.Name,
			"project_id":    environment.ProjectId,
			"template_id":   environment.LatestDeploymentLog.BlueprintId,
			"force_destroy": true,
		})
	}
	autoDeployOnPathChangesOnlyDefault := false
	autoDeployByCustomGlobDefault := ""
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
				Name:                        environment.Name,
				ProjectId:                   environment.ProjectId,
				AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
				AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
				Name: updatedEnvironment.Name,
			}).Times(1).Return(updatedEnvironment, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("Success in create and deploy with variables update", func(t *testing.T) {
		environment := client.Environment{
			Id:                          "id0",
			Name:                        "my-environment",
			ProjectId:                   "project-id",
			AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
			AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
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
			Value:  varSchema.Enum[0],
			Name:   "my env var",
			Type:   &varType,
			Schema: &varSchema,
		}
		formatVariables := func(variables []client.ConfigurationVariable) string {
			foramt := ""
			for _, variable := range variables {
				schemaFormat := ""
				if variable.Schema != nil {
					schemaFormat = fmt.Sprintf(`
									schema_type = "%s"
									schema_enum = ["%s"]
									`, variable.Schema.Type,
						strings.Join(variable.Schema.Enum, "\",\""))
				}
				varType := "environment"
				if *variable.Type == client.ConfigurationVariableTypeTerraform {
					varType = "terraform"
				}
				foramt += fmt.Sprintf(`configuration{
									name = "%s"
									value = "%s"
									type = "%s"
									%s
									}

							`, variable.Name,
					variable.Value, varType, schemaFormat)
			}

			return foramt
		}
		formatResourceWithConfiguration := func(env client.Environment, variables []client.ConfigurationVariable) string {
			return fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					force_destroy = true
					%s

				}`,
				resourceType, resourceName, env.Name,
				env.ProjectId, env.LatestDeploymentLog.BlueprintId,
				env.LatestDeploymentLog.BlueprintRevision, formatVariables(variables))
		}

		environmentResource := formatResourceWithConfiguration(environment, client.ConfigurationChanges{configurationVariables})
		newVarType := client.ConfigurationVariableTypeTerraform
		redeployConfigurationVariables := client.ConfigurationChanges{client.ConfigurationVariable{
			Value: "configurationVariables.Value",
			Name:  configurationVariables.Name,
			Type:  &newVarType,
		}}
		updatedEnvironmentResource := formatResourceWithConfiguration(updatedEnvironment, redeployConfigurationVariables)
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: environmentResource,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
						resource.TestCheckResourceAttr(accessor, "configuration.0.name", configurationVariables.Name),
						resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Schema.Enum[0]),
						resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", configurationVariables.Schema.Type),
						resource.TestCheckResourceAttr(accessor, "configuration.0.schema_enum.0", configurationVariables.Schema.Enum[0]),
						resource.TestCheckResourceAttr(accessor, "configuration.0.schema_enum.1", configurationVariables.Schema.Enum[1]),
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
						resource.TestCheckResourceAttr(accessor, "configuration.0.value", "configurationVariables.Value"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			isSensitive := false
			configurationVariables.Scope = client.ScopeDeployment
			configurationVariables.IsSensitive = &isSensitive
			configurationVariables.Value = configurationVariables.Schema.Enum[0]
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:                        environment.Name,
				ProjectId:                   environment.ProjectId,
				AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
				AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
				DeployRequest: &client.DeployRequest{
					BlueprintId:          environment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision:    environment.LatestDeploymentLog.BlueprintRevision,
					ConfigurationChanges: &client.ConfigurationChanges{configurationVariables},
				}, ConfigurationChanges: &client.ConfigurationChanges{configurationVariables},
			}).Times(1).Return(environment, nil)
			configurationVariables.Id = "generated-id-from-server"

			varTrue := true
			configurationVariables.ToDelete = &varTrue
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{configurationVariables}, nil), // read after create -> on update
				mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, updatedEnvironment.Id).Times(1).Return(redeployConfigurationVariables, nil),                      // read after create -> on update -> read after update
			)
			redeployConfigurationVariables[0].Scope = client.ScopeDeployment
			redeployConfigurationVariables[0].IsSensitive = &isSensitive
			mock.EXPECT().EnvironmentDeploy(environment.Id, client.DeployRequest{
				BlueprintId:          environment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision:    updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
				ConfigurationChanges: &client.ConfigurationChanges{redeployConfigurationVariables[0], configurationVariables},
			}).Times(1).Return(client.EnvironmentDeployResponse{
				Id: environment.Id,
			}, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("Update to: revision, configuration should trigger a deployment", func(t *testing.T) {
		environment := client.Environment{
			Id:                          "id0",
			Name:                        "my-environment",
			ProjectId:                   "project-id",
			AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
			AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
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
					force_destroy = true
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
						"name":          environment.Name,
						"project_id":    environment.ProjectId,
						"template_id":   environment.LatestDeploymentLog.BlueprintId,
						"revision":      environment.LatestDeploymentLog.BlueprintRevision,
						"force_destroy": true,
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
				Name:                        environment.Name,
				ProjectId:                   environment.ProjectId,
				AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
				AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
				DeployRequest: &client.DeployRequest{
					BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
				},
			}).Times(1).Return(environment, nil)

			mock.EXPECT().EnvironmentDeploy(environment.Id, gomock.Any()).Times(1).Return(client.EnvironmentDeployResponse{
				Id: "deployment-id",
			}, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, updatedEnvironment.Id).Times(4).Return(client.ConfigurationChanges{}, nil) // read after create -> on update -> read after update
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
						"name":          updatedEnvironment.Name,
						"project_id":    updatedEnvironment.ProjectId,
						"template_id":   updatedEnvironment.LatestDeploymentLog.BlueprintId,
						"ttl":           updatedEnvironment.LifespanEndAt,
						"force_destroy": true,
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
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil) // read after create -> on update -> read after update

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
						"name":          updatedEnvironment.Name,
						"project_id":    updatedEnvironment.ProjectId,
						"template_id":   updatedEnvironment.LatestDeploymentLog.BlueprintId,
						"force_destroy": true,
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
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("Deleting triggers should set them to false", func(t *testing.T) {
		falsey := false
		truthyFruity := true
		environment := client.Environment{
			Id:        "id0",
			Name:      "my-environment",
			ProjectId: "project-id",
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId: "template-id",
			},

			ContinuousDeployment:        &truthyFruity,
			AutoDeployOnPathChangesOnly: &truthyFruity,
			RequiresApproval:            &falsey,
			PullRequestPlanDeployments:  &truthyFruity,
			AutoDeployByCustomGlob:      ".*",
		}
		environmentAfterUpdate := client.Environment{
			Id:        environment.Id,
			Name:      environment.Name,
			ProjectId: environment.ProjectId,
			LatestDeploymentLog: client.DeploymentLog{
				BlueprintId: environment.LatestDeploymentLog.BlueprintId,
			},

			ContinuousDeployment:        &falsey,
			AutoDeployOnPathChangesOnly: &falsey,
			RequiresApproval:            &truthyFruity,
			PullRequestPlanDeployments:  &falsey,
			AutoDeployByCustomGlob:      "",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                             environment.Name,
						"project_id":                       environment.ProjectId,
						"template_id":                      environment.LatestDeploymentLog.BlueprintId,
						"deploy_on_push":                   *environment.ContinuousDeployment,
						"approve_plan_automatically":       !*environment.RequiresApproval,
						"run_plan_on_pull_requests":        *environment.PullRequestPlanDeployments,
						"auto_deploy_on_path_changes_only": *environment.AutoDeployOnPathChangesOnly,
						"auto_deploy_by_custom_glob":       environment.AutoDeployByCustomGlob,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", "true"),
						resource.TestCheckResourceAttr(accessor, "run_plan_on_pull_requests", "true"),
						resource.TestCheckResourceAttr(accessor, "auto_deploy_on_path_changes_only", "true"),
						resource.TestCheckResourceAttr(accessor, "auto_deploy_by_custom_glob", ".*"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":          environment.Name,
						"project_id":    environment.ProjectId,
						"template_id":   environment.LatestDeploymentLog.BlueprintId,
						"force_destroy": true,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", "false"),
						resource.TestCheckResourceAttr(accessor, "run_plan_on_pull_requests", "false"),
						resource.TestCheckResourceAttr(accessor, "auto_deploy_on_path_changes_only", "false"),
						resource.TestCheckResourceAttr(accessor, "auto_deploy_by_custom_glob", ""),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdate(environment.Id, client.EnvironmentUpdate{
				Name: environment.Name,

				ContinuousDeployment:        &falsey,
				AutoDeployOnPathChangesOnly: &falsey,
				RequiresApproval:            &truthyFruity,
				PullRequestPlanDeployments:  &falsey,
				AutoDeployByCustomGlob:      "",
			}).Times(1).Return(environmentAfterUpdate, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),            // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(environmentAfterUpdate, nil), // 1 after update
			)
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, environment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("should only allow destroy when force destroy is enabled", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        environment.Name,
						"project_id":  environment.ProjectId,
						"template_id": environment.LatestDeploymentLog.BlueprintId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", templateId),
					),
				},
				{
					Destroy: true,
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":          environment.Name,
						"project_id":    environment.ProjectId,
						"template_id":   environment.LatestDeploymentLog.BlueprintId,
						"force_destroy": true,
					}),
					ExpectError: regexp.MustCompile(`must enable "force_destroy" safeguard in order to destroy`),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":          environment.Name,
						"project_id":    environment.ProjectId,
						"template_id":   environment.LatestDeploymentLog.BlueprintId,
						"force_destroy": true,
					}),
				},
				{
					Destroy: true,
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":          environment.Name,
						"project_id":    environment.ProjectId,
						"template_id":   environment.LatestDeploymentLog.BlueprintId,
						"force_destroy": true,
					}),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, environment.Id).Times(5).Return(client.ConfigurationChanges{}, nil)
			mock.EXPECT().Environment(gomock.Any()).Times(5).Return(environment, nil)
			mock.EXPECT().EnvironmentDestroy(gomock.Any()).Times(1).Return(environment, nil)

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
				Name:                        environment.Name,
				ProjectId:                   environment.ProjectId,
				AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
				AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
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
				Name:                        environment.Name,
				ProjectId:                   environment.ProjectId,
				AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
				AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, environment.Id).Times(2).Return(client.ConfigurationChanges{}, nil)
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
						"name":          updatedEnvironment.Name,
						"project_id":    updatedEnvironment.ProjectId,
						"template_id":   updatedEnvironment.LatestDeploymentLog.BlueprintId,
						"revision":      updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
						"force_destroy": true,
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
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, environment.Id).Times(2).Return(client.ConfigurationChanges{}, nil)
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
				Name:                        environment.Name,
				ProjectId:                   environment.ProjectId,
				AutoDeployOnPathChangesOnly: &autoDeployOnPathChangesOnlyDefault,
				AutoDeployByCustomGlob:      autoDeployByCustomGlobDefault,
				DeployRequest: &client.DeployRequest{
					BlueprintId: templateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().Environment(gomock.Any()).Return(client.Environment{}, errors.New("error"))
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			mock.EXPECT().ConfigurationVariables(client.ScopeEnvironment, environment.Id).Times(0).Return(client.ConfigurationChanges{}, nil)
		})

	})
}
