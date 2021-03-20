package api

import (
	"errors"
	"net/http"
	"os"
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

func NewClientFromEnv() (*ApiClient, error) {
	result := &ApiClient{
		ApiKey: os.Getenv("ENV0_API_KEY"),
		ApiSecret: os.Getenv("ENV0_API_SECRET"),
		Endpoint: "https://api.env0.com/",
	}
	if len(result.ApiKey) == 0 {
		return nil, errors.New("ENV0_API_KEY must be specified in environment")
	}
	if len(result.ApiSecret) == 0 {
		return nil, errors.New("ENV0_API_SECRET must be specified in environment")
	}
	return result, nil
}