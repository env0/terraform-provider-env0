package client_test

import (
	"fmt"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Approval Policy Client", func() {
	mockApprovalPolicy := ApprovalPolicy{
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

	Describe("Get Custom Flows By Name", func() {
		var returnedApprovalPolicies []ApprovalPolicy
		mockApprovalPolicies := []ApprovalPolicy{mockApprovalPolicy}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/approval-policy", map[string]string{"organizationId": organizationId, "name": mockApprovalPolicy.Name}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApprovalPolicy) {
					*response = mockApprovalPolicies
				})
			organizationIdCall.Times(1)
			httpCall.Times(1)
			returnedApprovalPolicies, _ = apiClient.ApprovalPolicies(mockApprovalPolicy.Name)
		})

		It("Should return approval policies", func() {
			Expect(returnedApprovalPolicies).To(Equal(mockApprovalPolicies))
		})
	})

	mockAssignment := ApprovalPolicyAssignment{
		Scope:       ApprovalPolicyProjectScope,
		ScopeId:     "scope_id",
		BlueprintId: "blueprint_id",
	}

	Describe("Assign Approval Policy", func() {
		var returnedApprovalPolicyAssignment *ApprovalPolicyAssignment

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Post("/approval-policy/assignment", &mockAssignment, gomock.Any()).
				Do(func(path string, request interface{}, response *ApprovalPolicyAssignment) {
					*response = mockAssignment
				})
			httpCall.Times(1)
			returnedApprovalPolicyAssignment, _ = apiClient.ApprovalPolicyAssign(&mockAssignment)
		})

		It("Should return approval policy assignment", func() {
			Expect(*returnedApprovalPolicyAssignment).To(Equal(mockAssignment))
		})
	})

	Describe("Unassign Custom Flow", func() {
		var err error
		id := fmt.Sprintf("%s#%s#%s", mockAssignment.Scope, mockAssignment.ScopeId, mockAssignment.BlueprintId)

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/approval-policy/assignment", map[string]string{"id": id})
			httpCall.Times(1)
			err = apiClient.ApprovalPolicyUnassign(id)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("Get Approval Policy By Scope", func() {
		var ret []ApprovalPolicyByScope

		scope := string(mockAssignment.Scope)
		scopeId := mockAssignment.ScopeId

		mockApprovalPolicyByScope := ApprovalPolicyByScope{
			Scope:          scope,
			ScopeId:        scopeId,
			ApprovalPolicy: &mockApprovalPolicy,
		}

		mockApprovalPolicyByScopeArr := []ApprovalPolicyByScope{mockApprovalPolicyByScope}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get(fmt.Sprintf("/approval-policy/%s/%s", scope, scopeId), nil, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ApprovalPolicyByScope) {
					*response = mockApprovalPolicyByScopeArr
				})
			httpCall.Times(1)
			ret, _ = apiClient.ApprovalPolicyByScope(scope, scopeId)
		})

		It("Should return approval policy assignment", func() {
			Expect(ret).To(Equal(mockApprovalPolicyByScopeArr))
		})
	})

	Describe("Create ApprovalPolicy", func() {
		var result *ApprovalPolicy
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId).Times(1)

			createApprovalPolicyPayload := ApprovalPolicyCreatePayload{}
			copier.Copy(&createApprovalPolicyPayload, &mockApprovalPolicy)

			expectedCreateRequest := struct {
				ApprovalPolicyCreatePayload
				OrganizationId string `json:"organizationId"`
			}{
				createApprovalPolicyPayload,
				organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/approval-policy", &expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *ApprovalPolicy) {
					*response = mockApprovalPolicy
				}).Times(1)

			result, err = apiClient.ApprovalPolicyCreate(&createApprovalPolicyPayload)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return created Approval Policy", func() {
			Expect(*result).To(Equal(mockApprovalPolicy))
		})
	})

	Describe("Update ApprovalPolicy", func() {
		var result *ApprovalPolicy
		var err error

		BeforeEach(func() {
			updateApprovalPolicyPayload := ApprovalPolicyUpdatePayload{}
			copier.Copy(&updateApprovalPolicyPayload, &mockApprovalPolicy)

			httpCall = mockHttpClient.EXPECT().
				Put("/approval-policy", &updateApprovalPolicyPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *ApprovalPolicy) {
					*response = mockApprovalPolicy
				}).Times(1)

			result, err = apiClient.ApprovalPolicyUpdate(&updateApprovalPolicyPayload)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return updated Approval Policy", func() {
			Expect(*result).To(Equal(mockApprovalPolicy))
		})
	})
})
