package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Environment Discovery", func() {
	mockError := errors.New("error")

	projectId := "pid"

	mockEnvironmentDiscovery := EnvironmentDiscoveryPayload{
		Id:                     "id",
		GlobPattern:            "**",
		EnvironmentPlacement:   "topProject",
		WorkspaceNaming:        "default",
		AutoDeployByCustomGlob: "**",
		Repository:             "https://re.po",
		TerraformVersion:       "1.5.6",
		Type:                   "terraform",
		GithubInstallationId:   1000,
	}

	Describe("PUT", func() {
		putPayload := EnvironmentDiscoveryPutPayload{
			GlobPattern:            mockEnvironmentDiscovery.GlobPattern,
			EnvironmentPlacement:   mockEnvironmentDiscovery.EnvironmentPlacement,
			WorkspaceNaming:        mockEnvironmentDiscovery.WorkspaceNaming,
			AutoDeployByCustomGlob: mockEnvironmentDiscovery.AutoDeployByCustomGlob,
			Repository:             mockEnvironmentDiscovery.Repository,
			TerraformVersion:       mockEnvironmentDiscovery.TerraformVersion,
			Type:                   mockEnvironmentDiscovery.Type,
			GithubInstallationId:   mockEnvironmentDiscovery.GithubInstallationId,
		}

		Describe("success", func() {
			var ret *EnvironmentDiscoveryPayload

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/environment-discovery/projects/"+projectId, &putPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentDiscoveryPayload) {
						*response = mockEnvironmentDiscovery
					}).Times(1)
				ret, _ = apiClient.PutEnvironmentDiscovery(projectId, &putPayload)
			})

			It("Should return environment discovery", func() {
				Expect(*ret).To(Equal(mockEnvironmentDiscovery))
			})
		})

		Describe("failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/environment-discovery/projects/"+projectId, &putPayload, gomock.Any()).Return(mockError).Times(1)
				_, err = apiClient.PutEnvironmentDiscovery(projectId, &putPayload)
			})

			It("Should return error", func() {
				Expect(err).To(Equal(mockError))
			})
		})
	})

	Describe("GET", func() {
		Describe("success", func() {
			var ret *EnvironmentDiscoveryPayload

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/environment-discovery/projects/"+projectId, nil, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentDiscoveryPayload) {
						*response = mockEnvironmentDiscovery
					}).Times(1)
				ret, _ = apiClient.GetEnvironmentDiscovery(projectId)
			})

			It("Should return environment discovery", func() {
				Expect(*ret).To(Equal(mockEnvironmentDiscovery))
			})
		})

		Describe("failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/environment-discovery/projects/"+projectId, nil, gomock.Any()).Return(mockError).Times(1)
				_, err = apiClient.GetEnvironmentDiscovery(projectId)
			})

			It("Should return error", func() {
				Expect(err).To(Equal(mockError))
			})
		})
	})

	Describe("DELETE", func() {
		Describe("success", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().Delete("/environment-discovery/projects/"+projectId, nil).Times(1)
				err = apiClient.DeleteEnvironmentDiscovery(projectId)
			})

			It("Should send DELETE request", func() {})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Delete("/environment-discovery/projects/"+projectId, nil).Return(mockError).Times(1)
				err = apiClient.DeleteEnvironmentDiscovery(projectId)
			})

			It("Should return error", func() {
				Expect(err).To(Equal(mockError))
			})
		})
	})
})
