package client

import (
	"github.com/env0/terraform-provider-env0/client/http"
)

type ApiClient struct {
	client               http.HttpClientInterface
	cachedOrganizationId string
}

func NewApiClient(client http.HttpClientInterface, organizationId string) *ApiClient {
	return &ApiClient{
		client:               client,
		cachedOrganizationId: organizationId,
	}
}
