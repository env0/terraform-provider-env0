package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

const projectName = "project_test"
const projectDescription = "project description"
const parentProjectId = "parent_project_id"

var _ = Describe("Project", func() {
	var project Project
	var moduleTestingProject *ModuleTestingProject
	mockProject := Project{
		Id:             "idX",
		Name:           "projectX",
		Description:    "descriptionX",
		OrganizationId: organizationId,
	}

	mockModuleTestingProject := ModuleTestingProject{
		Id:   "idx",
		Name: "namex",
	}

	Describe("ProjectCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall()

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
				Do(func(path string, request any, response *Project) {
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
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/projects/"+mockProject.Id, nil)
			err = apiClient.ProjectDelete(mockProject.Id)
		})

		It("Should send DELETE request with project id", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
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
				Do(func(path string, request any, response *Project) {
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
				Do(func(path string, request any, response *Project) {
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
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/projects", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]Project) {
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
		var err error

		targetProjectId := "targetid"

		payload := struct {
			TargetProjectId *string `json:"targetProjectId"`
		}{
			&targetProjectId,
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/projects/"+mockProject.Id+"/move", payload, nil).Times(1)
			err = apiClient.ProjectMove(mockProject.Id, targetProjectId)
		})

		It("Should send POST request with project id and target project id", func() {})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("ProjectMove with no target project id", func() {
		var err error

		payload := struct {
			TargetProjectId *string `json:"targetProjectId"`
		}{
			nil,
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/projects/"+mockProject.Id+"/move", payload, nil).Times(1)
			err = apiClient.ProjectMove(mockProject.Id, "")
		})

		It("Should send POST request with project id and nil target project id", func() {})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("ModuleTestingProject", func() {
		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/projects/modules/testing/"+organizationId, nil, gomock.Any()).
				Do(func(path string, request any, response *ModuleTestingProject) {
					*response = mockModuleTestingProject
				})
			moduleTestingProject, _ = apiClient.ModuleTestingProject()
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return module testing project", func() {
			Expect(mockModuleTestingProject).To(Equal(*moduleTestingProject))
		})
	})
})
