package env0

import (
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestPolicyDataSource(t *testing.T) {
	policy := client.Policy{
		Id:                          "id0",
		ProjectId:                   "project0",
		NumberOfEnvironments:        intPtr(1),
		NumberOfEnvironmentsTotal:   intPtr(2),
		RequiresApprovalDefault:     true,
		IncludeCostEstimation:       true,
		SkipApplyWhenPlanIsEmpty:    true,
		DisableDestroyEnvironments:  true,
		ContinuousDeploymentDefault: true,
		RunPullRequestPlanDefault:   false,
		MaxTtl:                      stringPtr("3h"),
		DefaultTtl:                  stringPtr("1h"),
		ForceRemoteBackend:          true,
		DriftDetectionEnabled:       true,
		AutoDriftRemediation:        "CODE_TO_CLOUD",
	}

	resourceType := "env0_project_policy"
	resourceName := "test_policy"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": policy.ProjectId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", policy.Id),
						resource.TestCheckResourceAttr(accessor, "project_id", policy.ProjectId),
						resource.TestCheckResourceAttr(accessor, "number_of_environments", strconv.Itoa(*policy.NumberOfEnvironments)),
						resource.TestCheckResourceAttr(accessor, "number_of_environments_total", strconv.Itoa(*policy.NumberOfEnvironmentsTotal)),
						resource.TestCheckResourceAttr(accessor, "requires_approval_default", strconv.FormatBool(policy.RequiresApprovalDefault)),
						resource.TestCheckResourceAttr(accessor, "include_cost_estimation", strconv.FormatBool(policy.IncludeCostEstimation)),
						resource.TestCheckResourceAttr(accessor, "skip_apply_when_plan_is_empty", strconv.FormatBool(policy.SkipApplyWhenPlanIsEmpty)),
						resource.TestCheckResourceAttr(accessor, "disable_destroy_environments", strconv.FormatBool(policy.DisableDestroyEnvironments)),
						resource.TestCheckResourceAttr(accessor, "run_pull_request_plan_default", strconv.FormatBool(policy.RunPullRequestPlanDefault)),
						resource.TestCheckResourceAttr(accessor, "continuous_deployment_default", strconv.FormatBool(policy.ContinuousDeploymentDefault)),
						resource.TestCheckResourceAttr(accessor, "max_ttl", *policy.MaxTtl),
						resource.TestCheckResourceAttr(accessor, "default_ttl", *policy.DefaultTtl),
						resource.TestCheckResourceAttr(accessor, "force_remote_backend", strconv.FormatBool(policy.ForceRemoteBackend)),
						resource.TestCheckResourceAttr(accessor, "auto_drift_remediation", policy.AutoDriftRemediation),
					),
				},
			},
		}
	}

	t.Run("valid", func(t *testing.T) {
		runUnitTest(t, getValidTestCase(), func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Policy(policy.ProjectId).AnyTimes().Return(policy, nil)
		})
	})
}
