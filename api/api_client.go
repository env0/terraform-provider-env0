package api

import (
	"errors"
	"os"
)

type ApiClient struct {
	httpClient           *HttpClient
	cachedOrganizationId string
}

func NewClientFromEnv() (*ApiClient, error) {
	apiKey := os.Getenv("ENV0_API_KEY")
	apiSecret := os.Getenv("ENV0_API_SECRET")

	if len(apiKey) == 0 {
		return nil, errors.New("ENV0_API_KEY must be specified in environment")
	}
	if len(apiSecret) == 0 {
		return nil, errors.New("ENV0_API_SECRET must be specified in environment")
	}

	result := &ApiClient{
		httpClient: newHttpClient(apiKey, apiSecret),
	}

	return result, nil
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
