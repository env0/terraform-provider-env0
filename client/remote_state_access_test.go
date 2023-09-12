package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("RemoteStateAccess", func() {
	environmentId := "environmnet_id"

	remoteStateAccess := RemoteStateAccessConfiguration{
		EnvironmentId:                    environmentId,
		AccessibleFromEntireOrganization: false,
		AllowedProjectIds: []string{
			"pid1",
		},
	}

	Describe("Create", func() {
		var err error
		var remoteStateAccessResponse *RemoteStateAccessConfiguration

		BeforeEach(func() {
			createRequest := RemoteStateAccessConfigurationCreate{
				AllowedProjectIds: remoteStateAccess.AllowedProjectIds,
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/remote-backend/states/"+environmentId+"/access-control", createRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *RemoteStateAccessConfiguration) {
					*response = remoteStateAccess
				})
			httpCall.Times(1)
			remoteStateAccessResponse, err = apiClient.RemoteStateAccessConfigurationCreate(environmentId, createRequest)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return created remote state access configuration", func() {
			Expect(*remoteStateAccessResponse).To(Equal(remoteStateAccess))
		})
	})

	Describe("Get", func() {
		var err error
		var remoteStateAccessResponse *RemoteStateAccessConfiguration

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/remote-backend/states/"+environmentId+"/access-control", gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *RemoteStateAccessConfiguration) {
					*response = remoteStateAccess
				})
			httpCall.Times(1)
			remoteStateAccessResponse, err = apiClient.RemoteStateAccessConfiguration(environmentId)
		})

		It("Should return remote state access configuration", func() {
			Expect(*remoteStateAccessResponse).To(Equal(remoteStateAccess))
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Delete", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/remote-backend/states/"+environmentId+"/access-control", nil)
			httpCall.Times(1)
			err = apiClient.RemoteStateAccessConfigurationDelete(environmentId)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})
})
