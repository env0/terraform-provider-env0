package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Custom Flow Client", func() {
	mockCustomFlow := CustomFlow{
		Id:         "id",
		Name:       "name",
		Repository: "repository",
		Path:       "path",
		Revision:   "revision",
		TokenId:    "tokenId",
		SshKeys: []TemplateSshKey{
			{Id: "id", Name: "name"},
		},
		GithubInstallationId: 1,
		BitbucketClientKey:   "bitbucket-key",
		IsBitbucketServer:    true,
		IsGitlabEnterprise:   false,
		IsGithubEnterprise:   true,
		IsGitLab:             false,
		IsAzureDevOps:        true,
		IsTerragruntRunAll:   false,
	}

	Describe("Get Custom Flow", func() {
		var returnedCustomFlow *CustomFlow

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/custom-flow/"+mockCustomFlow.Id, gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *CustomFlow) {
					*response = mockCustomFlow
				})
			returnedCustomFlow, _ = apiClient.CustomFlow(mockCustomFlow.Id)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return custom flow", func() {
			Expect(*returnedCustomFlow).To(Equal(mockCustomFlow))
		})
	})

	Describe("Get Custom Flows By Name", func() {
		var returnedCustomFlows []CustomFlow
		mockCustomFlows := []CustomFlow{mockCustomFlow}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/custom-flows", map[string]string{"organizationId": organizationId, "name": mockCustomFlow.Name}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]CustomFlow) {
					*response = mockCustomFlows
				})
			returnedCustomFlows, _ = apiClient.CustomFlows(mockCustomFlow.Name)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return custom flows", func() {
			Expect(returnedCustomFlows).To(Equal(mockCustomFlows))
		})
	})

	Describe("Create Custom Flow", func() {
		var createdCustomFlow *CustomFlow
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			var createCustomFlowPayload CustomFlowCreatePayload
			copier.Copy(&createCustomFlowPayload, &mockCustomFlow)

			expectedCreateRequest := struct {
				OrganizationId string `json:"organizationId"`
				CustomFlowCreatePayload
			}{
				organizationId,
				createCustomFlowPayload,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/custom-flow", &expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *CustomFlow) {
					*response = mockCustomFlow
				})

			createdCustomFlow, err = apiClient.CustomFlowCreate(createCustomFlowPayload)
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

		It("Should return created custom flow", func() {
			Expect(*createdCustomFlow).To(Equal(mockCustomFlow))
		})
	})

	Describe("Delete Custom Flow", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/custom-flow/" + mockCustomFlow.Id)
			apiClient.CustomFlowDelete(mockCustomFlow.Id)
		})

		It("Should send DELETE request with custom flow id", func() {
			httpCall.Times(1)
		})
	})

	Describe("Update Custom Flow", func() {
		var updatedCustomFlow *CustomFlow
		var err error

		updatedMockCustomFlow := mockCustomFlow
		updatedMockCustomFlow.Path = "updated-path"

		var updateCustomFlowPayload CustomFlowCreatePayload
		copier.Copy(&updateCustomFlowPayload, &updatedMockCustomFlow)

		BeforeEach(func() {
			expectedUpdateRequest := struct {
				Id string `json:"id"`
				CustomFlowCreatePayload
			}{
				updatedMockCustomFlow.Id,
				updateCustomFlowPayload,
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/custom-flow", &expectedUpdateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *CustomFlow) {
					*response = updatedMockCustomFlow
				})

			updatedCustomFlow, err = apiClient.CustomFlowUpdate(updatedMockCustomFlow.Id, updateCustomFlowPayload)
		})

		It("Should send PUT request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return updated custom flow", func() {
			Expect(*updatedCustomFlow).To(Equal(updatedMockCustomFlow))
		})
	})
})
