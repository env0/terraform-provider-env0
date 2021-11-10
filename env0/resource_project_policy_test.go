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
		NumberOfEnvironments:       1,
		NumberOfEnvironmentsTotal:  1,
		RequiresApprovalDefault:    false,
		IncludeCostEstimation:      false,
		SkipApplyWhenPlanIsEmpty:   false,
		DisableDestroyEnvironments: false,
		SkipRedundantDepolyments:   false,
		UpdatedBy:                  "updater0",
	}

	updatedPolicy := client.Policy{
		Id:                         policy.Id,
		ProjectId:                  policy.ProjectId,
		NumberOfEnvironments:       0,
		NumberOfEnvironmentsTotal:  0,
		RequiresApprovalDefault:    false,
		IncludeCostEstimation:      false,
		SkipApplyWhenPlanIsEmpty:   false,
		DisableDestroyEnvironments: false,
		SkipRedundantDepolyments:   false,
		UpdatedBy:                  "updater0",
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"project_id":                    policy.ProjectId,
					"number_of_environments":        policy.NumberOfEnvironments,
					"number_of_environments_total":  policy.NumberOfEnvironmentsTotal,
					"requires_approval_default":     policy.RequiresApprovalDefault,
					"include_cost_estimation":       policy.IncludeCostEstimation,
					"skip_apply_when_plan_is_empty": policy.SkipApplyWhenPlanIsEmpty,
					"disable_destroy_environments":  policy.DisableDestroyEnvironments,
					"skip_redundant_deployments":    policy.SkipRedundantDepolyments,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "project_id", policy.ProjectId),
					resource.TestCheckResourceAttr(accessor, "number_of_environments", strconv.Itoa(policy.NumberOfEnvironments)),
					resource.TestCheckResourceAttr(accessor, "number_of_environments_total", strconv.Itoa(policy.NumberOfEnvironmentsTotal)),
					resource.TestCheckResourceAttr(accessor, "requires_approval_default", strconv.FormatBool(policy.RequiresApprovalDefault)),
					resource.TestCheckResourceAttr(accessor, "include_cost_estimation", strconv.FormatBool(policy.IncludeCostEstimation)),
					resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", strconv.FormatBool(policy.SkipApplyWhenPlanIsEmpty)),
					resource.TestCheckResourceAttr(accessor, "disable_destroy_environments", strconv.FormatBool(policy.DisableDestroyEnvironments)),
					resource.TestCheckResourceAttr(accessor, "skip_redundant_deployments", strconv.FormatBool(policy.SkipRedundantDepolyments)),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"project_id":                    updatedPolicy.ProjectId,
					"number_of_environments":        updatedPolicy.NumberOfEnvironments,
					"number_of_environments_total":  updatedPolicy.NumberOfEnvironmentsTotal,
					"requires_approval_default":     updatedPolicy.RequiresApprovalDefault,
					"include_cost_estimation":       updatedPolicy.IncludeCostEstimation,
					"skip_apply_when_plan_is_empty": updatedPolicy.SkipApplyWhenPlanIsEmpty,
					"disable_destroy_environments":  updatedPolicy.DisableDestroyEnvironments,
					"skip_redundant_deployments":    updatedPolicy.SkipRedundantDepolyments,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "project_id", updatedPolicy.ProjectId),
					resource.TestCheckResourceAttr(accessor, "number_of_environments", strconv.Itoa(updatedPolicy.NumberOfEnvironments)),
					resource.TestCheckResourceAttr(accessor, "number_of_environments_total", strconv.Itoa(updatedPolicy.NumberOfEnvironmentsTotal)),
					resource.TestCheckResourceAttr(accessor, "requires_approval_default", strconv.FormatBool(updatedPolicy.RequiresApprovalDefault)),
					resource.TestCheckResourceAttr(accessor, "include_cost_estimation", strconv.FormatBool(updatedPolicy.IncludeCostEstimation)),
					resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", strconv.FormatBool(updatedPolicy.SkipApplyWhenPlanIsEmpty)),
					resource.TestCheckResourceAttr(accessor, "disable_destroy_environments", strconv.FormatBool(updatedPolicy.DisableDestroyEnvironments)),
					resource.TestCheckResourceAttr(accessor, "skip_redundant_deployments", strconv.FormatBool(updatedPolicy.SkipRedundantDepolyments)),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
			ProjectId:                  policy.ProjectId,
			NumberOfEnvironments:       policy.NumberOfEnvironments,
			NumberOfEnvironmentsTotal:  policy.NumberOfEnvironmentsTotal,
			RequiresApprovalDefault:    policy.RequiresApprovalDefault,
			IncludeCostEstimation:      policy.IncludeCostEstimation,
			SkipApplyWhenPlanIsEmpty:   policy.SkipApplyWhenPlanIsEmpty,
			DisableDestroyEnvironments: policy.DisableDestroyEnvironments,
			SkipRedundantDepolyments:   policy.SkipRedundantDepolyments,
		}).Times(1).Return(policy, nil)
		mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
			ProjectId:                  updatedPolicy.ProjectId,
			NumberOfEnvironments:       updatedPolicy.NumberOfEnvironments,
			NumberOfEnvironmentsTotal:  updatedPolicy.NumberOfEnvironmentsTotal,
			RequiresApprovalDefault:    updatedPolicy.RequiresApprovalDefault,
			IncludeCostEstimation:      updatedPolicy.IncludeCostEstimation,
			SkipApplyWhenPlanIsEmpty:   updatedPolicy.SkipApplyWhenPlanIsEmpty,
			DisableDestroyEnvironments: updatedPolicy.DisableDestroyEnvironments,
			SkipRedundantDepolyments:   updatedPolicy.SkipRedundantDepolyments,
		}).Times(2).Return(updatedPolicy, nil)

		gomock.InOrder(
			mock.EXPECT().Policy(gomock.Any()).Times(2).Return(policy, nil),        // 1 after create, 1 before update
			mock.EXPECT().Policy(gomock.Any()).Times(1).Return(updatedPolicy, nil), // 1 after update
		)
	})
}

func TestUnitPolicyInvalidParams(t *testing.T) {
	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate("env0_project_policy", "test", map[string]interface{}{"project_id": ""}),
				ExpectError: regexp.MustCompile("project id not ok"),
			},
		},
	}

	runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
}
