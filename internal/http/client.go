package http

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"os"
)

type HttpClientInterface interface {
	Get(path string, params map[string]string, response interface{}) error
	Post(path string, request interface{}, response interface{}) error
	Put(path string, request interface{}, response interface{}) error
	Delete(path string) error
}

type HttpClient struct {
	ApiKey    string
	ApiSecret string
	Endpoint  string
	client    *resty.Client
}

func NewRestClientFromEnv() (*HttpClient, error) {
	apiKey := os.Getenv("ENV0_API_KEY")
	apiSecret := os.Getenv("ENV0_API_SECRET")

	if len(apiKey) == 0 {
		return nil, errors.New("ENV0_API_KEY must be specified in environment")
	}
	if len(apiSecret) == 0 {
		return nil, errors.New("ENV0_API_SECRET must be specified in environment")
	}

	return &HttpClient{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		client:    resty.New().SetHostURL("https://api.env0.com/"),
	}, nil
}

func (self *HttpClient) request() *resty.Request {
	return self.client.R().SetBasicAuth(self.ApiKey, self.ApiSecret)
}

func (self *HttpClient) httpResult(response *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if response.StatusCode() < 200 || response.StatusCode() > 299 {
		return errors.New(response.Status() + ": " + string(response.Body()))
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
