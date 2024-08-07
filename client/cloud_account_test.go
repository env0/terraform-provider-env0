package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CloudAccount", func() {
	var account *CloudAccount
	var accounts []CloudAccount
	var err error

	awsConfiguration := AWSCloudAccountConfiguration{
		AccountId:  "a",
		BucketName: "b",
		Prefix:     "prefix",
		Regions:    []string{"us-east-10", "us-west-24"},
	}

	awsConfigurationUpdated := AWSCloudAccountConfiguration{
		AccountId:  "c",
		BucketName: "d",
		Prefix:     "prefix2",
		Regions:    []string{"us-east-10"},
	}

	account1 := CloudAccount{
		Id:            "id1",
		Provider:      "AWS",
		Name:          "name1",
		Health:        true,
		Configuration: &awsConfiguration,
	}

	account1Updated := account1
	account1Updated.Name = "updatedname1"
	account1Updated.Configuration = awsConfigurationUpdated

	account2 := CloudAccount{
		Id:            "id2",
		Provider:      "GCP",
		Name:          "name2",
		Configuration: []string{"some random configuration"},
	}

	Describe("create", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			payload := CloudAccountCreatePayload{
				Provider:      account1.Provider,
				Name:          account1.Name,
				Configuration: account1.Configuration,
			}

			payloadWithOrganizationId := struct {
				*CloudAccountCreatePayload
				OrganizationId string `json:"organizationId"`
			}{
				&payload,
				organizationId,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/cloud/configurations", &payloadWithOrganizationId, gomock.Any()).
				Do(func(path string, request interface{}, response *CloudAccount) {
					*response = account1
				}).Times(1)

			account, err = apiClient.CloudAccountCreate(&payload)
		})

		It("should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("should return account", func() {
			Expect(*account).To(Equal(account1))
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("update", func() {
		BeforeEach(func() {
			payload := CloudAccountUpdatePayload{
				Name:          account1Updated.Name,
				Configuration: account1Updated.Configuration,
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/cloud/configurations/"+account.Id, &payload, gomock.Any()).
				Do(func(path string, request interface{}, response *CloudAccount) {
					*response = account1Updated
				}).Times(1)

			account, err = apiClient.CloudAccountUpdate(account.Id, &payload)
		})

		It("should return updated account", func() {
			Expect(*account).To(Equal(account1Updated))
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("delete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Delete("/cloud/configurations/"+account.Id, nil).
				Do(func(path string, request interface{}) {}).Times(1)

			err = apiClient.CloudAccountDelete(account.Id)
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("get", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/cloud/configurations/"+account.Id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *CloudAccount) {
					*response = account1
				}).Times(1)

			account, err = apiClient.CloudAccount(account.Id)
		})

		It("should return account", func() {
			Expect(*account).To(Equal(account1))
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("list", func() {
		mockedAccounts := []CloudAccount{
			account1,
			account2,
		}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/cloud/configurations", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]CloudAccount) {
					*response = mockedAccounts
				}).Times(1)

			accounts, err = apiClient.CloudAccounts()
		})

		It("should return accounts", func() {
			Expect(accounts).To(Equal(mockedAccounts))
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})
})
