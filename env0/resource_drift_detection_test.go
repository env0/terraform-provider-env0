package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitEnvironmentDriftDetectionResource(t *testing.T) {
	environmentId := "environment0"
	resourceType := "env0_environment_drift_detection"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	drift := client.EnvironmentSchedulingExpression{Cron: "2 * * * *", Enabled: true}
	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id": environmentId,
						"cron":           drift.Cron,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "cron", drift.Cron),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentUpdateDriftDetection(environmentId, drift).Times(1).Return(drift, nil)
			mock.EXPECT().EnvironmentDriftDetection(environmentId).Times(1).Return(drift, nil)
			mock.EXPECT().EnvironmentStopDriftDetection(environmentId).Times(1).Return(nil)
		})
	})
	t.Run("When received Enabled = false from BE (drift)", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id": environmentId,
						"cron":           drift.Cron,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentUpdateDriftDetection(environmentId, drift).Times(1).Return(drift, nil)
			mock.EXPECT().EnvironmentDriftDetection(environmentId).Times(1).Return(client.EnvironmentSchedulingExpression{Cron: "2 * * * *", Enabled: false}, nil)
			mock.EXPECT().EnvironmentStopDriftDetection(environmentId).Times(1).Return(nil)
		})
	})

}
