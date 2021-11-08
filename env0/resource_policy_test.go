package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitPolicyResource(t *testing.T) {
	resourceType := "env0_policy"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	policy := client.Policy{
		Id:        "id0",
		ProjectId: "project0",
	}

	updatedPolicy := client.Policy{
		Id:        policy.Id,
		ProjectId: policy.ProjectId,
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"project_id": policy.ProjectId,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "project_id", updatedPolicy.ProjectId),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"project_id":             policy.ProjectId,
					"number_of_environments": fmt.Sprintf("%d", policy.NumberOfEnvironmentsTotal),
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "number_of_environments", fmt.Sprintf("%d", updatedPolicy.NumberOfEnvironmentsTotal)),
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
		}).Times(1).Return(updatedPolicy, nil)

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
				Config:      resourceConfigCreate("env0_policy", "test", map[string]interface{}{"project_id": ""}),
				ExpectError: regexp.MustCompile("Project id cannot be empty"),
			},
		},
	}

	runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
}
