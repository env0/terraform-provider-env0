package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CloudCredentials", func() {
	const credentialsName = "credential_test"
	var credentials Credentials

	mockCredentialsForGoogleCostCred := Credentials{
		Id:             "id1",
		Name:           "key1",
		OrganizationId: organizationId,
		Type:           "GCP_CREDENTIALS",
	}

	mockCredentials := Credentials{
		Id:             "id1",
		Name:           "key1",
		OrganizationId: organizationId,
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	mockCredentialsSecond := Credentials{
		Id:             "id2",
		Name:           "key2",
		OrganizationId: organizationId,
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	keys := []Credentials{mockCredentials, mockCredentialsSecond}

	Describe("GoogleCostCredentialsCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall()

			payloadValue := GoogleCostCredentialsValuePayload{
				TableId: "table",
				Secret:  "secret",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &GoogleCostCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           "GCP_CREDENTIALS",
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentialsForGoogleCostCred
				})

			credentials, _ = apiClient.CredentialsCreate(&GoogleCostCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  "GCP_CREDENTIALS",
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockCredentialsForGoogleCostCred))
		})
	})

	Describe("AwsCredentialsCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall()

			payloadValue := AwsCredentialsValuePayload{
				RoleArn:  "role",
				Duration: 1,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &AwsCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentials
				})

			credentials, _ = apiClient.CredentialsCreate(&AwsCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockCredentials))
		})
	})

	Describe("AwsCredentialsCreate with ProjectId", func() {
		const projectId = "project-123"

		BeforeEach(func() {
			// Note: No organization ID call should be made

			payloadValue := AwsCredentialsValuePayload{
				RoleArn:  "role",
				Duration: 1,
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &AwsCredentialsCreatePayload{
					Name:      credentialsName,
					Type:      "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
					Value:     payloadValue,
					ProjectId: projectId,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentials
				})

			credentials, _ = apiClient.CredentialsCreate(&AwsCredentialsCreatePayload{
				Name:      credentialsName,
				Value:     payloadValue,
				Type:      "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
				ProjectId: projectId,
			})
		})

		It("Should not get organization id", func() {
			// Verify organizationId was not called
			organizationIdCall.Times(0)
		})

		It("Should send POST request with projectId and without organizationId", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockCredentials))
		})
	})

	Describe("AwsCredentialsUpdate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall()

			payloadValue := AwsCredentialsValuePayload{
				RoleArn:  "role",
				Duration: 1,
			}

			httpCall = mockHttpClient.EXPECT().
				Patch("/credentials/"+mockCredentials.Id, &AwsCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentials
				})

			credentials, _ = apiClient.CredentialsUpdate(credentials.Id, &AwsCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send PATCH request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockCredentials))
		})
	})

	Describe("AwsCredentialsUpdate with ProjectId", func() {
		const projectId = "project-456"

		BeforeEach(func() {
			// Note: No organization ID call should be made

			payloadValue := AwsCredentialsValuePayload{
				RoleArn:  "role",
				Duration: 1,
			}

			httpCall = mockHttpClient.EXPECT().
				Patch("/credentials/"+mockCredentials.Id, &AwsCredentialsCreatePayload{
					Name:      credentialsName,
					Type:      "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
					Value:     payloadValue,
					ProjectId: projectId,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentials
				})

			credentials, _ = apiClient.CredentialsUpdate(mockCredentials.Id, &AwsCredentialsCreatePayload{
				Name:      credentialsName,
				Value:     payloadValue,
				Type:      "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
				ProjectId: projectId,
			})
		})

		It("Should not get organization id", func() {
			organizationIdCall.Times(0)
		})

		It("Should send PATCH request with projectId and without organizationId", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockCredentials))
		})
	})

	Describe("GcpCredentialsCreate", func() {
		const gcpRequestType = "GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT"
		mockGcpCredentials := mockCredentials
		mockGcpCredentials.Type = gcpRequestType
		BeforeEach(func() {
			mockOrganizationIdCall()

			payloadValue := GcpCredentialsValuePayload{
				ProjectId:         "projectId",
				ServiceAccountKey: "serviceAccountKey",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &GcpCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           gcpRequestType,
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockGcpCredentials
				})

			credentials, _ = apiClient.CredentialsCreate(&GcpCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  gcpRequestType,
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockGcpCredentials))
		})
	})

	Describe("AzureCredentialsCreate", func() {
		const azureRequestType = "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT"
		mockAzureCredentials := mockCredentials
		mockAzureCredentials.Type = azureRequestType
		BeforeEach(func() {

			mockOrganizationIdCall().Times(1)

			payloadValue := AzureCredentialsValuePayload{
				ClientId:       "fakeClientId",
				ClientSecret:   "fakeClientSecret",
				SubscriptionId: "fakeSubscriptionId",
				TenantId:       "fakeTenantId",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &AzureCredentialsCreatePayload{
					Name:           credentialsName,
					OrganizationId: organizationId,
					Type:           azureRequestType,
					Value:          payloadValue,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockAzureCredentials
				}).Times(1)

			credentials, _ = apiClient.CredentialsCreate(&AzureCredentialsCreatePayload{
				Name:  credentialsName,
				Value: payloadValue,
				Type:  azureRequestType,
			})
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockAzureCredentials))
		})
	})

	Describe("AzureCredentialsCreate with ProjectId", func() {
		const projectId = "project-abc"
		const azureRequestType = "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT"
		mockAzureCredentialsWithProject := mockCredentials
		mockAzureCredentialsWithProject.Type = azureRequestType

		BeforeEach(func() {
			// No organization ID call expected

			payloadValue := AzureCredentialsValuePayload{
				ClientId:       "fakeClientId",
				ClientSecret:   "fakeClientSecret",
				SubscriptionId: "fakeSubscriptionId",
				TenantId:       "fakeTenantId",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &AzureCredentialsCreatePayload{
					Name:      credentialsName,
					Type:      azureRequestType,
					Value:     payloadValue,
					ProjectId: projectId,
				},
					gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockAzureCredentialsWithProject
				}).Times(1)

			credentials, _ = apiClient.CredentialsCreate(&AzureCredentialsCreatePayload{
				Name:      credentialsName,
				Value:     payloadValue,
				Type:      azureRequestType,
				ProjectId: projectId,
			})
		})

		It("Should not call organization id", func() {
			organizationIdCall.Times(0)
		})

		It("Should send POST request with projectId", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(mockAzureCredentialsWithProject))
		})
	})

	Describe("CloudCredentialsDelete", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/credentials/"+mockCredentials.Id, nil)
			err = apiClient.CloudCredentialsDelete(mockCredentials.Id)
		})

		It("Should send DELETE request with project id", func() {
			httpCall.Times(1)
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("CloudCredentials", func() {
		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/credentials", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]Credentials) {
					*response = keys
				})
			credentials, _ = apiClient.CloudCredentials(mockCredentials.Id)
		})

		It("Should send GET request with project id", func() {
			httpCall.Times(1)
		})

		It("Should return correct key", func() {
			Expect(credentials).To(Equal(mockCredentials))
		})
	})

	Describe("CloudCredentialsList", func() {
		var credentials []Credentials

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/credentials", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]Credentials) {
					*response = keys
				})
			credentials, _ = apiClient.CloudCredentialsList()
		})

		It("Should send GET request with organization id param", func() {
			httpCall.Times(1)
		})

		It("Should return all credentials", func() {
			Expect(credentials).To(Equal(keys))
		})
	})
})
