package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitEnvironmentDriftDetectionResource(t *testing.T) {
	environmentId := "environment0"
	resourceType := "env0_environment_drift_detection"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	drift := client.EnvironmentSchedulingExpression{Cron: "2 * * * *", Enabled: true}
	updateDrift := client.EnvironmentSchedulingExpression{Cron: "2 2 * * *", Enabled: true}
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
	t.Run("Update", func(t *testing.T) {
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
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id": environmentId,
						"cron":           updateDrift.Cron,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "cron", updateDrift.Cron),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentUpdateDriftDetection(environmentId, drift).Return(drift, nil)
			mock.EXPECT().EnvironmentUpdateDriftDetection(environmentId, updateDrift).Return(updateDrift, nil)

			gomock.InOrder(
				mock.EXPECT().EnvironmentDriftDetection(environmentId).Times(2).Return(drift, nil),
				mock.EXPECT().EnvironmentDriftDetection(environmentId).Return(updateDrift, nil),
			)
			mock.EXPECT().EnvironmentStopDriftDetection(environmentId).Return(nil)
		})
	})
}
