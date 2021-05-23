package client

import (
	"errors"
	"net/http"
)

type ApiClient struct {
	Endpoint             string
	ApiKey               string
	ApiSecret            string
	client               *http.Client
	cachedOrganizationId string
}

func (self *ApiClient) Organization() (Organization, error) {
	var result []Organization
	err := self.getJSON("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}
	if len(result) != 1 {
		return Organization{}, errors.New("Server responded with a too many organizations")
	}
	return result[0], nil
}
