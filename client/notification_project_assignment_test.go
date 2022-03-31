package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification Project Assignment Client", func() {
	projectId := "pid"

	mockNotificationProjectAssignment := NotificationProjectAssignment{
		Id:                     "id",
		NotificationEndpointId: "nid",
		EventNames:             []string{"name1", "name2"},
	}

	Describe("Get Notification Project Assignments", func() {
		var returnedAssignments []NotificationProjectAssignment

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/notifications/projects/"+projectId, gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *[]NotificationProjectAssignment) {
					*response = []NotificationProjectAssignment{mockNotificationProjectAssignment}
				})
			returnedAssignments, _ = apiClient.NotificationProjectAssignments(projectId)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return assignment", func() {
			Expect(returnedAssignments).To(Equal([]NotificationProjectAssignment{mockNotificationProjectAssignment}))
		})
	})

	Describe("Update Notification Project Assignment", func() {
		var updatedAssignment *NotificationProjectAssignment
		var err error

		updatedMockNotificationProjectAssignment := mockNotificationProjectAssignment
		updatedMockNotificationProjectAssignment.EventNames = []string{"name3"}

		BeforeEach(func() {
			updateAssignmentPayload := NotificationProjectAssignmentUpdatePayload{
				EventNames: []string{"name3"},
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/notifications/projects/"+projectId+"/endpoints/"+mockNotificationProjectAssignment.Id, updateAssignmentPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *NotificationProjectAssignment) {
					*response = updatedMockNotificationProjectAssignment
				})

			updatedAssignment, err = apiClient.NotificationProjectAssignmentUpdate(projectId, mockNotificationProjectAssignment.Id, updateAssignmentPayload)
		})

		It("Should send Put request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return assignment received from API", func() {
			Expect(*updatedAssignment).To(Equal(updatedMockNotificationProjectAssignment))
		})
	})
})
