package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GitToken Client", func() {
	mockGitToken := GitToken{
		Id:             "id",
		Name:           "name",
		Value:          "value",
		OrganizationId: organizationId,
	}

	Describe("Get Single GitToken", func() {
		var returnedGitToken *GitToken

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/tokens/"+mockGitToken.Id, gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *GitToken) {
					*response = mockGitToken
				})
			returnedGitToken, _ = apiClient.GitToken(mockGitToken.Id)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return GitToken", func() {
			Expect(*returnedGitToken).To(Equal(mockGitToken))
		})
	})

	Describe("Get All GitTokens", func() {
		var returnedGitTokens []GitToken
		mockGitTokens := []GitToken{mockGitToken}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/tokens", map[string]string{"organizationId": organizationId, "type": "GIT"}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]GitToken) {
					*response = mockGitTokens
				})
			returnedGitTokens, _ = apiClient.GitTokens()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return GitTokens", func() {
			Expect(returnedGitTokens).To(Equal(mockGitTokens))
		})
	})

	Describe("Create GitToken", func() {
		var createdGitToken *GitToken
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			createGitTokenPayload := GitTokenCreatePayload{}
			copier.Copy(&createGitTokenPayload, &mockGitToken)

			expectedCreateRequest := GitTokenCreatePayloadWith{
				GitTokenCreatePayload: createGitTokenPayload,
				OrganizationId:        organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/tokens", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *GitToken) {
					*response = mockGitToken
				})

			createdGitToken, err = apiClient.GitTokenCreate(createGitTokenPayload)
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

		It("Should return created GitToken", func() {
			Expect(*createdGitToken).To(Equal(mockGitToken))
		})
	})

	Describe("Delete GitToken", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/tokens/" + mockGitToken.Id)
			apiClient.GitTokenDelete(mockGitToken.Id)
		})

		It("Should send DELETE request with GitToken id", func() {
			httpCall.Times(1)
		})
	})
})
