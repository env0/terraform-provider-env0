package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
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
				Do(func(path string, request any, response *CustomFlow) {
					*response = mockCustomFlow
				})
			httpCall.Times(1)

			returnedCustomFlow, _ = apiClient.CustomFlow(mockCustomFlow.Id)
		})

		It("Should return custom flow", func() {
			Expect(*returnedCustomFlow).To(Equal(mockCustomFlow))
		})
	})

	Describe("Get Custom Flows By Name", func() {
		var returnedCustomFlows []CustomFlow

		mockCustomFlows := []CustomFlow{mockCustomFlow}

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/custom-flows", map[string]string{"organizationId": organizationId, "name": mockCustomFlow.Name}, gomock.Any()).
				Do(func(path string, request any, response *[]CustomFlow) {
					*response = mockCustomFlows
				})

			organizationIdCall.Times(1)
			httpCall.Times(1)

			returnedCustomFlows, _ = apiClient.CustomFlows(mockCustomFlow.Name)
		})

		It("Should return custom flows", func() {
			Expect(returnedCustomFlows).To(Equal(mockCustomFlows))
		})
	})

	Describe("Create Custom Flow", func() {
		var (
			createdCustomFlow *CustomFlow
			err               error
		)

		BeforeEach(func() {
			mockOrganizationIdCall()

			var createCustomFlowPayload CustomFlowCreatePayload

			_ = copier.Copy(&createCustomFlowPayload, &mockCustomFlow)

			expectedCreateRequest := struct {
				OrganizationId string `json:"organizationId"`
				CustomFlowCreatePayload
			}{
				organizationId,
				createCustomFlowPayload,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/custom-flow", &expectedCreateRequest, gomock.Any()).
				Do(func(path string, request any, response *CustomFlow) {
					*response = mockCustomFlow
				})
			httpCall.Times(1)
			organizationIdCall.Times(1)

			createdCustomFlow, err = apiClient.CustomFlowCreate(createCustomFlowPayload)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return created custom flow", func() {
			Expect(*createdCustomFlow).To(Equal(mockCustomFlow))
		})
	})

	Describe("Delete Custom Flow", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/custom-flow/"+mockCustomFlow.Id, nil)
			httpCall.Times(1)

			err = apiClient.CustomFlowDelete(mockCustomFlow.Id)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Update Custom Flow", func() {
		var (
			updatedCustomFlow *CustomFlow
			err               error
		)

		updatedMockCustomFlow := mockCustomFlow
		updatedMockCustomFlow.Path = "updated-path"

		var updateCustomFlowPayload CustomFlowCreatePayload

		_ = copier.Copy(&updateCustomFlowPayload, &updatedMockCustomFlow)

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
				Do(func(path string, request any, response *CustomFlow) {
					*response = updatedMockCustomFlow
				})
			httpCall.Times(1)

			updatedCustomFlow, err = apiClient.CustomFlowUpdate(updatedMockCustomFlow.Id, updateCustomFlowPayload)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return updated custom flow", func() {
			Expect(*updatedCustomFlow).To(Equal(updatedMockCustomFlow))
		})
	})

	mockAssignment := CustomFlowAssignment{
		Scope:       CustomFlowProjectScope,
		ScopeId:     "scope_id",
		BlueprintId: "blueprint_id",
	}

	mockAssignmentList := []CustomFlowAssignment{mockAssignment}

	Describe("Assign Custom Flow", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/custom-flow/assign", mockAssignmentList, nil)
			httpCall.Times(1)

			err = apiClient.CustomFlowAssign(mockAssignmentList)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Unassign Custom Flow", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/custom-flow/unassign", mockAssignmentList, nil)
			httpCall.Times(1)

			err = apiClient.CustomFlowUnassign(mockAssignmentList)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Custom Flow Assignments", func() {
		var (
			err         error
			assignments []CustomFlowAssignment
		)

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Post("/custom-flow/get-assignments", mockAssignmentList, gomock.Any()).
				Do(func(path string, request any, response *[]CustomFlowAssignment) {
					*response = mockAssignmentList
				})
			httpCall.Times(1)

			assignments, err = apiClient.CustomFlowGetAssignments(mockAssignmentList)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return custom flow assignments", func() {
			Expect(assignments).To(Equal(mockAssignmentList))
		})
	})
})
