package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Kubernetes Credentials", func() {
	var credentials *Credentials

	mockCredentials := Credentials{
		Id: "id111",
	}

	value := AzureAksValue{
		ClusterName:   "cc11",
		ResourceGroup: "rg11",
	}

	Describe("KubernetesCredentialsCreate", func() {
		createPayload := KubernetesCredentialsCreatePayload{
			Name:  "n1",
			Type:  "K8S_AZURE_AKS_AUTH",
			Value: &value,
		}

		createPayloadWithOrganizationId := struct {
			OrganizationId string `json:"organizationId"`
			Name           string `json:"name"`
			Type           string `json:"type"`
			Value          any    `json:"value"`
		}{
			OrganizationId: organizationId,
			Name:           createPayload.Name,
			Type:           string(createPayload.Type),
			Value:          createPayload.Value,
		}

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &createPayloadWithOrganizationId, gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentials
				})

			credentials, _ = apiClient.KubernetesCredentialsCreate(&createPayload)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(&mockCredentials))
		})
	})

	Describe("KubernetesCredentialsUpdate", func() {
		updatePayload := KubernetesCredentialsUpdatePayload{
			Type:  "K8S_AZURE_AKS_AUTH",
			Value: value,
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Patch("/credentials/"+mockCredentials.Id, &updatePayload, gomock.Any()).
				Do(func(path string, request any, response *Credentials) {
					*response = mockCredentials
				})

			credentials, _ = apiClient.KubernetesCredentialsUpdate(mockCredentials.Id, &updatePayload)
		})

		It("Should send PATCH request with params", func() {
			httpCall.Times(1)
		})

		It("Should return key", func() {
			Expect(credentials).To(Equal(&mockCredentials))
		})
	})
})
