package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("VcsConnection Client", func() {
	mockVcsConnection := VcsConnection{
		Id:          "id0",
		Name:        "test-connection",
		Type:        "GitHubEnterprise",
		Url:         "https://github.example.com",
		VcsAgentKey: "ENV0_DEFAULT",
	}

	Describe("Get VcsConnection", func() {
		var returnedVcsConnection *VcsConnection
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/vcs/connections/"+mockVcsConnection.Id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *VcsConnection) {
					*response = mockVcsConnection
				})
			returnedVcsConnection, err = apiClient.VcsConnection(mockVcsConnection.Id)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return vcs connection", func() {
			Expect(*returnedVcsConnection).To(Equal(mockVcsConnection))
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Create VcsConnection", func() {
		var createdVcsConnection *VcsConnection
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall()

			createPayload := VcsConnectionCreatePayload{}
			_ = copier.Copy(&createPayload, &mockVcsConnection)

			expectedCreateRequest := struct {
				VcsConnectionCreatePayload
				OrganizationId string `json:"organizationId"`
			}{
				VcsConnectionCreatePayload: createPayload,
				OrganizationId:             organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/vcs/connections", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *VcsConnection) {
					*response = mockVcsConnection
				})

			createdVcsConnection, err = apiClient.VcsConnectionCreate(createPayload)
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

		It("Should return created vcs connection", func() {
			Expect(*createdVcsConnection).To(Equal(mockVcsConnection))
		})
	})

	Describe("Update VcsConnection", func() {
		var updatedVcsConnection *VcsConnection
		var err error

		BeforeEach(func() {
			updatePayload := VcsConnectionUpdatePayload{
				Name:        mockVcsConnection.Name,
				VcsAgentKey: mockVcsConnection.VcsAgentKey,
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/vcs/connections/"+mockVcsConnection.Id, updatePayload, gomock.Any()).
				Do(func(path string, request interface{}, response *VcsConnection) {
					*response = mockVcsConnection
				})

			updatedVcsConnection, err = apiClient.VcsConnectionUpdate(mockVcsConnection.Id, updatePayload)
		})

		It("Should send PUT request with params", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return updated vcs connection", func() {
			Expect(*updatedVcsConnection).To(Equal(mockVcsConnection))
		})
	})

	Describe("Delete VcsConnection", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/vcs/connections/"+mockVcsConnection.Id, nil)
			err = apiClient.VcsConnectionDelete(mockVcsConnection.Id)
		})

		It("Should send DELETE request", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("List VcsConnections", func() {
		var returnedVcsConnections []VcsConnection
		var err error
		mockVcsConnections := []VcsConnection{mockVcsConnection}

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/vcs/connections", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]VcsConnection) {
					*response = mockVcsConnections
				})
			returnedVcsConnections, err = apiClient.VcsConnections()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return vcs connections", func() {
			Expect(returnedVcsConnections).To(Equal(mockVcsConnections))
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})
})
