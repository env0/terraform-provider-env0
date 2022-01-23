package env0

import (
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestUnitWorkflowTriggerResource(t *testing.T) {
	resourceType := "env0_workflow_triggers"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	environmentId := "environment_id"
	trigger := client.WorkflowTrigger{
		Id: "id0",
	}

	otherTrigger := client.WorkflowTrigger{
		Id: "id1",
	}
	createHCL := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"environment_id":             environmentId,
		"downstream_environment_ids": []string{trigger.Id},
	})
	t.Run("Success", func(t *testing.T) {

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createHCL,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "downstream_environment_ids.0", trigger.Id),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":             environmentId,
						"downstream_environment_ids": []string{otherTrigger.Id},
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "downstream_environment_ids.0", otherTrigger.Id),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().WorkflowTriggerUpsert(environmentId, client.WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{trigger.Id},
			}).Times(1).Return([]client.WorkflowTrigger{trigger}, nil)
			mock.EXPECT().WorkflowTriggerUpsert(environmentId, client.WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{otherTrigger.Id},
			}).Times(1).Return([]client.WorkflowTrigger{otherTrigger}, nil)
			mock.EXPECT().WorkflowTriggerUpsert(environmentId, client.WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{},
			}).Times(1).Return([]client.WorkflowTrigger{}, nil)

			gomock.InOrder(
				mock.EXPECT().WorkflowTrigger(environmentId).Times(2).Return([]client.WorkflowTrigger{trigger}, nil),      // 1 after createHCL, 1 before update
				mock.EXPECT().WorkflowTrigger(environmentId).Times(1).Return([]client.WorkflowTrigger{otherTrigger}, nil), // 1 after update
			)

		})
	})

	t.Run("Failure in upsert", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createHCL,
					ExpectError: regexp.MustCompile("could not Create or Update workflow triggers: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().WorkflowTriggerUpsert(environmentId, client.WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{trigger.Id},
			}).Times(1).Return([]client.WorkflowTrigger{}, errors.New("error"))
		})

	})

	t.Run("Failure in read", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createHCL,
					ExpectError: regexp.MustCompile("could not get workflow triggers: error "),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().WorkflowTriggerUpsert(environmentId, client.WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{trigger.Id},
			}).Times(1).Return([]client.WorkflowTrigger{trigger}, nil)
			mock.EXPECT().WorkflowTriggerUpsert(environmentId, client.WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{},
			}).Times(1).Return([]client.WorkflowTrigger{}, nil)

			mock.EXPECT().WorkflowTrigger(environmentId).Return([]client.WorkflowTrigger{}, errors.New("error"))
		})

	})
}
