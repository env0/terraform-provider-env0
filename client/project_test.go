package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const projectName = "project_test"
const projectDescription = "project description"
const parentProjectId = "parent_project_id"

var _ = Describe("Project", func() {
	var project Project
	mockProject := Project{
		Id:             "idX",
		Name:           "projectX",
		Description:    "descriptionX",
		OrganizationId: organizationId,
	}

	Describe("ProjectCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			payload := struct {
				ProjectCreatePayload
				OrganizationId string `json:"organizationId"`
			}{
				ProjectCreatePayload{
					Name:            projectName,
					Description:     projectDescription,
					ParentProjectId: parentProjectId,
				},
				organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/projects", payload, gomock.Any()).
				Do(func(path string, request interface{}, response *Project) {
					*response = mockProject
				})

			project, _ = apiClient.ProjectCreate(ProjectCreatePayload{
				Name:            projectName,
				Description:     projectDescription,
				ParentProjectId: parentProjectId,
			})
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

	Describe("ProjectUpdate", func() {
		var mockedResponse Project
		BeforeEach(func() {
			payload := ProjectUpdatePayload{
				Name:        "newName",
				Description: "newDesc",
			}

			mockedResponse = mockProject
			mockedResponse.Name = payload.Name
			mockedResponse.Description = payload.Description

			httpCall = mockHttpClient.EXPECT().
				Put("/projects/"+mockProject.Id, payload, gomock.Any()).
				Do(func(path string, request interface{}, response *Project) {
					*response = mockedResponse
				})
			project, _ = apiClient.ProjectUpdate(mockProject.Id, payload)
		})

		It("Should send PUT request with project ID and expected payload", func() {
			httpCall.Times(1)
		})

		It("Should return project received from API", func() {
			Expect(project).To(Equal(mockedResponse))
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

	Describe("ProjectMove", func() {
		targetProjectId := "targetid"

		payload := struct {
			TargetProjectId string `json:"targetProjectId"`
		}{
			targetProjectId,
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/projects/"+mockProject.Id+"/move", payload, nil).Times(1)
			apiClient.ProjectMove(mockProject.Id, targetProjectId)
		})

		It("Should send POST request with project id and target project id", func() {})
	})
})
