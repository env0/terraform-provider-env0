package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	gomock "github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notification Client", func() {
	mockNotification := Notification{
		Id:             "notification-id",
		Name:           "notification-name",
		Type:           "Slack",
		Value:          "https://some.url.com",
		OrganizationId: "organization-id",
	}

	Describe("NotificationCreate", func() {
		Describe("Success", func() {
			var createdNotification *Notification
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				createNotificationPayload := NotificationCreate{}
				copier.Copy(&createNotificationPayload, &mockNotification)

				expectedCreateRequest := NotificationCreateWithOrganizationId{
					NotificationCreate: createNotificationPayload,
					OrganizationId:     organizationId,
				}

				httpCall = mockHttpClient.EXPECT().
					Post("/notifications/endpoints", expectedCreateRequest, gomock.Any()).
					Do(func(path string, request interface{}, response *Notification) {
						*response = mockNotification
					})

				createdNotification, err = apiClient.NotificationCreate(createNotificationPayload)
			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return created notification", func() {
				Expect(*createdNotification).To(Equal(mockNotification))
			})
		})
	})

	Describe("NotificationDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/notifications/endpoints/" + mockNotification.Id)
			apiClient.NotificationDelete(mockNotification.Id)
		})

		It("Should send DELETE request with notification id", func() {
			httpCall.Times(1)
		})
	})

	Describe("NotificationUpdate", func() {
		Describe("Success", func() {
			updateMockNotification := mockNotification
			updateMockNotification.Name = "updated-name"
			var updatedNotification *Notification
			var err error

			BeforeEach(func() {
				updateNotificationPayload := NotificationUpdate{Name: "updated-name"}

				httpCall = mockHttpClient.EXPECT().
					Patch("/notifications/endpoints/"+mockNotification.Id, updateNotificationPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *Notification) {
						*response = updateMockNotification
					})

				updatedNotification, err = apiClient.NotificationUpdate(mockNotification.Id, updateNotificationPayload)
			})

			It("Should send Patch request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return notification received from API", func() {
				Expect(*updatedNotification).To(Equal(updateMockNotification))
			})
		})
	})
})
