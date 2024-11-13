package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Agent Project Assignment", func() {
	projectId := "projectId"
	userId := "userId"

	assignPayload := &AssignUserToProjectPayload{
		UserId: userId,
		Role:   string(AdminRole),
	}

	updatePayload := &UpdateUserProjectAssignmentPayload{
		Role: string(AdminRole),
	}

	expectedResponse := &UserProjectAssignment{
		Id: "id",
	}

	errorMock := errors.New("error")

	Describe("AssignUserToProject", func() {

		Describe("Successful", func() {
			var actualResult *UserProjectAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Post("/permissions/projects/"+projectId, assignPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *UserProjectAssignment) {
						*response = *expectedResponse
					}).Times(1)
				actualResult, err = apiClient.AssignUserToProject(projectId, assignPayload)

			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("should return the POST result", func() {
				Expect(*actualResult).To(Equal(*expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult *UserProjectAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Post("/permissions/projects/"+projectId, gomock.Any(), gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.AssignUserToProject(projectId, assignPayload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("RemoveUserFromProject", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/permissions/projects/"+projectId+"/users/"+expectedResponse.Id, nil)
			err = apiClient.RemoveUserFromProject(projectId, expectedResponse.Id)
		})

		It("Should send DELETE request with assignment id", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).Should(BeNil())
		})
	})

	Describe("UserProjectAssignments", func() {

		Describe("Successful", func() {
			var actualResult []UserProjectAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/permissions/projects/"+projectId, nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]UserProjectAssignment) {
						*response = []UserProjectAssignment{*expectedResponse}
					}).Times(1)
				actualResult, err = apiClient.UserProjectAssignments(projectId)

			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("should return the GET result", func() {
				Expect(actualResult).To(Equal([]UserProjectAssignment{*expectedResponse}))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult []UserProjectAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/permissions/projects/"+projectId, nil, gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.UserProjectAssignments(projectId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("UpdateUserProjectAssignment", func() {

		Describe("Successful", func() {
			var actualResult *UserProjectAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/permissions/projects/"+projectId+"/users/"+userId, updatePayload, gomock.Any()).
					Do(func(path string, request interface{}, response *UserProjectAssignment) {
						*response = *expectedResponse
					}).Times(1)
				actualResult, err = apiClient.UpdateUserProjectAssignment(projectId, userId, updatePayload)

			})

			It("Should send PUT request with params", func() {
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
			var actualResult *UserProjectAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/permissions/projects/"+projectId+"/users/"+userId, gomock.Any(), gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.UpdateUserProjectAssignment(projectId, userId, updatePayload)
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
