package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitEnvironmentSchedulingResource(t *testing.T) {
	environmentId := "environment0"
	resourceType := "env0_environment_scheduling"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	environmentScheduling := client.EnvironmentScheduling{
		Deploy:  &client.EnvironmentSchedulingExpression{Cron: "1 * * * *", Enabled: true},
		Destroy: &client.EnvironmentSchedulingExpression{Cron: "2 * * * *", Enabled: true},
	}

	cronExprKeys := []string{"deploy_cron", "destroy_cron"}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"environment_id": environmentId,
						"deploy_cron":    environmentScheduling.Deploy.Cron,
						"destroy_cron":   environmentScheduling.Destroy.Cron,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "deploy_cron", environmentScheduling.Deploy.Cron),
						resource.TestCheckResourceAttr(accessor, "destroy_cron", environmentScheduling.Destroy.Cron),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentSchedulingUpdate(environmentId, environmentScheduling).Times(1).Return(environmentScheduling, nil)
			mock.EXPECT().EnvironmentScheduling(environmentId).Times(1).Return(environmentScheduling, nil)
			mock.EXPECT().EnvironmentSchedulingDelete(environmentId).Times(1).Return(nil)
		})
	})

	for _, key := range cronExprKeys {
		t.Run("Failure due to invalid cron expression for "+key, func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
							"environment_id": environmentId,
							key:              "not_a_valid_cron",
						}),
						ExpectError: regexp.MustCompile("Invalid cron expression"),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
		})
	}

	t.Run("Failure when both deploy_cron and destroy_cron attributes are missing", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"environment_id": environmentId,
					}),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}
