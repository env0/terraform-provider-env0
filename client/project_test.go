package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const projectName = "project_test"

var _ = Describe("Project", func() {
	var project Project
	mockProject := Project{
		Id:             "idX",
		Name:           "projectX",
		OrganizationId: organizationId,
	}

	Describe("ProjectCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Post("/projects", map[string]interface{}{
					"name":           projectName,
					"organizationId": organizationId,
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *Project) {
					*response = mockProject
				})

			project, _ = apiClient.ProjectCreate(projectName)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return project", func() {
			Expect(project).To(Equal(mockProject))
		})
	})

	Describe("ProjectDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/projects/" + mockProject.Id)
			apiClient.ProjectDelete(mockProject.Id)
		})

		It("Should send DELETE request with project id", func() {
			httpCall.Times(1)
		})
	})

	Describe("Project", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/projects/"+mockProject.Id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *Project) {
					*response = mockProject
				})
			project, _ = apiClient.Project(mockProject.Id)
		})

		It("Should send GET request with project id", func() {
			httpCall.Times(1)
		})

		It("Should return project", func() {
			Expect(project).To(Equal(mockProject))
		})
	})

	Describe("Projects", func() {
		var projects []Project
		mockProjects := []Project{mockProject}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/projects", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]Project) {
					*response = mockProjects
				})
			projects, _ = apiClient.Projects()
		})

		It("Should send GET request with organization id param", func() {
			httpCall.Times(1)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should return projects", func() {
			Expect(projects).To(Equal(mockProjects))
		})
	})
})
