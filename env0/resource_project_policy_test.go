package env0

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitPolicyResource(t *testing.T) {
	resourceType := "env0_project_policy"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	policy := client.Policy{
		Id:                         "id0",
		ProjectId:                  "project0",
		NumberOfEnvironments:       intPtr(1),
		NumberOfEnvironmentsTotal:  intPtr(1),
		IncludeCostEstimation:      true,
		SkipApplyWhenPlanIsEmpty:   true,
		DisableDestroyEnvironments: true,
		SkipRedundantDeployments:   true,
		UpdatedBy:                  "updater0",
		MaxTtl:                     stringPtr("inherit"),
		DefaultTtl:                 stringPtr("inherit"),
		ForceRemoteBackend:         true,
	}

	updatedPolicy := client.Policy{
		Id:                         policy.Id,
		ProjectId:                  policy.ProjectId,
		NumberOfEnvironments:       nil,
		NumberOfEnvironmentsTotal:  nil,
		RequiresApprovalDefault:    false,
		IncludeCostEstimation:      false,
		SkipApplyWhenPlanIsEmpty:   false,
		DisableDestroyEnvironments: false,
		SkipRedundantDeployments:   false,
		UpdatedBy:                  "updater0",
		MaxTtl:                     nil,
		DefaultTtl:                 stringPtr("7-h"),
		ForceRemoteBackend:         false,
	}

	resetPolicy := client.Policy{
		ProjectId:               policy.ProjectId,
		RequiresApprovalDefault: true,
	}

	t.Run("create and update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":                    policy.ProjectId,
						"number_of_environments":        *policy.NumberOfEnvironments,
						"number_of_environments_total":  *policy.NumberOfEnvironmentsTotal,
						"requires_approval_default":     policy.RequiresApprovalDefault,
						"include_cost_estimation":       policy.IncludeCostEstimation,
						"skip_apply_when_plan_is_empty": policy.SkipApplyWhenPlanIsEmpty,
						"disable_destroy_environments":  policy.DisableDestroyEnvironments,
						"skip_redundant_deployments":    policy.SkipRedundantDeployments,
						"run_pull_request_plan_default": policy.RunPullRequestPlanDefault,
						"continuous_deployment_default": policy.ContinuousDeploymentDefault,
						"force_remote_backend":          policy.ForceRemoteBackend,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", policy.ProjectId),
						resource.TestCheckResourceAttr(accessor, "number_of_environments", strconv.Itoa(*policy.NumberOfEnvironments)),
						resource.TestCheckResourceAttr(accessor, "number_of_environments_total", strconv.Itoa(*policy.NumberOfEnvironmentsTotal)),
						resource.TestCheckResourceAttr(accessor, "requires_approval_default", strconv.FormatBool(policy.RequiresApprovalDefault)),
						resource.TestCheckResourceAttr(accessor, "include_cost_estimation", strconv.FormatBool(policy.IncludeCostEstimation)),
						resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", strconv.FormatBool(policy.SkipApplyWhenPlanIsEmpty)),
						resource.TestCheckResourceAttr(accessor, "disable_destroy_environments", strconv.FormatBool(policy.DisableDestroyEnvironments)),
						resource.TestCheckResourceAttr(accessor, "skip_redundant_deployments", strconv.FormatBool(policy.SkipRedundantDeployments)),
						resource.TestCheckResourceAttr(accessor, "run_pull_request_plan_default", strconv.FormatBool(policy.RunPullRequestPlanDefault)),
						resource.TestCheckResourceAttr(accessor, "continuous_deployment_default", strconv.FormatBool(policy.ContinuousDeploymentDefault)),
						resource.TestCheckResourceAttr(accessor, "max_ttl", "inherit"),
						resource.TestCheckResourceAttr(accessor, "default_ttl", "inherit"),
						resource.TestCheckResourceAttr(accessor, "force_remote_backend", strconv.FormatBool(policy.ForceRemoteBackend)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":                    updatedPolicy.ProjectId,
						"requires_approval_default":     updatedPolicy.RequiresApprovalDefault,
						"include_cost_estimation":       updatedPolicy.IncludeCostEstimation,
						"skip_apply_when_plan_is_empty": updatedPolicy.SkipApplyWhenPlanIsEmpty,
						"disable_destroy_environments":  updatedPolicy.DisableDestroyEnvironments,
						"skip_redundant_deployments":    updatedPolicy.SkipRedundantDeployments,
						"max_ttl":                       "Infinite",
						"default_ttl":                   *updatedPolicy.DefaultTtl,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", updatedPolicy.ProjectId),
						resource.TestCheckResourceAttr(accessor, "number_of_environments", "0"),
						resource.TestCheckResourceAttr(accessor, "number_of_environments_total", "0"),
						resource.TestCheckResourceAttr(accessor, "requires_approval_default", strconv.FormatBool(updatedPolicy.RequiresApprovalDefault)),
						resource.TestCheckResourceAttr(accessor, "include_cost_estimation", strconv.FormatBool(updatedPolicy.IncludeCostEstimation)),
						resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", strconv.FormatBool(updatedPolicy.SkipApplyWhenPlanIsEmpty)),
						resource.TestCheckResourceAttr(accessor, "disable_destroy_environments", strconv.FormatBool(updatedPolicy.DisableDestroyEnvironments)),
						resource.TestCheckResourceAttr(accessor, "skip_redundant_deployments", strconv.FormatBool(updatedPolicy.SkipRedundantDeployments)),
						resource.TestCheckResourceAttr(accessor, "max_ttl", "Infinite"),
						resource.TestCheckResourceAttr(accessor, "default_ttl", *updatedPolicy.DefaultTtl),
						resource.TestCheckResourceAttr(accessor, "force_remote_backend", strconv.FormatBool(updatedPolicy.ForceRemoteBackend)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				// Create
				mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
					ProjectId:                  policy.ProjectId,
					NumberOfEnvironments:       *policy.NumberOfEnvironments,
					NumberOfEnvironmentsTotal:  *policy.NumberOfEnvironmentsTotal,
					RequiresApprovalDefault:    policy.RequiresApprovalDefault,
					IncludeCostEstimation:      policy.IncludeCostEstimation,
					SkipApplyWhenPlanIsEmpty:   policy.SkipApplyWhenPlanIsEmpty,
					DisableDestroyEnvironments: policy.DisableDestroyEnvironments,
					SkipRedundantDeployments:   policy.SkipRedundantDeployments,
					MaxTtl:                     "inherit",
					DefaultTtl:                 "inherit",
					ForceRemoteBackend:         true,
				}).Times(1).Return(policy, nil),
				mock.EXPECT().Policy(gomock.Any()).Times(2).Return(policy, nil), // 1 after create, 1 before update
				// Update
				mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
					ProjectId:                  updatedPolicy.ProjectId,
					NumberOfEnvironments:       0,
					NumberOfEnvironmentsTotal:  0,
					RequiresApprovalDefault:    updatedPolicy.RequiresApprovalDefault,
					IncludeCostEstimation:      updatedPolicy.IncludeCostEstimation,
					SkipApplyWhenPlanIsEmpty:   updatedPolicy.SkipApplyWhenPlanIsEmpty,
					DisableDestroyEnvironments: updatedPolicy.DisableDestroyEnvironments,
					SkipRedundantDeployments:   updatedPolicy.SkipRedundantDeployments,
					MaxTtl:                     "",
					DefaultTtl:                 *updatedPolicy.DefaultTtl,
					ForceRemoteBackend:         updatedPolicy.ForceRemoteBackend,
				}).Times(1).Return(updatedPolicy, nil),
				mock.EXPECT().Policy(gomock.Any()).Times(1).Return(updatedPolicy, nil), // 1 after update.
				// Delete
				mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
					ProjectId:                  resetPolicy.ProjectId,
					NumberOfEnvironments:       0,
					NumberOfEnvironmentsTotal:  0,
					RequiresApprovalDefault:    resetPolicy.RequiresApprovalDefault,
					IncludeCostEstimation:      resetPolicy.IncludeCostEstimation,
					SkipApplyWhenPlanIsEmpty:   resetPolicy.SkipApplyWhenPlanIsEmpty,
					DisableDestroyEnvironments: resetPolicy.DisableDestroyEnvironments,
					SkipRedundantDeployments:   resetPolicy.SkipRedundantDeployments,
					ForceRemoteBackend:         resetPolicy.ForceRemoteBackend,
				}).Times(1).Return(resetPolicy, nil),
			)
		})
	})

	t.Run("requires_approval_default default is true", func(t *testing.T) {
		expectedPolicy := client.Policy{
			Id:                         "id0",
			ProjectId:                  "project0",
			NumberOfEnvironments:       intPtr(1),
			NumberOfEnvironmentsTotal:  intPtr(1),
			RequiresApprovalDefault:    true,
			IncludeCostEstimation:      true,
			SkipApplyWhenPlanIsEmpty:   true,
			DisableDestroyEnvironments: true,
			SkipRedundantDeployments:   true,
			UpdatedBy:                  "updater0",
		}

		testCaseForDefault := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":                    policy.ProjectId,
						"number_of_environments":        *policy.NumberOfEnvironments,
						"number_of_environments_total":  *policy.NumberOfEnvironmentsTotal,
						"include_cost_estimation":       policy.IncludeCostEstimation,
						"skip_apply_when_plan_is_empty": policy.SkipApplyWhenPlanIsEmpty,
						"disable_destroy_environments":  policy.DisableDestroyEnvironments,
						"skip_redundant_deployments":    policy.SkipRedundantDeployments,
						"run_pull_request_plan_default": policy.RunPullRequestPlanDefault,
						"continuous_deployment_default": policy.ContinuousDeploymentDefault,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", policy.ProjectId),
						resource.TestCheckResourceAttr(accessor, "number_of_environments", strconv.Itoa(*policy.NumberOfEnvironments)),
						resource.TestCheckResourceAttr(accessor, "number_of_environments_total", strconv.Itoa(*policy.NumberOfEnvironmentsTotal)),
						resource.TestCheckResourceAttr(accessor, "requires_approval_default", strconv.FormatBool(true)),
						resource.TestCheckResourceAttr(accessor, "include_cost_estimation", strconv.FormatBool(policy.IncludeCostEstimation)),
						resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", strconv.FormatBool(policy.SkipApplyWhenPlanIsEmpty)),
						resource.TestCheckResourceAttr(accessor, "disable_destroy_environments", strconv.FormatBool(policy.DisableDestroyEnvironments)),
						resource.TestCheckResourceAttr(accessor, "skip_redundant_deployments", strconv.FormatBool(policy.SkipRedundantDeployments)),
						resource.TestCheckResourceAttr(accessor, "run_pull_request_plan_default", strconv.FormatBool(policy.RunPullRequestPlanDefault)),
						resource.TestCheckResourceAttr(accessor, "continuous_deployment_default", strconv.FormatBool(policy.ContinuousDeploymentDefault)),
					),
				},
			},
		}

		runUnitTest(t, testCaseForDefault, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
				ProjectId:                  policy.ProjectId,
				NumberOfEnvironments:       *policy.NumberOfEnvironments,
				NumberOfEnvironmentsTotal:  *policy.NumberOfEnvironmentsTotal,
				RequiresApprovalDefault:    true,
				IncludeCostEstimation:      policy.IncludeCostEstimation,
				SkipApplyWhenPlanIsEmpty:   policy.SkipApplyWhenPlanIsEmpty,
				DisableDestroyEnvironments: policy.DisableDestroyEnvironments,
				SkipRedundantDeployments:   policy.SkipRedundantDeployments,
				DefaultTtl:                 "inherit",
				MaxTtl:                     "inherit",
			}).Times(1).Return(policy, nil)

			mock.EXPECT().Policy(gomock.Any()).Times(1).Return(expectedPolicy, nil)

			mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
				ProjectId:                  resetPolicy.ProjectId,
				NumberOfEnvironments:       0,
				NumberOfEnvironmentsTotal:  0,
				RequiresApprovalDefault:    resetPolicy.RequiresApprovalDefault,
				IncludeCostEstimation:      resetPolicy.IncludeCostEstimation,
				SkipApplyWhenPlanIsEmpty:   resetPolicy.SkipApplyWhenPlanIsEmpty,
				DisableDestroyEnvironments: resetPolicy.DisableDestroyEnvironments,
				SkipRedundantDeployments:   resetPolicy.SkipRedundantDeployments,
			}).Times(1).Return(resetPolicy, nil)

		})
	})

	t.Run("Create Failure - max smaller than default", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":                    policy.ProjectId,
						"number_of_environments":        *policy.NumberOfEnvironments,
						"number_of_environments_total":  *policy.NumberOfEnvironmentsTotal,
						"include_cost_estimation":       policy.IncludeCostEstimation,
						"skip_apply_when_plan_is_empty": policy.SkipApplyWhenPlanIsEmpty,
						"disable_destroy_environments":  policy.DisableDestroyEnvironments,
						"skip_redundant_deployments":    policy.SkipRedundantDeployments,
						"run_pull_request_plan_default": policy.RunPullRequestPlanDefault,
						"continuous_deployment_default": policy.ContinuousDeploymentDefault,
						"max_ttl":                       "12-h",
						"default_ttl":                   "1-d",
					}),
					ExpectError: regexp.MustCompile("default ttl must not be larger than max ttl"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - max smaller than default (Infinite)", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":                    policy.ProjectId,
						"number_of_environments":        *policy.NumberOfEnvironments,
						"number_of_environments_total":  *policy.NumberOfEnvironmentsTotal,
						"include_cost_estimation":       policy.IncludeCostEstimation,
						"skip_apply_when_plan_is_empty": policy.SkipApplyWhenPlanIsEmpty,
						"disable_destroy_environments":  policy.DisableDestroyEnvironments,
						"skip_redundant_deployments":    policy.SkipRedundantDeployments,
						"run_pull_request_plan_default": policy.RunPullRequestPlanDefault,
						"continuous_deployment_default": policy.ContinuousDeploymentDefault,
						"max_ttl":                       "1-M",
						"default_ttl":                   "Infinite",
					}),
					ExpectError: regexp.MustCompile("default ttl must not be larger than max ttl"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - only one is inherit", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":                    policy.ProjectId,
						"number_of_environments":        *policy.NumberOfEnvironments,
						"number_of_environments_total":  *policy.NumberOfEnvironmentsTotal,
						"include_cost_estimation":       policy.IncludeCostEstimation,
						"skip_apply_when_plan_is_empty": policy.SkipApplyWhenPlanIsEmpty,
						"disable_destroy_environments":  policy.DisableDestroyEnvironments,
						"skip_redundant_deployments":    policy.SkipRedundantDeployments,
						"run_pull_request_plan_default": policy.RunPullRequestPlanDefault,
						"continuous_deployment_default": policy.ContinuousDeploymentDefault,
						"default_ttl":                   "1-d",
					}),
					ExpectError: regexp.MustCompile("max_ttl and default_ttl must both inherit organization settings or override them"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}

func TestUnitPolicyInvalidParams(t *testing.T) {
	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate("env0_project_policy", "test", map[string]interface{}{"project_id": ""}),
				ExpectError: regexp.MustCompile("may not be empty"),
			},
		},
	}

	runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
}
