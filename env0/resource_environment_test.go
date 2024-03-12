package env0

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitEnvironmentResource(t *testing.T) {
	resourceType := "env0_environment"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)
	templateId := "template-id"
	deploymentLogId := "deploymentLogId0"
	isRemoteBackendTrue := true
	isRemoteBackendFalse := false

	driftDetectionCron := "*/5 * * * *"
	updatedDriftDetectionCron := "*/10 1 * * *"

	environment := client.Environment{
		Id:            uuid.New().String(),
		Name:          "my-environment",
		ProjectId:     "project-id",
		WorkspaceName: "workspace-name",
		LatestDeploymentLog: client.DeploymentLog{
			Id:                deploymentLogId,
			BlueprintId:       templateId,
			BlueprintRevision: "revision",
			Output:            []byte(`{"a": "b"}`),
		},
		TerragruntWorkingDirectory: "/terragrunt/directory/",
		VcsCommandsAlias:           "alias",
		IsRemoteBackend:            &isRemoteBackendFalse,
	}

	updatedEnvironment := client.Environment{
		Id:            environment.Id,
		Name:          "my-updated-environment-name",
		ProjectId:     "project-id",
		WorkspaceName: environment.WorkspaceName,
		LatestDeploymentLog: client.DeploymentLog{
			Id:                deploymentLogId,
			BlueprintId:       templateId,
			BlueprintRevision: "revision",
			Output:            []byte(`{"a": "b"}`),
		},
		TerragruntWorkingDirectory: "/terragrunt/directory/",
		VcsCommandsAlias:           "alias2",
		IsRemoteBackend:            &isRemoteBackendTrue,
		IsArchived:                 boolPtr(true),
	}

	template := client.Template{
		ProjectId: updatedEnvironment.ProjectId,
	}

	templateInSlice := client.Template{
		ProjectIds: []string{updatedEnvironment.ProjectId},
	}

	createEnvironmentResourceConfig := func(environment client.Environment) string {
		config := map[string]interface{}{
			"name":                         environment.Name,
			"project_id":                   environment.ProjectId,
			"template_id":                  environment.LatestDeploymentLog.BlueprintId,
			"workspace":                    environment.WorkspaceName,
			"terragrunt_working_directory": environment.TerragruntWorkingDirectory,
			"force_destroy":                true,
			"vcs_commands_alias":           environment.VcsCommandsAlias,
			"is_remote_backend":            *(environment.IsRemoteBackend),
		}

		if environment.IsArchived != nil {
			config["is_inactive"] = *(environment.IsArchived)
		}

		return resourceConfigCreate(resourceType, resourceName, config)
	}

	createEnvironmentResourceConfigDriftDetection := func(environment client.Environment, cron string) string {
		config := map[string]interface{}{
			"name":                         environment.Name,
			"project_id":                   environment.ProjectId,
			"template_id":                  environment.LatestDeploymentLog.BlueprintId,
			"workspace":                    environment.WorkspaceName,
			"terragrunt_working_directory": environment.TerragruntWorkingDirectory,
			"force_destroy":                true,
			"vcs_commands_alias":           environment.VcsCommandsAlias,
			"is_remote_backend":            *(environment.IsRemoteBackend),
		}

		if environment.IsArchived != nil {
			config["is_inactive"] = *(environment.IsArchived)
		}

		config["drift_detection_cron"] = cron

		return resourceConfigCreate(resourceType, resourceName, config)
	}

	autoDeployByCustomGlobDefault := ""

	testSuccess := func() {
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
							resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
							resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", environment.TerragruntWorkingDirectory),
							resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", environment.VcsCommandsAlias),
							resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "false"),
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckNoResourceAttr(accessor, "deploy_on_push"),
							resource.TestCheckNoResourceAttr(accessor, "run_plan_on_pull_requests"),
						),
					},
					{
						Config: createEnvironmentResourceConfig(updatedEnvironment),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
							resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
							resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", updatedEnvironment.TerragruntWorkingDirectory),
							resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", updatedEnvironment.VcsCommandsAlias),
							resource.TestCheckResourceAttr(accessor, "revision", updatedEnvironment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "true"),
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckResourceAttr(accessor, "is_inactive", "true"),
							resource.TestCheckNoResourceAttr(accessor, "deploy_on_push"),
							resource.TestCheckNoResourceAttr(accessor, "run_plan_on_pull_requests"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                       environment.Name,
					ProjectId:                  environment.ProjectId,
					WorkspaceName:              environment.WorkspaceName,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					IsRemoteBackend: &isRemoteBackendFalse,
				}).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
					Name:                       updatedEnvironment.Name,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: updatedEnvironment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           updatedEnvironment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendTrue,
					IsArchived:                 updatedEnvironment.IsArchived,
				}).Times(1).Return(updatedEnvironment, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("prevent auto deploy", func(t *testing.T) {
			templateId := "template-id"
			newTemplateId := "new-template-id"
			truthyFruity := true

			environment := client.Environment{
				Id:        uuid.New().String(),
				Name:      "name",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: templateId,
				},
			}

			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                environment.Name,
							"project_id":          environment.ProjectId,
							"template_id":         templateId,
							"force_destroy":       true,
							"prevent_auto_deploy": true,
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "prevent_auto_deploy", "true"),
						),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                environment.Name,
							"project_id":          environment.ProjectId,
							"template_id":         newTemplateId,
							"force_destroy":       true,
							"prevent_auto_deploy": true,
						}),
						ExpectError: regexp.MustCompile("template_id may not be modified, create a new environment instead"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil),
					mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
						Name:      environment.Name,
						ProjectId: environment.ProjectId,
						DeployRequest: &client.DeployRequest{
							BlueprintId: templateId,
						},

						PreventAutoDeploy: &truthyFruity,
					}).Times(1).Return(environment, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1),
				)
			})
		})

		t.Run("avoid modifying template id", func(t *testing.T) {
			templateId := "template-id"
			newTemplateId := "new-template-id"

			environment := client.Environment{
				Id:        uuid.New().String(),
				Name:      "name",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: templateId,
				},
			}

			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":          environment.Name,
							"project_id":    environment.ProjectId,
							"template_id":   templateId,
							"force_destroy": true,
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
						),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":          environment.Name,
							"project_id":    environment.ProjectId,
							"template_id":   newTemplateId,
							"force_destroy": true,
						}),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile("template_id may not be modified, create a new environment instead"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil),
					mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
						Name:      environment.Name,
						ProjectId: environment.ProjectId,

						DeployRequest: &client.DeployRequest{
							BlueprintId: templateId,
						},
					}).Times(1).Return(environment, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1),
				)
			})
		})

		t.Run("remote apply is enabled", func(t *testing.T) {
			templateId := "template-id"

			environment := client.Environment{
				Id:        uuid.New().String(),
				Name:      "name",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: templateId,
				},
				IsRemoteBackend:  boolPtr(true),
				RequiresApproval: boolPtr(false),
			}

			updatedEnvironment := client.Environment{
				Id:        environment.Id,
				Name:      "name",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: templateId,
				},
				IsRemoteApplyEnabled: true,
			}

			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                       environment.Name,
							"project_id":                 environment.ProjectId,
							"template_id":                templateId,
							"is_remote_backend":          *environment.IsRemoteBackend,
							"approve_plan_automatically": !*environment.RequiresApproval,
							"force_destroy":              true,
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "true"),
							resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", "true"),
							resource.TestCheckResourceAttr(accessor, "is_remote_apply_enabled", "false"),
						),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                       environment.Name,
							"project_id":                 environment.ProjectId,
							"template_id":                templateId,
							"is_remote_backend":          *environment.IsRemoteBackend,
							"approve_plan_automatically": !*environment.RequiresApproval,
							"force_destroy":              true,
							"is_remote_apply_enabled":    true,
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "true"),
							resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", "true"),
							resource.TestCheckResourceAttr(accessor, "is_remote_apply_enabled", "true"),
						),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                       environment.Name,
							"project_id":                 environment.ProjectId,
							"template_id":                templateId,
							"is_remote_backend":          !*environment.IsRemoteBackend,
							"approve_plan_automatically": !*environment.RequiresApproval,
							"force_destroy":              true,
							"is_remote_apply_enabled":    true,
						}),
						ExpectError: regexp.MustCompile("cannot set is_remote_apply_enabled when approve_plan_automatically or is_remote_backend are disabled"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil),
					mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
						Name:      environment.Name,
						ProjectId: environment.ProjectId,

						DeployRequest: &client.DeployRequest{
							BlueprintId: templateId,
						},
						IsRemoteBackend:  environment.IsRemoteBackend,
						RequiresApproval: environment.RequiresApproval,
					}).Times(1).Return(environment, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
						Name:                 updatedEnvironment.Name,
						IsRemoteBackend:      updatedEnvironment.IsRemoteBackend,
						RequiresApproval:     updatedEnvironment.RequiresApproval,
						IsRemoteApplyEnabled: updatedEnvironment.IsRemoteApplyEnabled,
					}).Times(1).Return(updatedEnvironment, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(updatedEnvironment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(updatedEnvironment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
					mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1),
				)
			})
		})

		t.Run("Import By Id", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: createEnvironmentResourceConfig(environment),
					},
					{
						ResourceName:            resourceNameImport,
						ImportState:             true,
						ImportStateId:           environment.Id,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"force_destroy"},
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                       environment.Name,
					ProjectId:                  environment.ProjectId,
					WorkspaceName:              environment.WorkspaceName,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					IsRemoteBackend: &isRemoteBackendFalse,
				}).Times(1).Return(environment, nil)
				mock.EXPECT().Environment(environment.Id).Times(3).Return(environment, nil)
				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
			})
		})

		t.Run("Success create and remove drift cron", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: createEnvironmentResourceConfigDriftDetection(environment, driftDetectionCron),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
							resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", environment.TerragruntWorkingDirectory),
							resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", environment.VcsCommandsAlias),
							resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "false"),
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckResourceAttr(accessor, "drift_detection_cron", driftDetectionCron),
							resource.TestCheckNoResourceAttr(accessor, "deploy_on_push"),
							resource.TestCheckNoResourceAttr(accessor, "run_plan_on_pull_requests"),
						),
					},
					{
						Config: createEnvironmentResourceConfig(updatedEnvironment),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
							resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
							resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", updatedEnvironment.TerragruntWorkingDirectory),
							resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", updatedEnvironment.VcsCommandsAlias),
							resource.TestCheckResourceAttr(accessor, "revision", updatedEnvironment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "true"),
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckResourceAttr(accessor, "is_inactive", "true"),
							resource.TestCheckResourceAttr(accessor, "drift_detection_cron", ""),
							resource.TestCheckNoResourceAttr(accessor, "deploy_on_push"),
							resource.TestCheckNoResourceAttr(accessor, "run_plan_on_pull_requests"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                       environment.Name,
					ProjectId:                  environment.ProjectId,
					WorkspaceName:              environment.WorkspaceName,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					IsRemoteBackend: &isRemoteBackendFalse,
					DriftDetectionRequest: &client.DriftDetectionRequest{
						Enabled: true,
						Cron:    driftDetectionCron,
					},
				}).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
					Name:                       updatedEnvironment.Name,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: updatedEnvironment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           updatedEnvironment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendTrue,
					IsArchived:                 updatedEnvironment.IsArchived,
				}).Times(1).Return(updatedEnvironment, nil)
				mock.EXPECT().EnvironmentStopDriftDetection(environment.Id).Times(1).Return(nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Success create and update drift cron", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: createEnvironmentResourceConfigDriftDetection(environment, driftDetectionCron),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
							resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", environment.TerragruntWorkingDirectory),
							resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", environment.VcsCommandsAlias),
							resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "false"),
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckResourceAttr(accessor, "drift_detection_cron", driftDetectionCron),
							resource.TestCheckNoResourceAttr(accessor, "deploy_on_push"),
							resource.TestCheckNoResourceAttr(accessor, "run_plan_on_pull_requests"),
						),
					},
					{
						Config: createEnvironmentResourceConfigDriftDetection(updatedEnvironment, updatedDriftDetectionCron),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
							resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", templateId),
							resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
							resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", updatedEnvironment.TerragruntWorkingDirectory),
							resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", updatedEnvironment.VcsCommandsAlias),
							resource.TestCheckResourceAttr(accessor, "revision", updatedEnvironment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "is_remote_backend", "true"),
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckResourceAttr(accessor, "is_inactive", "true"),
							resource.TestCheckResourceAttr(accessor, "drift_detection_cron", updatedDriftDetectionCron),
							resource.TestCheckNoResourceAttr(accessor, "deploy_on_push"),
							resource.TestCheckNoResourceAttr(accessor, "run_plan_on_pull_requests"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                       environment.Name,
					ProjectId:                  environment.ProjectId,
					WorkspaceName:              environment.WorkspaceName,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					IsRemoteBackend: &isRemoteBackendFalse,
					DriftDetectionRequest: &client.DriftDetectionRequest{
						Enabled: true,
						Cron:    driftDetectionCron,
					},
				}).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
					Name:                       updatedEnvironment.Name,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: updatedEnvironment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           updatedEnvironment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendTrue,
					IsArchived:                 updatedEnvironment.IsArchived,
				}).Times(1).Return(updatedEnvironment, nil)
				mock.EXPECT().EnvironmentUpdateDriftDetection(environment.Id, client.EnvironmentSchedulingExpression{Cron: updatedDriftDetectionCron, Enabled: true}).Times(1).Return(client.EnvironmentSchedulingExpression{}, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Success in create and deploy with variables update", func(t *testing.T) {
			environment := client.Environment{
				Id:                     environment.Id,
				Name:                   "my-environment",
				ProjectId:              "project-id",
				AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
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
					Output:            []byte(`{"a": "b"}`),
				},
			}

			varType := client.ConfigurationVariableTypeEnvironment
			varSchema := client.ConfigurationVariableSchema{
				Type:   "string",
				Enum:   []string{"a", "b"},
				Format: client.HCL,
			}
			configurationVariables := client.ConfigurationVariable{
				Value:  varSchema.Enum[0],
				Name:   "my env var",
				Type:   &varType,
				Schema: &varSchema,
				Regex:  "regex",
			}
			formatVariables := func(variables []client.ConfigurationVariable) string {
				format := ""
				for _, variable := range variables {
					schemaFormat := ""
					if variable.Schema != nil {
						if variable.Schema.Enum != nil {
							schemaFormat = fmt.Sprintf(`
									schema_type = "%s"
									schema_enum = ["%s"]
									schema_format = "%s"
									`, variable.Schema.Type,
								strings.Join(variable.Schema.Enum, "\",\""), variable.Schema.Format)
						} else {
							schemaFormat = fmt.Sprintf(`
									schema_type = "%s"
									schema_format = "%s"
									`, variable.Schema.Type,
								variable.Schema.Format)
						}

					}
					varType := "environment"
					if *variable.Type == client.ConfigurationVariableTypeTerraform {
						varType = "terraform"
					}
					format += fmt.Sprintf(`configuration{
									name = "%s"
									value = "%s"
									type = "%s"
									regex = "%s"
									%s
									}

							`, variable.Name,
						variable.Value, varType, variable.Regex, schemaFormat)
				}

				return format
			}

			formatResourceWithConfiguration := func(env client.Environment, variables []client.ConfigurationVariable) string {
				output := "null"
				if len(env.LatestDeploymentLog.Output) > 0 {
					output = strings.ReplaceAll(string(env.LatestDeploymentLog.Output), `"`, `\"`)
				}

				return fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					output = "%s"
					force_destroy = true
					%s

				}`,
					resourceType,
					resourceName,
					env.Name,
					env.ProjectId,
					env.LatestDeploymentLog.BlueprintId,
					env.LatestDeploymentLog.BlueprintRevision,
					output,
					formatVariables(variables))
			}

			environmentResource := formatResourceWithConfiguration(environment, client.ConfigurationChanges{configurationVariables})
			newVarType := client.ConfigurationVariableTypeTerraform
			redeployConfigurationVariables := client.ConfigurationChanges{client.ConfigurationVariable{
				Value: configurationVariables.Value,
				Name:  configurationVariables.Name,
				Type:  &newVarType,
				Schema: &client.ConfigurationVariableSchema{
					Format: client.Text,
				},
				Regex: "regex2",
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
							resource.TestCheckResourceAttr(accessor, "output", "null"),
							resource.TestCheckResourceAttr(accessor, "configuration.0.name", configurationVariables.Name),
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Schema.Enum[0]),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", configurationVariables.Schema.Type),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_format", string(configurationVariables.Schema.Format)),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_enum.0", configurationVariables.Schema.Enum[0]),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_enum.1", configurationVariables.Schema.Enum[1]),
							resource.TestCheckResourceAttr(accessor, "configuration.0.regex", configurationVariables.Regex),
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
							resource.TestCheckResourceAttr(accessor, "output", string(updatedEnvironment.LatestDeploymentLog.Output)),
							resource.TestCheckResourceAttr(accessor, "configuration.0.name", configurationVariables.Name),
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Value),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_format", string(client.Text)),
							resource.TestCheckResourceAttr(accessor, "configuration.0.regex", "regex2"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				isSensitive := false
				configurationVariables.Scope = client.ScopeDeployment
				configurationVariables.IsSensitive = &isSensitive
				configurationVariables.IsReadOnly = boolPtr(false)
				configurationVariables.IsRequired = boolPtr(false)
				configurationVariables.Value = configurationVariables.Schema.Enum[0]
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(templateInSlice, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                   environment.Name,
					ProjectId:              environment.ProjectId,
					WorkspaceName:          environment.WorkspaceName,
					AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
					DeployRequest: &client.DeployRequest{
						BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
						BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
					}, ConfigurationChanges: &client.ConfigurationChanges{configurationVariables},
				}).Times(1).Return(environment, nil)
				configurationVariables.Id = "generated-id-from-server"

				varTrue := true
				configurationVariables.ToDelete = &varTrue
				gomock.InOrder(
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{configurationVariables}, nil), // read after create -> on update
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(1).Return(redeployConfigurationVariables, nil),                      // read after create -> on update -> read after update
				)
				redeployConfigurationVariables[0].Scope = client.ScopeDeployment
				redeployConfigurationVariables[0].IsSensitive = &isSensitive
				redeployConfigurationVariables[0].IsReadOnly = boolPtr(false)
				redeployConfigurationVariables[0].IsRequired = boolPtr(false)

				deployRequest := client.DeployRequest{
					BlueprintId:          environment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision:    updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
					ConfigurationChanges: &client.ConfigurationChanges{redeployConfigurationVariables[0], configurationVariables},
				}

				mock.EXPECT().EnvironmentDeploy(environment.Id, deployRequest).Times(1).Return(client.EnvironmentDeployResponse{
					Id: environment.Id,
				}, nil)

				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Create configuration variables - default values", func(t *testing.T) {
			environment := client.Environment{
				Id:        environment.Id,
				Name:      "my-environment",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId:       "template-id",
					BlueprintRevision: "revision",
				},
			}

			varType := client.ConfigurationVariableTypeEnvironment
			varSchema := client.ConfigurationVariableSchema{
				Type: "string",
			}
			configurationVariables := client.ConfigurationVariable{
				Value:  "my env var value",
				Name:   "my env var",
				Type:   &varType,
				Schema: &varSchema,
			}

			environmentResource := fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					force_destroy = true
					configuration {
						name = "%s"
						value = "%s"
					}
				}`,
				resourceType, resourceName, environment.Name,
				environment.ProjectId, environment.LatestDeploymentLog.BlueprintId,
				environment.LatestDeploymentLog.BlueprintRevision, configurationVariables.Name,
				configurationVariables.Value,
			)

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
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Value),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", "string"),
							resource.TestCheckNoResourceAttr(accessor, "configuration.0.schema_enum"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)

				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariables}, nil) // read after create -> on update -> read after update
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Create and redeploy configuration variables - sensitive values", func(t *testing.T) {
			environment := client.Environment{
				Id:        environment.Id,
				Name:      "my-environment",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId:       "template-id",
					BlueprintRevision: "revision",
				},
			}

			varType := client.ConfigurationVariableTypeEnvironment
			varSchema := client.ConfigurationVariableSchema{
				Type: "string",
			}
			configurationVariables := client.ConfigurationVariable{
				Value:       "my env var value",
				Name:        "my env var",
				Type:        &varType,
				Schema:      &varSchema,
				IsSensitive: boolPtr(true),
				Scope:       client.ScopeDeployment,
				IsReadOnly:  boolPtr(false),
				IsRequired:  boolPtr(false),
			}

			redeployConfigurationVariable := client.ConfigurationVariable{
				Value:       configurationVariables.Value + "1",
				Name:        "my env var",
				Type:        &varType,
				Schema:      &varSchema,
				IsSensitive: boolPtr(true),
				Scope:       client.ScopeDeployment,
				IsReadOnly:  boolPtr(false),
				IsRequired:  boolPtr(false),
			}

			environmentResource := fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					force_destroy = true
					configuration {
						name = "%s"
						value = "%s"
						is_sensitive = true
					}
				}`,
				resourceType, resourceName, environment.Name,
				environment.ProjectId, environment.LatestDeploymentLog.BlueprintId,
				environment.LatestDeploymentLog.BlueprintRevision, configurationVariables.Name,
				configurationVariables.Value,
			)

			updateEnvironmentResource := fmt.Sprintf(`
			resource "%s" "%s" {
				name = "%s"
				project_id = "%s"
				template_id = "%s"
				revision = "%s"
				force_destroy = true
				configuration {
					name = "%s"
					value = "%s"
					is_sensitive = true
				}
			}`,
				resourceType, resourceName, environment.Name,
				environment.ProjectId, environment.LatestDeploymentLog.BlueprintId,
				environment.LatestDeploymentLog.BlueprintRevision, configurationVariables.Name,
				redeployConfigurationVariable.Value,
			)

			environmentCreate := client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
				},
				ConfigurationChanges: &client.ConfigurationChanges{
					configurationVariables,
				},
			}

			environmentDeploy := client.DeployRequest{
				BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
				ConfigurationChanges: &client.ConfigurationChanges{
					redeployConfigurationVariable,
				},
			}

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
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Value),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", "string"),
							resource.TestCheckResourceAttr(accessor, "configuration.0.is_sensitive", "true"),
							resource.TestCheckNoResourceAttr(accessor, "configuration.0.schema_enum"),
						),
					},
					{
						Config: updateEnvironmentResource,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
							resource.TestCheckResourceAttr(accessor, "revision", environment.LatestDeploymentLog.BlueprintRevision),
							resource.TestCheckResourceAttr(accessor, "configuration.0.name", configurationVariables.Name),
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", redeployConfigurationVariable.Value),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", "string"),
							resource.TestCheckResourceAttr(accessor, "configuration.0.is_sensitive", "true"),
							resource.TestCheckResourceAttr(accessor, "deployment_id", "did"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				gomock.InOrder(
					mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil),
					mock.EXPECT().EnvironmentCreate(environmentCreate).Times(1).Return(environment, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariables}, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(2).Return(client.ConfigurationChanges{configurationVariables}, nil),
					mock.EXPECT().EnvironmentDeploy(environment.Id, environmentDeploy).Times(1).Return(client.EnvironmentDeployResponse{Id: "did"}, nil),
					mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{redeployConfigurationVariable}, nil),
					mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1),
				)
			})
		})

		t.Run("Create configuration variables - schema type string", func(t *testing.T) {
			environment := client.Environment{
				Id:        environment.Id,
				Name:      "my-environment",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId:       "template-id",
					BlueprintRevision: "revision",
				},
			}

			varType := client.ConfigurationVariableTypeEnvironment
			varSchema := client.ConfigurationVariableSchema{
				Type: "string",
			}
			configurationVariables := client.ConfigurationVariable{
				Value:  "my env var value",
				Name:   "my env var",
				Type:   &varType,
				Schema: &varSchema,
			}

			environmentResource := fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					force_destroy = true
					configuration {
						name = "%s"
						value = "%s"
						schema_type = "string"
					}
				}`,
				resourceType, resourceName, environment.Name,
				environment.ProjectId, environment.LatestDeploymentLog.BlueprintId,
				environment.LatestDeploymentLog.BlueprintRevision, configurationVariables.Name,
				configurationVariables.Value,
			)

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
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariables.Value),
							resource.TestCheckResourceAttr(accessor, "configuration.0.schema_type", "string"),
							resource.TestCheckNoResourceAttr(accessor, "configuration.0.schema_enum"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)

				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariables}, nil) // read after create -> on update -> read after update
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		// Tests use-cases where the response returned by the backend varies from the order of the state.
		t.Run("Create unordered configuration variables", func(t *testing.T) {
			environment := client.Environment{
				Id:        environment.Id,
				Name:      "my-environment",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId:       "template-id",
					BlueprintRevision: "revision",
				},
			}

			configurationVariable1 := client.ConfigurationVariable{
				Value: "my env var value",
				Name:  "my env var",
				Schema: &client.ConfigurationVariableSchema{
					Type: "string",
				},
			}

			configurationVariable2 := client.ConfigurationVariable{
				Value: "a",
				Name:  "b",
				Schema: &client.ConfigurationVariableSchema{
					Type: "string",
				},
			}

			environmentResource := fmt.Sprintf(`
				resource "%s" "%s" {
					name = "%s"
					project_id = "%s"
					template_id = "%s"
					revision = "%s"
					force_destroy = true
					configuration {
						name = "%s"
						value = "%s"
					}
					configuration {
						name = "%s"
						value = "%s"
					}
				}`,
				resourceType, resourceName, environment.Name,
				environment.ProjectId, environment.LatestDeploymentLog.BlueprintId,
				environment.LatestDeploymentLog.BlueprintRevision, configurationVariable1.Name,
				configurationVariable1.Value, configurationVariable2.Name,
				configurationVariable2.Value,
			)

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
							resource.TestCheckResourceAttr(accessor, "configuration.0.name", configurationVariable1.Name),
							resource.TestCheckResourceAttr(accessor, "configuration.0.value", configurationVariable1.Value),
							resource.TestCheckResourceAttr(accessor, "configuration.1.name", configurationVariable2.Name),
							resource.TestCheckResourceAttr(accessor, "configuration.1.value", configurationVariable2.Value),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)

				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariable2, configurationVariable1}, nil) // read after create -> on update -> read after update
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Update to: revision, configuration should trigger a deployment", func(t *testing.T) {
			environment := client.Environment{
				Id:                     environment.Id,
				Name:                   "my-environment",
				ProjectId:              "project-id",
				AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
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
				Type:   "string",
				Enum:   []string{"a", "b"},
				Format: client.Text,
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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                   environment.Name,
					ProjectId:              environment.ProjectId,
					AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
					DeployRequest: &client.DeployRequest{
						BlueprintId:       environment.LatestDeploymentLog.BlueprintId,
						BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
					},
				}).Times(1).Return(environment, nil)

				mock.EXPECT().EnvironmentDeploy(environment.Id, gomock.Any()).Times(1).Return(client.EnvironmentDeployResponse{
					Id: "deployment-id",
				}, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(2).Return(client.ConfigurationChanges{}, nil) // read after create -> on update -> read after update
				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil), // 1 after create, 1 before update
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariables}, nil),
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
					mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariables}, nil),
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

	}

	testTTL := func() {
		t.Run("TTL update", func(t *testing.T) {
			environment := client.Environment{
				Id:            environment.Id,
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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentUpdateTTL(environment.Id, client.TTL{
					Type:  client.TTLTypeDate,
					Value: updatedEnvironment.LifespanEndAt,
				}).Times(1).Return(environment, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil) // read after create -> on update -> read after update

				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Deleting TTL from environment should update ttl to infinite", func(t *testing.T) {
			environment := client.Environment{
				Id:            environment.Id,
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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentUpdateTTL(environment.Id, client.TTL{
					Type:  client.TTlTypeInfinite,
					Value: "",
				}).Times(1).Return(environment, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, updatedEnvironment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)

				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
				)

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})
	}

	testTriggers := func() {
		t.Run("Deleting triggers should set them to false", func(t *testing.T) {
			falsey := false
			truthyFruity := true
			environment := client.Environment{
				Id:        environment.Id,
				Name:      "my-environment",
				ProjectId: "project-id",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: "template-id",
				},

				AutoDeployOnPathChangesOnly: &falsey,
				ContinuousDeployment:        &falsey,
				RequiresApproval:            &falsey,
				PullRequestPlanDeployments:  &falsey,
			}
			environmentAfterUpdate := client.Environment{
				Id:        environment.Id,
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: environment.LatestDeploymentLog.BlueprintId,
				},

				ContinuousDeployment:       &truthyFruity,
				RequiresApproval:           &truthyFruity,
				PullRequestPlanDeployments: &truthyFruity,
				AutoDeployByCustomGlob:     ".*",
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
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
							resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", "true"),
							resource.TestCheckResourceAttr(accessor, "deploy_on_push", "false"),
							resource.TestCheckResourceAttr(accessor, "run_plan_on_pull_requests", "false"),
							resource.TestCheckResourceAttr(accessor, "auto_deploy_on_path_changes_only", "false"),
							resource.TestCheckResourceAttr(accessor, "auto_deploy_by_custom_glob", environment.AutoDeployByCustomGlob),
						),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                             environment.Name,
							"project_id":                       environment.ProjectId,
							"template_id":                      environment.LatestDeploymentLog.BlueprintId,
							"deploy_on_push":                   true,
							"run_plan_on_pull_requests":        true,
							"force_destroy":                    true,
							"auto_deploy_on_path_changes_only": true,
							"auto_deploy_by_custom_glob":       ".*",
						}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
							resource.TestCheckResourceAttr(accessor, "template_id", environment.LatestDeploymentLog.BlueprintId),
							resource.TestCheckResourceAttr(accessor, "approve_plan_automatically", "false"),
							resource.TestCheckResourceAttr(accessor, "deploy_on_push", "true"),
							resource.TestCheckResourceAttr(accessor, "run_plan_on_pull_requests", "true"),
							resource.TestCheckResourceAttr(accessor, "auto_deploy_on_path_changes_only", "true"),
							resource.TestCheckResourceAttr(accessor, "auto_deploy_by_custom_glob", ".*"),
						),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentUpdate(environment.Id, client.EnvironmentUpdate{
					Name:                        environment.Name,
					ContinuousDeployment:        &truthyFruity,
					AutoDeployOnPathChangesOnly: &truthyFruity,
					RequiresApproval:            &truthyFruity,
					PullRequestPlanDeployments:  &truthyFruity,
					AutoDeployByCustomGlob:      ".*",
				}).Times(1).Return(environmentAfterUpdate, nil)

				gomock.InOrder(
					mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),            // 1 after create, 1 before update
					mock.EXPECT().Environment(gomock.Any()).Times(1).Return(environmentAfterUpdate, nil), // 1 after update
				)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(3).Return(client.ConfigurationChanges{}, nil)
				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

	}

	testForceDestroy := func() {
		t.Run("should only allow destroy when force destroy is enabled", func(t *testing.T) {
			environment := client.Environment{
				Id:            environment.Id,
				Name:          "my-environment",
				ProjectId:     "project-id",
				WorkspaceName: "workspace-name",
				LatestDeploymentLog: client.DeploymentLog{
					BlueprintId: templateId,
				},
			}

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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(5).Return(client.ConfigurationChanges{}, nil)
				mock.EXPECT().Environment(gomock.Any()).Times(5).Return(environment, nil)
				mock.EXPECT().EnvironmentDestroy(gomock.Any()).Times(1)

			})
		})
	}

	testValidationFailures := func() {
		t.Run("create environment with is_remote_apply_enabled set to 'true'", func(t *testing.T) {
			runUnitTest(t, resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                    environment.Name,
							"project_id":              environment.ProjectId,
							"template_id":             environment.LatestDeploymentLog.BlueprintId,
							"is_remote_apply_enabled": true,
						}),
						ExpectError: regexp.MustCompile("is_remote_apply_enabled cannot be set when creating a new environment"),
					},
				},
			}, func(mockFunc *client.MockApiClientInterface) {})
		})

		t.Run("Failure in validation while glob is enabled and pathChanges no", func(t *testing.T) {
			autoDeployWithCustomGlobEnabled := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                             environment.Name,
							"project_id":                       environment.ProjectId,
							"template_id":                      environment.LatestDeploymentLog.BlueprintId,
							"auto_deploy_on_path_changes_only": false,
							"force_destroy":                    true,
							"run_plan_on_pull_requests":        true,
							"auto_deploy_by_custom_glob":       "/**",
						}),
						ExpectError: regexp.MustCompile("cannot set auto_deploy_by_custom_glob when auto_deploy_on_path_changes_only is disabled"),
					},
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                             environment.Name,
							"project_id":                       environment.ProjectId,
							"template_id":                      environment.LatestDeploymentLog.BlueprintId,
							"auto_deploy_on_path_changes_only": true,
							"run_plan_on_pull_requests":        true,
							"auto_deploy_by_custom_glob":       "/**",
							"force_destroy":                    true,
						}),
						ExpectNonEmptyPlan: true,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", environment.Id),
							resource.TestCheckResourceAttr(accessor, "name", environment.Name),
							resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						),
					},
				},
			}
			runUnitTest(t, autoDeployWithCustomGlobEnabled, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().ConfigurationVariablesByScope(gomock.Any(), gomock.Any()).Times(1).Return(client.ConfigurationChanges{}, nil)
				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
			})
		})

		t.Run("Failure in validation while prPlan and CD are disabled", func(t *testing.T) {
			autoDeployWithCustomGlobEnabled := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"name":                       environment.Name,
							"project_id":                 environment.ProjectId,
							"template_id":                environment.LatestDeploymentLog.BlueprintId,
							"force_destroy":              true,
							"auto_deploy_by_custom_glob": "/**",
						}),
						ExpectError: regexp.MustCompile("Missing required argument"),
					},
				},
			}
			runUnitTest(t, autoDeployWithCustomGlobEnabled, func(mock *client.MockApiClientInterface) {})
		})
	}

	testApiFailures := func() {
		t.Run("Failure template not assigned to project", func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      createEnvironmentResourceConfig(environment),
						ExpectError: regexp.MustCompile("could not create environment: template is not assigned to project"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(client.Template{
					ProjectId:  "no-match",
					ProjectIds: []string{"no-match"},
				}, nil)
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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                   environment.Name,
					ProjectId:              environment.ProjectId,
					WorkspaceName:          environment.WorkspaceName,
					AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendFalse,
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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                   environment.Name,
					ProjectId:              environment.ProjectId,
					WorkspaceName:          environment.WorkspaceName,
					AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendFalse,
				}).Times(1).Return(environment, nil)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(2).Return(client.ConfigurationChanges{}, nil)
				mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
					Name:                       updatedEnvironment.Name,
					AutoDeployByCustomGlob:     autoDeployByCustomGlobDefault,
					TerragruntWorkingDirectory: updatedEnvironment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           updatedEnvironment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendTrue,
					IsArchived:                 boolPtr(true),
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
							"name":                         updatedEnvironment.Name,
							"project_id":                   updatedEnvironment.ProjectId,
							"template_id":                  updatedEnvironment.LatestDeploymentLog.BlueprintId,
							"revision":                     updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
							"force_destroy":                true,
							"terragrunt_working_directory": environment.TerragruntWorkingDirectory,
							"vcs_commands_alias":           environment.VcsCommandsAlias,
						}),
						ExpectError: regexp.MustCompile("failed deploying environment: error"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
				mock.EXPECT().EnvironmentDeploy(updatedEnvironment.Id, client.DeployRequest{
					BlueprintId:       updatedEnvironment.LatestDeploymentLog.BlueprintId,
					BlueprintRevision: updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
				}).Times(1).Return(client.EnvironmentDeployResponse{}, errors.New("error"))
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil) // 1 after create, 1 before update
				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(2).Return(client.ConfigurationChanges{}, nil)
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
				mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
				mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
					Name:                   environment.Name,
					ProjectId:              environment.ProjectId,
					WorkspaceName:          environment.WorkspaceName,
					AutoDeployByCustomGlob: autoDeployByCustomGlobDefault,
					DeployRequest: &client.DeployRequest{
						BlueprintId: templateId,
					},
					TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
					VcsCommandsAlias:           environment.VcsCommandsAlias,
					IsRemoteBackend:            &isRemoteBackendFalse,
				}).Times(1).Return(environment, nil)
				mock.EXPECT().Environment(gomock.Any()).Return(client.Environment{}, errors.New("error"))
				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(0).Return(client.ConfigurationChanges{}, nil)
			})

		})
	}

	t.Run("Failure in delete", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createEnvironmentResourceConfig(environment),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(environment.LatestDeploymentLog.BlueprintId).Times(1).Return(template, nil)
			mock.EXPECT().EnvironmentCreate(gomock.Any()).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentDeploy(updatedEnvironment.Id, client.DeployRequest{
				BlueprintId:       updatedEnvironment.LatestDeploymentLog.BlueprintId,
				BlueprintRevision: updatedEnvironment.LatestDeploymentLog.BlueprintRevision,
			}).Times(1).Return(client.EnvironmentDeployResponse{}, errors.New("error"))
			mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil)
			mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil)
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1).Return(environment, http.NewMockFailedResponseError(400))
		})
	})

	testSuccess()
	testTTL()
	testTriggers()
	testForceDestroy()
	testValidationFailures()
	testApiFailures()
}

func TestUnitEnvironmentWithoutTemplateResource(t *testing.T) {
	resourceType := "env0_environment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	environment := client.Environment{
		Id:                         "id0",
		Name:                       "my-environment",
		ProjectId:                  "project-id",
		WorkspaceName:              "workspace-name",
		TerragruntWorkingDirectory: "/terragrunt/directory/",
		VcsCommandsAlias:           "alias",
		LatestDeploymentLog: client.DeploymentLog{
			BlueprintId: "id-template-0",
		},
	}

	environmentWithBluePrint := environment
	environmentWithBluePrint.BlueprintId = "id-template-0"

	template := client.Template{
		Id:          "id-template-0",
		Name:        "single-use-template-for-" + environment.Name,
		Description: "description0",
		Repository:  "env0/repo",
		Path:        "path/zero",
		Revision:    "branch-zero",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      2,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      1,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:                 "terraform",
		GithubInstallationId: 1,
		TerraformVersion:     "0.12.24",
	}

	updatedTemplate := client.Template{
		Id:          "id-template-0",
		Name:        "single-use-template-for-" + environment.Name,
		Description: "description1",
		Repository:  "env0/repo1",
		Path:        "path/zero1",
		Revision:    "branch-zero1",
		Retry: client.TemplateRetry{
			OnDeploy: &client.TemplateRetryOn{
				Times:      3,
				ErrorRegex: "RetryMeForDeploy.*",
			},
			OnDestroy: &client.TemplateRetryOn{
				Times:      3,
				ErrorRegex: "RetryMeForDestroy.*",
			},
		},
		Type:                 "terragrunt",
		GithubInstallationId: 2,
		TerraformVersion:     "0.12.25",
		TerragruntVersion:    "0.26.1",
		TerragruntTfBinary:   "terraform",
	}

	environmentCreatePayload := client.EnvironmentCreate{
		Name:                       environment.Name,
		ProjectId:                  environment.ProjectId,
		DeployRequest:              &client.DeployRequest{},
		WorkspaceName:              environment.WorkspaceName,
		RequiresApproval:           environment.RequiresApproval,
		ContinuousDeployment:       environment.ContinuousDeployment,
		PullRequestPlanDeployments: environment.PullRequestPlanDeployments,
		TerragruntWorkingDirectory: environment.TerragruntWorkingDirectory,
		VcsCommandsAlias:           environment.VcsCommandsAlias,
	}

	templateCreatePayload := client.TemplateCreatePayload{
		Repository:           template.Repository,
		Description:          template.Description,
		GithubInstallationId: template.GithubInstallationId,
		IsGitlabEnterprise:   template.IsGitlabEnterprise,
		IsGitLab:             template.TokenId != "",
		TokenId:              template.TokenId,
		Path:                 template.Path,
		Revision:             template.Revision,
		Type:                 "terraform",
		Retry:                template.Retry,
		TerraformVersion:     template.TerraformVersion,
		BitbucketClientKey:   template.BitbucketClientKey,
		IsGithubEnterprise:   template.IsGithubEnterprise,
		IsBitbucketServer:    template.IsBitbucketServer,
		FileName:             template.FileName,
		TerragruntVersion:    template.TerragruntVersion,
		IsTerragruntRunAll:   template.IsTerragruntRunAll,
		OrganizationId:       template.OrganizationId,
	}

	templateUpdatePayload := client.TemplateCreatePayload{
		Repository:           updatedTemplate.Repository,
		Description:          updatedTemplate.Description,
		GithubInstallationId: updatedTemplate.GithubInstallationId,
		IsGitlabEnterprise:   updatedTemplate.IsGitlabEnterprise,
		IsGitLab:             updatedTemplate.TokenId != "",
		TokenId:              updatedTemplate.TokenId,
		Path:                 updatedTemplate.Path,
		Revision:             updatedTemplate.Revision,
		Type:                 "terragrunt",
		Retry:                updatedTemplate.Retry,
		TerraformVersion:     updatedTemplate.TerraformVersion,
		BitbucketClientKey:   updatedTemplate.BitbucketClientKey,
		IsGithubEnterprise:   updatedTemplate.IsGithubEnterprise,
		IsBitbucketServer:    updatedTemplate.IsBitbucketServer,
		FileName:             updatedTemplate.FileName,
		TerragruntVersion:    updatedTemplate.TerragruntVersion,
		IsTerragruntRunAll:   updatedTemplate.IsTerragruntRunAll,
		OrganizationId:       updatedTemplate.OrganizationId,
		TerragruntTfBinary:   updatedTemplate.TerragruntTfBinary,
	}

	createPayload := client.EnvironmentCreateWithoutTemplate{
		EnvironmentCreate: environmentCreatePayload,
		TemplateCreate:    templateCreatePayload,
	}

	createEnvironmentResourceConfig := func(environment client.Environment, template client.Template) string {
		terragruntVersion := ""
		if template.TerragruntVersion != "" {
			terragruntVersion = "terragrunt_version = \"" + template.TerragruntVersion + "\""
		}

		terragruntTfBinary := ""
		if template.TerragruntTfBinary != "" {
			terragruntTfBinary = "terragrunt_tf_binary = \"" + template.TerragruntTfBinary + "\""
		}

		return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "%s"
			project_id = "%s"
			workspace = "%s"
			terragrunt_working_directory = "%s"
			force_destroy = true
			vcs_commands_alias = "%s"
			without_template_settings {
				repository = "%s"
				terraform_version = "%s"
				type = "%s"
				revision = "%s"
				path = "%s"
				retries_on_deploy = %d
				retry_on_deploy_only_when_matches_regex = "%s"
				retries_on_destroy = %d
				retry_on_destroy_only_when_matches_regex = "%s"
				description = "%s"
				github_installation_id = %d
				%s
				%s
			}
		}`,
			resourceType, resourceName,
			environment.Name,
			environment.ProjectId,
			environment.WorkspaceName,
			environment.TerragruntWorkingDirectory,
			environment.VcsCommandsAlias,
			template.Repository,
			template.TerraformVersion,
			template.Type,
			template.Revision,
			template.Path,
			template.Retry.OnDeploy.Times,
			template.Retry.OnDeploy.ErrorRegex,
			template.Retry.OnDestroy.Times,
			template.Retry.OnDestroy.ErrorRegex,
			template.Description,
			template.GithubInstallationId,
			terragruntVersion,
			terragruntTfBinary,
		)
	}

	t.Run("Success in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				// Create the environment and template
				{
					Config: createEnvironmentResourceConfig(environment, template),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckNoResourceAttr(accessor, "template_id"),
						resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
						resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", environment.TerragruntWorkingDirectory),
						resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", environment.VcsCommandsAlias),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.repository", template.Repository),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.terraform_version", template.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.type", template.Type),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.path", template.Path),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.revision", template.Revision),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retries_on_deploy", strconv.Itoa(template.Retry.OnDeploy.Times)),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retries_on_destroy", strconv.Itoa(template.Retry.OnDestroy.Times)),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex),
					),
				},

				// Update the template.
				{
					Config: createEnvironmentResourceConfig(environment, updatedTemplate),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckNoResourceAttr(accessor, "template_id"),
						resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
						resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", environment.TerragruntWorkingDirectory),
						resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", environment.VcsCommandsAlias),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.repository", updatedTemplate.Repository),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.terraform_version", updatedTemplate.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.type", updatedTemplate.Type),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.path", updatedTemplate.Path),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.terragrunt_version", updatedTemplate.TerragruntVersion),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.terragrunt_tf_binary", updatedTemplate.TerragruntTfBinary),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.revision", updatedTemplate.Revision),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retries_on_deploy", strconv.Itoa(updatedTemplate.Retry.OnDeploy.Times)),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retry_on_deploy_only_when_matches_regex", updatedTemplate.Retry.OnDeploy.ErrorRegex),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retries_on_destroy", strconv.Itoa(updatedTemplate.Retry.OnDestroy.Times)),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retry_on_destroy_only_when_matches_regex", updatedTemplate.Retry.OnDestroy.ErrorRegex),
					),
				},
				// No need to update template
				{
					Config: createEnvironmentResourceConfig(environment, updatedTemplate),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckNoResourceAttr(accessor, "template_id"),
						resource.TestCheckResourceAttr(accessor, "workspace", environment.WorkspaceName),
						resource.TestCheckResourceAttr(accessor, "terragrunt_working_directory", environment.TerragruntWorkingDirectory),
						resource.TestCheckResourceAttr(accessor, "vcs_commands_alias", environment.VcsCommandsAlias),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.repository", updatedTemplate.Repository),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.terraform_version", updatedTemplate.TerraformVersion),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.type", updatedTemplate.Type),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.path", updatedTemplate.Path),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.revision", updatedTemplate.Revision),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retries_on_deploy", strconv.Itoa(updatedTemplate.Retry.OnDeploy.Times)),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retry_on_deploy_only_when_matches_regex", updatedTemplate.Retry.OnDeploy.ErrorRegex),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retries_on_destroy", strconv.Itoa(updatedTemplate.Retry.OnDestroy.Times)),
						resource.TestCheckResourceAttr(accessor, "without_template_settings.0.retry_on_destroy_only_when_matches_regex", updatedTemplate.Retry.OnDestroy.ErrorRegex),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {

			gomock.InOrder(
				// Step1
				// Create
				mock.EXPECT().EnvironmentCreateWithoutTemplate(createPayload).Times(1).Return(environmentWithBluePrint, nil),

				// Read
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(template, nil),

				// Step2
				// Read
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(template, nil),

				// Update
				mock.EXPECT().TemplateUpdate(template.Id, templateUpdatePayload).Times(1).Return(updatedTemplate, nil),

				// Read
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(updatedTemplate, nil),

				// Step3
				// Read
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(updatedTemplate, nil),

				// Read
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeEnvironment, environment.Id).Times(1).Return(client.ConfigurationChanges{}, nil),
				mock.EXPECT().Template(template.Id).Times(1).Return(updatedTemplate, nil),

				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1),
			)
		})
	})

	t.Run("Revision conflict", func(t *testing.T) {
		createEnvironmentResourceConfigWithRevision := func(environment client.Environment, template client.Template, revision string) string {
			return fmt.Sprintf(`
			resource "%s" "%s" {
				name = "%s"
				revision = "%s"
				project_id = "%s"
				workspace = "%s"
				terragrunt_working_directory = "%s"
				force_destroy = true
				vcs_commands_alias = "%s"
				without_template_settings {
					repository = "%s"
					terraform_version = "%s"
					type = "%s"
					revision = "%s"
					path = "%s"
					retries_on_deploy = %d
					retry_on_deploy_only_when_matches_regex = "%s"
					retries_on_destroy = %d
					retry_on_destroy_only_when_matches_regex = "%s"
					description = "%s"
					github_installation_id = %d
				}
			}`,
				resourceType, resourceName,
				environment.Name,
				revision,
				environment.ProjectId,
				environment.WorkspaceName,
				environment.TerragruntWorkingDirectory,
				environment.VcsCommandsAlias,
				template.Repository,
				template.TerraformVersion,
				template.Type,
				template.Revision,
				template.Path,
				template.Retry.OnDeploy.Times,
				template.Retry.OnDeploy.ErrorRegex,
				template.Retry.OnDestroy.Times,
				template.Retry.OnDestroy.ErrorRegex,
				template.Description,
				template.GithubInstallationId,
			)
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createEnvironmentResourceConfigWithRevision(environment, template, "environment_revision"),
					ExpectError: regexp.MustCompile("conflicts with without_template_settings"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}

func TestUnitEnvironmentWithSubEnvironment(t *testing.T) {
	resourceType := "env0_environment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	workflowSubEnvironment := client.WorkflowSubEnvironment{
		EnvironmentId: "subenv1id",
	}

	subEnvironment := SubEnvironment{
		Alias:     "alias1",
		Revision:  "revision1",
		Workspace: "workspace1",
		Configuration: client.ConfigurationChanges{
			{
				Name:        "name",
				Value:       "value",
				Scope:       client.ScopeEnvironment,
				IsSensitive: boolPtr(false),
				IsReadOnly:  boolPtr(false),
				IsRequired:  boolPtr(false),
				Schema: &client.ConfigurationVariableSchema{
					Type: "string",
				},
				Type: (*client.ConfigurationVariableType)(intPtr(0)),
			},
		},
	}

	updatedSubEnvironment := subEnvironment
	updatedSubEnvironment.Configuration = append(updatedSubEnvironment.Configuration, client.ConfigurationVariable{
		Name:        "name2",
		Value:       "value2",
		Scope:       client.ScopeEnvironment,
		IsSensitive: boolPtr(false),
		IsReadOnly:  boolPtr(false),
		IsRequired:  boolPtr(false),
		Schema: &client.ConfigurationVariableSchema{
			Type: "string",
		},
		Type: (*client.ConfigurationVariableType)(intPtr(0)),
	})

	subEnvrionmentWithId := subEnvironment
	subEnvrionmentWithId.Id = workflowSubEnvironment.EnvironmentId

	environment := client.Environment{
		Id:          "id",
		Name:        "environment",
		ProjectId:   "project-id",
		BlueprintId: "template-id",

		LatestDeploymentLog: client.DeploymentLog{
			WorkflowFile: &client.WorkflowFile{
				Environments: map[string]client.WorkflowSubEnvironment{
					subEnvironment.Alias: workflowSubEnvironment,
				},
			},
		},
	}

	configurationVariable := client.ConfigurationVariable{
		Value: "v1",
		Name:  "n1",
		Type:  (*client.ConfigurationVariableType)(intPtr(0)),
		Schema: &client.ConfigurationVariableSchema{
			Type: "string",
		},
	}

	environmentCreatePayload := client.EnvironmentCreate{
		Name:      environment.Name,
		ProjectId: environment.ProjectId,
		ConfigurationChanges: &client.ConfigurationChanges{
			{
				Name:        "n1",
				Value:       "v1",
				Scope:       client.ScopeDeployment,
				IsSensitive: boolPtr(false),
				IsReadOnly:  boolPtr(false),
				IsRequired:  boolPtr(false),
				Schema: &client.ConfigurationVariableSchema{
					Type: "string",
				},
				Type: (*client.ConfigurationVariableType)(intPtr(0)),
			},
		},
		DeployRequest: &client.DeployRequest{
			BlueprintId: environment.BlueprintId,
			SubEnvironments: map[string]client.SubEnvironment{
				subEnvironment.Alias: {
					Revision:             subEnvironment.Revision,
					ConfigurationChanges: subEnvironment.Configuration,
				},
			},
		},
		Type: "workflow",
	}

	template := client.Template{
		ProjectId: environment.ProjectId,
	}

	deployRequest := client.DeployRequest{
		BlueprintId:       environment.BlueprintId,
		BlueprintRevision: environment.LatestDeploymentLog.BlueprintRevision,
		SubEnvironments: map[string]client.SubEnvironment{
			subEnvironment.Alias: {
				Revision:             subEnvironment.Revision,
				Workspace:            subEnvironment.Workspace,
				ConfigurationChanges: updatedSubEnvironment.Configuration,
			},
		},
	}

	t.Run("Success in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
					resource "%s" "%s" {
						name = "%s"
						project_id = "%s"
						template_id = "%s"
						force_destroy = true
						configuration {
							name = "n1"
							value = "v1"
						}
						sub_environment_configuration {
							alias = "%s"
							revision = "%s"
							workspace = "%s"
							configuration {
								name = "%s"
								value = "%s"
							}
						}
					}`,
						resourceType, resourceName,
						environmentCreatePayload.Name,
						environmentCreatePayload.ProjectId,
						environment.BlueprintId,
						subEnvironment.Alias,
						subEnvironment.Revision,
						subEnvironment.Workspace,
						subEnvironment.Configuration[0].Name,
						subEnvironment.Configuration[0].Value,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.id", workflowSubEnvironment.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.alias", subEnvironment.Alias),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.revision", subEnvironment.Revision),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.workspace", subEnvironment.Workspace),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.configuration.0.name", subEnvironment.Configuration[0].Name),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.configuration.0.value", subEnvironment.Configuration[0].Value),
					),
				},
				{
					Config: fmt.Sprintf(`
					resource "%s" "%s" {
						name = "%s"
						project_id = "%s"
						template_id = "%s"
						force_destroy = true
						sub_environment_configuration {
							alias = "%s"
							revision = "%s"
							workspace = "%s"
							configuration {
								name = "%s"
								value = "%s"
							}
							configuration {
								name = "%s"
								value = "%s"
							}
						}
					}`,
						resourceType, resourceName,
						environmentCreatePayload.Name,
						environmentCreatePayload.ProjectId,
						environment.BlueprintId,
						subEnvironment.Alias,
						subEnvironment.Revision,
						subEnvironment.Workspace,
						subEnvironment.Configuration[0].Name,
						subEnvironment.Configuration[0].Value,
						updatedSubEnvironment.Configuration[1].Name,
						updatedSubEnvironment.Configuration[1].Value,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.BlueprintId),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.id", workflowSubEnvironment.EnvironmentId),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.alias", subEnvironment.Alias),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.revision", subEnvironment.Revision),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.workspace", subEnvironment.Workspace),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.configuration.0.name", subEnvironment.Configuration[0].Name),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.configuration.0.value", subEnvironment.Configuration[0].Value),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.configuration.1.name", updatedSubEnvironment.Configuration[1].Name),
						resource.TestCheckResourceAttr(accessor, "sub_environment_configuration.0.configuration.1.value", updatedSubEnvironment.Configuration[1].Value),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Template(environmentCreatePayload.DeployRequest.BlueprintId).Times(1).Return(template, nil),
				mock.EXPECT().EnvironmentCreate(environmentCreatePayload).Times(1).Return(environment, nil),
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeWorkflow, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariable}, nil),
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeWorkflow, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariable}, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeWorkflow, subEnvrionmentWithId.Id).Times(1).Return(subEnvironment.Configuration, nil),
				mock.EXPECT().EnvironmentDeploy(environment.Id, deployRequest).Times(1).Return(client.EnvironmentDeployResponse{
					Id: environment.Id,
				}, nil),
				mock.EXPECT().Environment(environment.Id).Times(1).Return(environment, nil),
				mock.EXPECT().ConfigurationVariablesByScope(client.ScopeWorkflow, environment.Id).Times(1).Return(client.ConfigurationChanges{configurationVariable}, nil),
				mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1),
			)
		})
	})
}
