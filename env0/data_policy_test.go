package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestPolicyDataSource(t *testing.T) {
	policy := client.Policy{
		Id:                             "id0",
		ProjectId:                      "project0",
		NumberOfEnvironments:           1,
		NumberOfEnvironmentsPerProject: 2,
		RequiresApprovalDefault:        true,
		IncludeCostEstimation:          true,
		SkipApplyWhenPlanIsEmpty:       true,
		DisableDestroyEnvironments:     true,
	}

	resourceType := "env0_policy"
	resourceName := "test_policy"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": policy.ProjectId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", policy.Id),
						resource.TestCheckResourceAttr(accessor, "project_id", policy.ProjectId),
						resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", "true"),
					),
				},
			},
		}
	}

	t.Run("valid", func(t *testing.T) {
		runUnitTest(t, getValidTestCase(map[string]interface{}{}), func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Policy(policy.ProjectId).AnyTimes().Return(policy, nil)
		})
	})

}
