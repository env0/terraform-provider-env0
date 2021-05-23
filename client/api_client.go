package client

import (
	"github.com/env0/terraform-provider-env0/client/http"
)

type ApiClient struct {
	http                 http.HttpClientInterface
	cachedOrganizationId string
}

func NewApiClient(client http.HttpClientInterface) *ApiClient {
	return &ApiClient{
		http:                 client,
		cachedOrganizationId: "",
	}
}
