package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitWorkflowTriggerResource(t *testing.T) {
	resourceType := "env0_workflow_trigger"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	environmentId := "environment_id"
	triggerId := "trigger_environment_id"
	otherTriggerId := "other_trigger_environment_id"

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": triggerId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId+"_"+triggerId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "downstream_environment_id", triggerId),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": otherTriggerId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId+"_"+otherTriggerId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "downstream_environment_id", otherTriggerId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().SubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
				mock.EXPECT().WorkflowTrigger(environmentId).Times(2).Return([]client.WorkflowTrigger{
					{
						Id: triggerId,
					},
				}, nil),
				mock.EXPECT().UnsubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
				mock.EXPECT().SubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{otherTriggerId},
				}).Times(1).Return(nil),
				mock.EXPECT().WorkflowTrigger(environmentId).Times(1).Return([]client.WorkflowTrigger{
					{
						Id: otherTriggerId,
					},
				}, nil),
				mock.EXPECT().UnsubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{otherTriggerId},
				}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Failure in Get Triggers", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": triggerId,
					}),
					ExpectError: regexp.MustCompile("could not get workflow triggers: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().SubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
				mock.EXPECT().WorkflowTrigger(environmentId).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().UnsubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Failure in Unsubscribe", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": triggerId,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": otherTriggerId,
					}),
					ExpectError: regexp.MustCompile("failed to unsubscribe a workflow trigger: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().SubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
				mock.EXPECT().WorkflowTrigger(environmentId).Times(2).Return([]client.WorkflowTrigger{
					{
						Id: triggerId,
					},
				}, nil),
				mock.EXPECT().UnsubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(errors.New("error")),
				mock.EXPECT().UnsubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
			)
		})
	})

	t.Run("Failure in Subscribe", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": triggerId,
					}),
					ExpectError: regexp.MustCompile("failed to subscribe a workflow trigger: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().SubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(errors.New("error")),
			)
		})
	})

	t.Run("Drift", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"environment_id":            environmentId,
						"downstream_environment_id": triggerId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentId+"_"+triggerId),
						resource.TestCheckResourceAttr(accessor, "environment_id", environmentId),
						resource.TestCheckResourceAttr(accessor, "downstream_environment_id", triggerId),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().SubscribeWorkflowTrigger(environmentId, client.WorkflowTriggerEnvironments{
					DownstreamEnvironmentIds: []string{triggerId},
				}).Times(1).Return(nil),
				mock.EXPECT().WorkflowTrigger(environmentId).Times(2).Return([]client.WorkflowTrigger{
					{
						Id: otherTriggerId,
					},
				}, nil),
			)
		})
	})
}
