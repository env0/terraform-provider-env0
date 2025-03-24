package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User Environment Assignment", func() {
	environmentId := "environmentId"
	userId := "userId"

	assignPayload := &AssignUserRoleToEnvironmentPayload{
		UserId:        userId,
		EnvironmentId: environmentId,
		Role:          "role1",
	}

	expectedResponse := &UserRoleEnvironmentAssignment{
		Id: "id",
	}

	errorMock := errors.New("error")

	Describe("AssignUserToEnvironment", func() {

		Describe("Successful", func() {
			var actualResult *UserRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/users", assignPayload, gomock.Any()).
					Do(func(path string, request any, response *UserRoleEnvironmentAssignment) {
						*response = *expectedResponse
					}).Times(1)
				actualResult, err = apiClient.AssignUserRoleToEnvironment(assignPayload)

			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("should return the PUT result", func() {
				Expect(*actualResult).To(Equal(*expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult *UserRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/users", gomock.Any(), gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.AssignUserRoleToEnvironment(assignPayload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("RemoveUserFromEnvironment", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/users", map[string]string{"environmentId": environmentId, "userId": userId})
			err = apiClient.RemoveUserRoleFromEnvironment(environmentId, userId)
		})

		It("Should send DELETE request with assignment id", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Describe("UserEnvironmentAssignments", func() {

		Describe("Successful", func() {
			var actualResult []UserRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/roles/assignments/users", map[string]string{"environmentId": environmentId}, gomock.Any()).
					Do(func(path string, request any, response *[]UserRoleEnvironmentAssignment) {
						*response = []UserRoleEnvironmentAssignment{*expectedResponse}
					}).Times(1)
				actualResult, err = apiClient.UserRoleEnvironmentAssignments(environmentId)

			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("should return the GET result", func() {
				Expect(actualResult).To(Equal([]UserRoleEnvironmentAssignment{*expectedResponse}))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult []UserRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/roles/assignments/users", map[string]string{"environmentId": environmentId}, gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.UserRoleEnvironmentAssignments(environmentId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})
})
