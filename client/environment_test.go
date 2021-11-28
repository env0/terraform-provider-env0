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

			It("Should return teams", func() {
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

		Describe("Success", func() {
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
	})
})
