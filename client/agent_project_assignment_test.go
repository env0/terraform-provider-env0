package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Agent Project Assignment", func() {
	mapping := map[string]string{
		"pid1": "aid1",
		"pid2": "aid2",
	}

	imapping := map[string]interface{}{
		"pid1": "aid1",
		"pid2": "aid2",
	}

	expectedResponse := ProjectsAgentsAssignments{
		DefaultAgent:   "default-agent",
		ProjectsAgents: imapping,
	}

	errorMock := errors.New("error")

	Describe("assign agents to projects", func() {

		Describe("success", func() {
			var results *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)

				httpCall = mockHttpClient.EXPECT().
					Post("/agents/projects-assignments?organizationId="+organizationId, mapping, gomock.Any()).Times(1).
					Do(func(path string, request interface{}, response *ProjectsAgentsAssignments) {
						*response = expectedResponse
					})
				results, err = apiClient.AssignAgentsToProjects(mapping)

			})

			It("should return results", func() {
				Expect(*results).To(Equal(expectedResponse))
			})

			It("should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("failure", func() {
			var results *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)

				httpCall = mockHttpClient.EXPECT().
					Post("/agents/projects-assignments?organizationId="+organizationId, mapping, gomock.Any()).Times(1).Return(errorMock)

				results, err = apiClient.AssignAgentsToProjects(mapping)
			})

			It("should return an error", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("should not return results", func() {
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("ProjectsAgentsAssignments", func() {
		Describe("success", func() {
			var results *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)

				httpCall = mockHttpClient.EXPECT().
					Get("/agents/projects-assignments", map[string]string{"organizationId": organizationId}, gomock.Any()).Times(1).
					Do(func(path string, request interface{}, response *ProjectsAgentsAssignments) {
						*response = expectedResponse
					})
				results, err = apiClient.ProjectsAgentsAssignments()
			})

			It("should return results", func() {
				Expect(*results).To(Equal(expectedResponse))
			})

			It("should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("failure", func() {
			var results *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)

				httpCall = mockHttpClient.EXPECT().
					Get("/agents/projects-assignments", map[string]string{"organizationId": organizationId}, gomock.Any()).Times(1).
					Return(errorMock)

				results, err = apiClient.ProjectsAgentsAssignments()
			})

			It("should return an error", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("should not return results", func() {
				Expect(results).To(BeNil())
			})
		})
	})
})
