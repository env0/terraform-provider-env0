package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Kubernetes Credentials", func() {
	var credentials *Credentials

	Describe("KubernetesCredentialsCreate", func() {
		value := AzureAksValue{
			ClusterName:   "cc11",
			ResourceGroup: "rg11",
		}

		createPayload := KubernetesCredentialsCreatePayload{
			Name:  "n1",
			Type:  "K8S_AZURE_AKS_AUTH",
			Value: &value,
		}

		createPayloadWithOrganizationId := struct {
			OrganizationId string      `json:"organizationId"`
			Name           string      `json:"name"`
			Type           string      `json:"type"`
			Value          interface{} `json:"value"`
		}{
			OrganizationId: organizationId,
			Name:           createPayload.Name,
			Type:           string(createPayload.Type),
			Value:          createPayload.Value,
		}

		mockCredentials := Credentials{
			Id: "id111",
		}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Post("/credentials", &createPayloadWithOrganizationId, gomock.Any()).
				Do(func(path string, request interface{}, response *Credentials) {
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
})
