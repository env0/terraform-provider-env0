package client_test

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const full_page = 100
const partial_page = 33

var _ = Describe("Environment Client", func() {
	const (
		environmentId = "env-id"
	)

	mockEnvironment := Environment{
		Id:   environmentId,
		Name: "env0",
	}

	mockTemplate := Template{
		Id:         "template-id",
		Name:       "template-name",
		Repository: "https://re.po",
	}

	Describe("Environments", func() {
		var environments []Environment
		mockEnvironments := []Environment{mockEnvironment}
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)
				httpCall = mockHttpClient.EXPECT().
					Get("/environments", map[string]string{
						"limit":          "100",
						"offset":         "0",
						"organizationId": organizationId,
						"name":           mockEnvironment.Name,
					}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Environment) {
						*response = mockEnvironments
					})

				environments, err = apiClient.EnvironmentsByName(mockEnvironment.Name)
			})

			It("Should send GET request", func() {
				httpCall.Times(1)
			})

			It("Should return the environment", func() {
				Expect(environments).To(Equal(mockEnvironments))
			})
		})

		Describe("SuccessMultiPages", func() {
			var environmentsP1, environmentsP2 []Environment
			for i := 0; i < full_page; i++ {
				environmentsP1 = append(environmentsP1, mockEnvironment)
			}

			for i := 0; i < partial_page; i++ {
				environmentsP2 = append(environmentsP2, mockEnvironment)
			}

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)
				httpCall = mockHttpClient.EXPECT().
					Get("/environments", map[string]string{
						"offset":         "0",
						"limit":          "100",
						"organizationId": organizationId,
						"name":           mockEnvironment.Name,
					}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Environment) {
						*response = environmentsP1
					}).Times(1)

				httpCall2 = mockHttpClient.EXPECT().
					Get("/environments", map[string]string{
						"offset":         "100",
						"limit":          "100",
						"organizationId": organizationId,
						"name":           mockEnvironment.Name,
					}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Environment) {
						*response = environmentsP2
					}).Times(1)

				environments, err = apiClient.EnvironmentsByName(mockEnvironment.Name)
			})

			It("Should return the environments", func() {
				Expect(environments).To(Equal(append(environmentsP1, environmentsP2...)))
			})
		})

		Describe("SuccessMultiPagesWithProject", func() {
			projectId := "proj123"
			var environmentsP1, environmentsP2 []Environment
			for i := 0; i < full_page; i++ {
				environmentsP1 = append(environmentsP1, mockEnvironment)
			}

			for i := 0; i < partial_page; i++ {
				environmentsP2 = append(environmentsP2, mockEnvironment)
			}

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/environments", map[string]string{
						"offset":    "0",
						"limit":     "100",
						"projectId": projectId,
					}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Environment) {
						*response = environmentsP1
					}).Times(1)

				httpCall2 = mockHttpClient.EXPECT().
					Get("/environments", map[string]string{
						"offset":    "100",
						"limit":     "100",
						"projectId": projectId,
					}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Environment) {
						*response = environmentsP2
					}).Times(1)

				environments, err = apiClient.ProjectEnvironments(projectId)
			})

			It("Should return the environments", func() {
				Expect(environments).To(Equal(append(environmentsP1, environmentsP2...)))
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				mockOrganizationIdCall(organizationId)
				httpCall = mockHttpClient.EXPECT().
					Get("/environments", map[string]string{
						"limit":          "100",
						"offset":         "0",
						"organizationId": organizationId,
						"name":           mockEnvironment.Name,
					}, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.EnvironmentsByName(mockEnvironment.Name)
				Expect(expectedErr).Should(Equal(err))
			})
		})
	})

	Describe("Environment", func() {
		var environment Environment
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/environments/"+mockEnvironment.Id, nil, gomock.Any()).
					Do(func(path string, request interface{}, response *Environment) {
						*response = mockEnvironment
					})

				environment, err = apiClient.Environment(mockEnvironment.Id)
			})

			It("Should send GET request", func() {
				httpCall.Times(1)
			})

			It("Should return environments", func() {
				Expect(environment).To(Equal(mockEnvironment))
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get("/environments/"+mockEnvironment.Id, nil, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.Environment(mockEnvironment.Id)
				Expect(expectedErr).Should(Equal(err))
			})
		})
	})

	Describe("EnvironmentCreate", func() {
		var createdEnvironment Environment
		var err error

		BeforeEach(func() {
			createEnvironmentPayload := EnvironmentCreate{}
			copier.Copy(&createEnvironmentPayload, &mockEnvironment)

			expectedCreateRequest := createEnvironmentPayload

			httpCall = mockHttpClient.EXPECT().
				Post("/environments", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *Environment) {
					*response = mockEnvironment
				})

			createdEnvironment, err = apiClient.EnvironmentCreate(createEnvironmentPayload)
		})

		It("Should send POST request", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return the created environment", func() {
			Expect(createdEnvironment).To(Equal(mockEnvironment))
		})
	})

	Describe("EnvironmentCreateWithoutTemplate", func() {
		var createdEnvironment Environment
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			createEnvironmentPayload := EnvironmentCreate{}
			copier.Copy(&createEnvironmentPayload, &mockEnvironment)
			createTemplatePayload := TemplateCreatePayload{}
			copier.Copy(&createTemplatePayload, &mockTemplate)

			createRequest := EnvironmentCreateWithoutTemplate{
				EnvironmentCreate: createEnvironmentPayload,
				TemplateCreate:    createTemplatePayload,
			}

			expectedCreateRequest := createRequest
			expectedCreateRequest.TemplateCreate.OrganizationId = organizationId

			httpCall = mockHttpClient.EXPECT().
				Post("/environments/without-template", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *Environment) {
					*response = mockEnvironment
				})

			createdEnvironment, err = apiClient.EnvironmentCreateWithoutTemplate(createRequest)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return the created environment", func() {
			Expect(createdEnvironment).To(Equal(mockEnvironment))
		})
	})

	Describe("EnvironmentDelete", func() {
		var err error
		var res *EnvironmentDestroyResponse

		mockedRes := EnvironmentDestroyResponse{
			Id: "id123",
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/environments/"+mockEnvironment.Id+"/destroy", nil, gomock.Any()).Times(1).
				Do((func(path string, request interface{}, response *EnvironmentDestroyResponse) {
					*response = mockedRes
				}))
			res, err = apiClient.EnvironmentDestroy(mockEnvironment.Id)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return the expected response", func() {
			Expect(*res).To(Equal(mockedRes))
		})
	})

	Describe("EnvironmentMarkAsArchived", func() {
		payload := struct {
			IsArchived bool `json:"isArchived"`
		}{true}

		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Put("/environments/"+mockEnvironment.Id, &payload, nil).Return(nil).Times(1)
			err = apiClient.EnvironmentMarkAsArchived(mockEnvironment.Id)
		})

		It("Should send an archive put request", func() {})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

	})

	Describe("EnvironmentUpdate", func() {
		Describe("Success", func() {
			var updatedEnvironment Environment
			var err error

			BeforeEach(func() {
				updateEnvironmentPayload := EnvironmentUpdate{Name: "updated-name"}

				httpCall = mockHttpClient.EXPECT().
					Put("/environments/"+mockEnvironment.Id, updateEnvironmentPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *Environment) {
						*response = mockEnvironment
					})

				updatedEnvironment, err = apiClient.EnvironmentUpdate(mockEnvironment.Id, updateEnvironmentPayload)
			})

			It("Should send Put request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return the environment received from API", func() {
				Expect(updatedEnvironment).To(Equal(mockEnvironment))
			})
		})
	})

	Describe("EnvironmentDeploy", func() {
		Describe("Success", func() {
			var response EnvironmentDeployResponse
			var err error
			deployResponseMock := EnvironmentDeployResponse{
				Id: "deployment-id",
			}

			BeforeEach(func() {
				userRequiresApproval := false
				deployRequest := DeployRequest{
					BlueprintId:          "",
					BlueprintRevision:    "",
					BlueprintRepository:  "",
					ConfigurationChanges: nil,
					TTL:                  nil,
					EnvName:              "",
					UserRequiresApproval: &userRequiresApproval,
				}

				httpCall = mockHttpClient.EXPECT().
					Post("/environments/"+mockEnvironment.Id+"/deployments", deployRequest, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentDeployResponse) {
						*response = deployResponseMock
					})

				response, err = apiClient.EnvironmentDeploy(mockEnvironment.Id, deployRequest)
			})

			It("Should send post request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return the deployment id received from API", func() {
				Expect(response).To(Equal(deployResponseMock))
			})
		})
	})

	Describe("EnvironmentUpdateTTL", func() {
		Describe("Success", func() {
			var updatedEnvironment Environment
			var err error

			BeforeEach(func() {
				updateTTLRequest := TTL{
					Type:  "",
					Value: "",
				}

				httpCall = mockHttpClient.EXPECT().
					Put("/environments/"+mockEnvironment.Id+"/ttl", updateTTLRequest, gomock.Any()).
					Do(func(path string, request interface{}, response *Environment) {
						*response = mockEnvironment
					})

				updatedEnvironment, err = apiClient.EnvironmentUpdateTTL(mockEnvironment.Id, updateTTLRequest)
			})

			It("Should send Put request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return the deployment id received from API", func() {
				Expect(updatedEnvironment).To(Equal(mockEnvironment))
			})
		})
	})

	Describe("Environment Move", func() {
		var err error

		environmentId := "envid"
		projectId := "projid"

		request := EnvironmentMoveRequest{
			ProjectId: projectId,
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/environments/"+environmentId+"/move", &request, nil).Times(1)
			err = apiClient.EnvironmentMove(environmentId, projectId)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("EnvironmentDeployment", func() {
		var deployment *DeploymentLog
		var err error

		mockDeployment := DeploymentLog{
			Id:     "id12345",
			Status: "IN_PROGRESS",
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/environments/deployments/"+mockDeployment.Id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *DeploymentLog) {
					*response = mockDeployment
				}).Times(1)

			deployment, err = apiClient.EnvironmentDeployment(mockDeployment.Id)
		})

		It("Should return deployment", func() {
			Expect(*deployment).To(Equal(mockDeployment))
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})
})

func TestMarshalEnvironmentCreateWithoutTemplate(t *testing.T) {
	templateCreate := TemplateCreatePayload{
		Name: "name",
		SshKeys: []TemplateSshKey{
			{Id: "id1", Name: "name1"},
		},
		Type: "terraform",
	}
	environmentCreate := EnvironmentCreate{
		Name:      "name",
		ProjectId: "project_id",
	}

	environmentCreateWithoutTemplate := EnvironmentCreateWithoutTemplate{
		EnvironmentCreate: environmentCreate,
		TemplateCreate:    templateCreate,
	}

	b, err := json.Marshal(&environmentCreateWithoutTemplate)
	require.NoError(t, err)

	var templateCreateFromJson TemplateCreatePayload

	require.NoError(t, json.Unmarshal(b, &templateCreateFromJson))
	require.Equal(t, templateCreate, templateCreateFromJson)

	var environmentCreateFromJSON EnvironmentCreate

	require.NoError(t, json.Unmarshal(b, &environmentCreateFromJSON))

	environmentCreateWithType := environmentCreate
	environmentCreateWithType.Type = templateCreate.Type
	require.Equal(t, environmentCreateWithType, environmentCreateFromJSON)
}
