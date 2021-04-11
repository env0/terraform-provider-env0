package api

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

type HttpClient struct {
	ApiKey    string
	ApiSecret string
	Endpoint  string
	client    *resty.Client
}

func newHttpClient(apiKey string, apiSecret string) *HttpClient {
	return &HttpClient{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		client:    resty.New().SetHostURL("https://api.env0.com/"),
	}
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

func (self *HttpClient) postJSON(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Post(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) putJSON(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Put(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) getJSON(path string, params map[string]string, response interface{}) error {
	result, err := self.request().
		SetQueryParams(params).
		SetResult(response).
		Get(path)
	return self.httpResult(result, err)
}

func (self *HttpClient) delete(path string) error {
	result, err := self.request().Delete(path)
	return self.httpResult(result, err)
}
