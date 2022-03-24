package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitNotificationProjectAssignmentResource(t *testing.T) {
	resourceType := "env0_notification_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	projectId := "pid1"

	assignment := client.NotificationProjectAssignment{
		Id:                     "id1",
		NotificationEndpointId: "nid1",
		EventNames: []string{
			"driftDetected",
		},
	}

	assignmentEventNamesUpdated := client.NotificationProjectAssignment{
		Id:                     "id1",
		NotificationEndpointId: "nid1",
		EventNames: []string{
			"driftDetected", "environmentDestroyStarted",
		},
	}

	assignmentNotificationEndpointIdUpdated := client.NotificationProjectAssignment{
		Id:                     "id2",
		NotificationEndpointId: "nid2",
		EventNames: []string{
			"driftDetected",
		},
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":               projectId,
						"notification_endpoint_id": assignment.NotificationEndpointId,
						"event_names":              []string{"driftDetected"},
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.Id),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "notification_endpoint_id", assignment.NotificationEndpointId),
						resource.TestCheckResourceAttr(accessor, "event_names.#", "1"),
						resource.TestCheckResourceAttr(accessor, "event_names.0", "driftDetected"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":               projectId,
						"notification_endpoint_id": assignment.NotificationEndpointId,
						"event_names":              []string{"driftDetected", "environmentDestroyStarted"},
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.Id),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "notification_endpoint_id", assignment.NotificationEndpointId),
						resource.TestCheckResourceAttr(accessor, "event_names.#", "2"),
						resource.TestCheckResourceAttr(accessor, "event_names.0", "driftDetected"),
						resource.TestCheckResourceAttr(accessor, "event_names.1", "environmentDestroyStarted"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignment.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: assignment.EventNames,
			}).Times(1).Return(&assignment, nil)

			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignment.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: assignmentEventNamesUpdated.EventNames,
			}).Times(1).Return(&assignmentEventNamesUpdated, nil)

			gomock.InOrder(
				mock.EXPECT().NotificationProjectAssignments(projectId).Times(2).Return([]client.NotificationProjectAssignment{assignment}, nil),
				mock.EXPECT().NotificationProjectAssignments(projectId).Times(1).Return([]client.NotificationProjectAssignment{assignmentEventNamesUpdated}, nil),
			)

			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignment.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: []string{},
			}).Times(1)
		})
	})

	t.Run("Success - notification endpoint id updated", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":               projectId,
						"notification_endpoint_id": assignment.NotificationEndpointId,
						"event_names":              []string{"driftDetected"},
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignment.Id),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "notification_endpoint_id", assignment.NotificationEndpointId),
						resource.TestCheckResourceAttr(accessor, "event_names.#", "1"),
						resource.TestCheckResourceAttr(accessor, "event_names.0", "driftDetected"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":               projectId,
						"notification_endpoint_id": assignmentNotificationEndpointIdUpdated.NotificationEndpointId,
						"event_names":              []string{"driftDetected"},
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", assignmentNotificationEndpointIdUpdated.Id),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "notification_endpoint_id", assignmentNotificationEndpointIdUpdated.NotificationEndpointId),
						resource.TestCheckResourceAttr(accessor, "event_names.#", "1"),
						resource.TestCheckResourceAttr(accessor, "event_names.0", "driftDetected"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignment.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: assignment.EventNames,
			}).Times(1).Return(&assignment, nil)

			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignmentNotificationEndpointIdUpdated.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: assignmentNotificationEndpointIdUpdated.EventNames,
			}).Times(1).Return(&assignmentNotificationEndpointIdUpdated, nil)

			gomock.InOrder(
				mock.EXPECT().NotificationProjectAssignments(projectId).Times(2).Return([]client.NotificationProjectAssignment{assignment}, nil),
				mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignment.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
					EventNames: []string{},
				}).Times(1),
				mock.EXPECT().NotificationProjectAssignments(projectId).Times(1).Return([]client.NotificationProjectAssignment{assignment, assignmentNotificationEndpointIdUpdated}, nil),
			)

			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignmentNotificationEndpointIdUpdated.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: []string{},
			}).Times(1)
		})
	})

	t.Run("Create Failure - invalid event name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":               projectId,
						"notification_endpoint_id": assignment.NotificationEndpointId,
						"event_names":              []string{"driftDetected", "invalidname"},
					}),
					ExpectError: regexp.MustCompile("'invalidname' must be one of:"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":               projectId,
						"notification_endpoint_id": assignment.NotificationEndpointId,
						"event_names":              []string{"driftDetected"},
					}),
					ExpectError: regexp.MustCompile("could not create or update notification project assignment"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().NotificationProjectAssignmentUpdate(projectId, assignment.NotificationEndpointId, client.NotificationProjectAssignmentUpdatePayload{
				EventNames: []string{"driftDetected"},
			}).Times(1).Return(nil, errors.New("error"))
		})
	})
}
