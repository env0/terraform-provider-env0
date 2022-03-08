package http

//go:generate mockgen -destination=client_mock.go -package=http . HttpClientInterface

import (
	"github.com/go-resty/resty/v2"
)

type HttpClientInterface interface {
	Get(path string, params map[string]string, response interface{}) error
	Post(path string, request interface{}, response interface{}) error
	Put(path string, request interface{}, response interface{}) error
	Delete(path string) error
	Patch(path string, request interface{}, response interface{}) error
}

type HttpClient struct {
	ApiKey    string
	ApiSecret string
	Endpoint  string
	client    *resty.Client
}

type HttpClientConfig struct {
	ApiKey      string
	ApiSecret   string
	ApiEndpoint string
	UserAgent   string
	RestClient  *resty.Client
}

func NewHttpClient(config HttpClientConfig) (*HttpClient, error) {
	return &HttpClient{
		ApiKey:    config.ApiKey,
		ApiSecret: config.ApiSecret,
		client:    config.RestClient.SetHostURL(config.ApiEndpoint).SetHeader("User-Agent", config.UserAgent),
	}, nil
}

func (self *HttpClient) request() *resty.Request {
	return self.client.R().SetBasicAuth(self.ApiKey, self.ApiSecret)
}

func (self *HttpClient) httpResult(response *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return &FailedResponseError{res: response}
	}
	return nil
}

func (self *HttpClient) Post(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Post(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Put(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Put(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Get(path string, params map[string]string, response interface{}) error {
	result, err := self.request().
		SetQueryParams(params).
		SetResult(response).
		Get(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Delete(path string) error {
	result, err := self.request().Delete(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Patch(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Patch(path)
	return self.httpResult(result, err)
}
