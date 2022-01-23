package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestWorkflowTriggerDataSource(t *testing.T) {
	trigger := client.WorkflowTrigger{
		Id: "id0",
	}

	environmentId := "environment_id"
	resourceType := "env0_workflow_triggers"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	t.Run("By environment id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id": environmentId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "downstream_environment_ids.0", trigger.Id),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().WorkflowTrigger(environmentId).AnyTimes().Return([]client.WorkflowTrigger{trigger}, nil)
		})
	})

	t.Run("When Error", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id": environmentId,
					}),
					ExpectError: regexp.MustCompile("could not get workflow triggers: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().WorkflowTrigger(environmentId).AnyTimes().Return([]client.WorkflowTrigger{}, errors.New("error"))
		})
	})
}
