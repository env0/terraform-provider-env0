package client_test

import (
	"encoding/json"
	"errors"
	"strconv"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Templates Client", func() {
	mockTemplate := Template{
		Id:         "template-id",
		Name:       "template-name",
		Repository: "https://re.po",
	}

	Describe("TemplateCreatePayload", func() {
		DescribeTable("Github Installation Id",
			func(value int, expected types.GomegaMatcher) {
				payload := TemplateCreatePayload{
					GithubInstallationId: value,
				}
				jsonPayload, _ := json.Marshal(payload)
				var parsedPayload map[string]interface{}
				json.Unmarshal(jsonPayload, &parsedPayload)
				Expect(parsedPayload["githubInstallationId"]).To(expected)
			},
			Entry("Has value", 123, BeEquivalentTo(123)),
			Entry("No value", nil, BeNil()),
		)
	})

	Describe("Template", func() {
		var returnedTemplate Template

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/blueprints/"+mockTemplate.Id, gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *Template) {
					*response = mockTemplate
				})
			returnedTemplate, _ = apiClient.Template(mockTemplate.Id)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return template", func() {
			Expect(returnedTemplate).To(Equal(mockTemplate))
		})
	})

	Describe("Templates", func() {
		var returnedTemplates []Template
		mockTemplates := []Template{mockTemplate}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			expectedPayload := map[string]string{"organizationId": organizationId}
			httpCall = mockHttpClient.EXPECT().
				Get("/blueprints", expectedPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *[]Template) {
					*response = mockTemplates
				})
			returnedTemplates, _ = apiClient.Templates()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return template", func() {
			Expect(returnedTemplates).To(Equal(mockTemplates))
		})
	})

	Describe("TemplateCreate", func() {
		var createdTemplate Template
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			createTemplatePayload := TemplateCreatePayload{}
			copier.Copy(&createTemplatePayload, &mockTemplate)

			expectedCreateRequest := createTemplatePayload
			expectedCreateRequest.OrganizationId = organizationId

			httpCall = mockHttpClient.EXPECT().
				Post("/blueprints", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *Template) {
					*response = mockTemplate
				})

			createdTemplate, err = apiClient.TemplateCreate(createTemplatePayload)
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

		It("Should return created configuration variable", func() {
			Expect(createdTemplate).To(Equal(mockTemplate))
		})
	})

	Describe("TemplateDelete", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/blueprints/"+mockTemplate.Id, nil).Times(1)
			err = apiClient.TemplateDelete(mockTemplate.Id)
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("TemplateUpdate", func() {
		var updatedTemplate Template
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			updateTemplatePayload := TemplateCreatePayload{}
			copier.Copy(&updateTemplatePayload, &mockTemplate)

			expectedUpdateRequest := updateTemplatePayload
			expectedUpdateRequest.OrganizationId = organizationId

			httpCall = mockHttpClient.EXPECT().
				Put("/blueprints/"+mockTemplate.Id, expectedUpdateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *Template) {
					*response = mockTemplate
				})

			updatedTemplate, err = apiClient.TemplateUpdate(mockTemplate.Id, updateTemplatePayload)
		})

		It("Should send POST request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return configuration value received from API", func() {
			Expect(updatedTemplate).To(Equal(mockTemplate))
		})
	})

	Describe("add template to project", func() {
		templateAssignmentToProjectPayload := TemplateAssignmentToProjectPayload{
			ProjectId: "project-id",
		}
		var assignedTemplate Template
		var err error

		BeforeEach(func() {

			expectedAssignRequest := templateAssignmentToProjectPayload

			httpCall = mockHttpClient.EXPECT().
				Patch("/blueprints/"+mockTemplate.Id+"/projects", expectedAssignRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *Template) {
					*response = mockTemplate
				})

			assignedTemplate, err = apiClient.AssignTemplateToProject(mockTemplate.Id, templateAssignmentToProjectPayload)
		})

		It("Should send POST request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return template from API", func() {
			Expect(assignedTemplate).To(Equal(mockTemplate))
		})
		It("should return an error on empty projectId", func() {
			assignedTemplate, err = apiClient.AssignTemplateToProject(mockTemplate.Id, TemplateAssignmentToProjectPayload{})
			Expect(err).To(BeEquivalentTo(errors.New("must specify projectId on assignment to a template")))
		})
	})

	Describe("remove template from project", func() {
		var err error

		projectId := "project-id"
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/blueprints/"+mockTemplate.Id+"/projects/"+projectId, nil).Times(1)
			err = apiClient.RemoveTemplateFromProject(mockTemplate.Id, projectId)
		})

		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("VariablesFromRepository", func() {
		var returnedVariables []ConfigurationVariable
		var err error

		payload := &VariablesFromRepositoryPayload{
			GithubInstallationId: 1111,
			Path:                 "path",
			Revision:             "1",
			Repository:           "main",
			SshKeyIds: []string{
				"1", "2",
			},
		}

		expectedParams := map[string]string{
			"githubInstallationId": strconv.Itoa(payload.GithubInstallationId),
			"path":                 payload.Path,
			"revision":             payload.Revision,
			"repository":           payload.Repository,
			"sshKeyIds":            `["1","2"]`,
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/blueprints/variables-from-repository", expectedParams, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationVariable) {
					*response = []ConfigurationVariable{}
				})
			returnedVariables, err = apiClient.VariablesFromRepository(payload)
		})

		It("Should return variables", func() {
			Expect(err).To(BeNil())
			Expect(returnedVariables).To(Equal([]ConfigurationVariable{}))
		})
	})
})
