package api

import (
	. "github.com/env0/terraform-provider-env0/internal/http"
)

type ApiClient struct {
	client         HttpClientInterface
	organizationId string
}

func NewApiClient(client HttpClientInterface, organizationId string) *ApiClient {
	return &ApiClient{
		client:         client,
		organizationId: organizationId,
	}
}

func (self *ApiClient) getOrganizationId() (string, error) {
	if self.organizationId != "" {
		return self.organizationId, nil
	}
	organization, err := self.Organization()
	if err != nil {
		return "", nil
	}
	self.organizationId = organization.Id
	return self.organizationId, nil
}
