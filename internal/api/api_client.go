package api

import (
	. "github.com/env0/terraform-provider-env0/internal/http"
)

type ApiClient struct {
	client               HttpClientInterface
	cachedOrganizationId string
}

func NewApiClient(client HttpClientInterface) *ApiClient {
	return &ApiClient{
		client: client,
	}
}

func (self *ApiClient) organizationId() (string, error) {
	if self.cachedOrganizationId != "" {
		return self.cachedOrganizationId, nil
	}
	organization, err := self.Organization()
	if err != nil {
		return "", nil
	}
	self.cachedOrganizationId = organization.Id
	return self.cachedOrganizationId, nil
}
