package http

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

type HttpClientInterface interface {
	Get(path string, params map[string]string) (interface{}, error)
	Post(path string, request interface{}) (interface{}, error)
	Put(path string, request interface{}) (interface{}, error)
	Delete(path string) error
}

type HttpClient struct {
	ApiKey    string
	ApiSecret string
	Endpoint  string
	client    *resty.Client
}

func NewHttpClient(apiKey string, apiSecret string, apiEndpoint string) (*HttpClient, error) {
	return &HttpClient{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		client:    resty.New().SetHostURL(apiEndpoint),
	}, nil
}

func (self *HttpClient) request() *resty.Request {
	return self.client.R().SetBasicAuth(self.ApiKey, self.ApiSecret)
}

func (self *HttpClient) httpResult(response *resty.Response, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		return nil, errors.New(response.Status() + ": " + string(response.Body()))
	}
	return response.Result(), nil
}

func (self *HttpClient) Post(path string, request interface{}) (interface{}, error) {
	result, err := self.request().
		SetBody(request).
		Post(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Put(path string, request interface{}) (interface{}, error) {
	result, err := self.request().
		SetBody(request).
		Put(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Get(path string, params map[string]string) (interface{}, error) {
	result, err := self.request().
		SetQueryParams(params).
		Get(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) Delete(path string) error {
	result, err := self.request().Delete(path)
	_, httpErr := self.httpResult(result, err)

	return httpErr
}
