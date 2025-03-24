package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitEnvironmentDriftDetectionResource(t *testing.T) {
	environmentId := "environment0"
	resourceType := "env0_environment_drift_detection"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	drift := client.EnvironmentSchedulingExpression{
		Cron:                 "2 * * * *",
		Enabled:              true,
		AutoDriftRemediation: "CODE_TO_CLOUD",
	}
	updateDrift := client.EnvironmentSchedulingExpression{
		Cron:                 "2 2 * * *",
		Enabled:              true,
		AutoDriftRemediation: "DISABLED",
	}
	autoDriftRemediation := "CODE_TO_CLOUD"

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"environment_id":         environmentId,
						"cron":                   drift.Cron,
						"auto_drift_remediation": autoDriftRemediation,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "cron", drift.Cron),
						resource.TestCheckResourceAttr(accessor, "auto_drift_remediation", autoDriftRemediation),
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"environment_id":         environmentId,
						"cron":                   drift.Cron,
						"auto_drift_remediation": autoDriftRemediation,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "cron", drift.Cron),
						resource.TestCheckResourceAttr(accessor, "auto_drift_remediation", autoDriftRemediation),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"environment_id":         environmentId,
						"cron":                   updateDrift.Cron,
						"auto_drift_remediation": "DISABLED",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "cron", updateDrift.Cron),
						resource.TestCheckResourceAttr(accessor, "auto_drift_remediation", "DISABLED"),
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
