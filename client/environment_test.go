package client_test

import (
	"errors"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	environmentId = "env-id"
)

var _ = Describe("Environment Client", func() {
	mockEnvironment := Environment{
		Id:   environmentId,
		Name: "env0",
	}

	Describe("Environments", func() {
		var environments []Environment
		mockEnvironments := []Environment{mockEnvironment}
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/environments", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Environment) {
						*response = mockEnvironments
					})

				environments, err = apiClient.Environments()
			})

			It("Should send GET request", func() {
				httpCall.Times(1)
			})

			It("Should return the environment", func() {
				Expect(environments).To(Equal(mockEnvironments))
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get("/environments", nil, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.Environments()
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

	Describe("EnvironmentDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/environments/"+mockEnvironment.Id+"/destroy", nil, gomock.Any())
			apiClient.EnvironmentDestroy(mockEnvironment.Id)
		})

		It("Should send a destroy request", func() {
			httpCall.Times(1)
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
				deployRequest := DeployRequest{
					BlueprintId:          "",
					BlueprintRevision:    "",
					BlueprintRepository:  "",
					ConfigurationChanges: nil,
					TTL:                  nil,
					EnvName:              "",
					UserRequiresApproval: false,
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
				updateTTLRequest := EnvironmentUpdateTTL{
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
})
