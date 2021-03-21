package api

import (
	"errors"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ApiClient struct {
	Endpoint             string
	ApiKey               string
	ApiSecret            string
	client               *resty.Client
	cachedOrganizationId string
}

func NewClientFromEnv() (*ApiClient, error) {
	result := &ApiClient{
		ApiKey:    os.Getenv("ENV0_API_KEY"),
		ApiSecret: os.Getenv("ENV0_API_SECRET"),
		Endpoint:  "https://api.env0.com/",
		client:    resty.New(),
	}
	result.normalizeEndpoint()
	if len(result.ApiKey) == 0 {
		return nil, errors.New("ENV0_API_KEY must be specified in environment")
	}
	if len(result.ApiSecret) == 0 {
		return nil, errors.New("ENV0_API_SECRET must be specified in environment")
	}
	return result, nil
}

func (self *ApiClient) normalizeEndpoint() {
	for strings.HasSuffix(self.Endpoint, "/") {
		self.Endpoint = self.Endpoint[:len(self.Endpoint)-1]
	}
}
