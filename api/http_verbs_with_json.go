package api

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

func (self *ApiClient) request() *resty.Request {
	self.normalizeEndpoint()
	return self.client.R().SetBasicAuth(self.ApiKey, self.ApiSecret)
}

func httpResult(response *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if response.StatusCode() < 200 || response.StatusCode() > 299 {
		return errors.New(response.Status() + ": " + string(response.Body()))
	}
	return nil
}

func (self *ApiClient) postJSON(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Post(path)
	return httpResult(result, err)
}

func (self *ApiClient) putJSON(path string, request interface{}, response interface{}) error {
	result, err := self.request().
		SetBody(request).
		SetResult(response).
		Put(path)
	return httpResult(result, err)
}

func (self *ApiClient) getJSON(path string, params map[string]string, response interface{}) error {
	result, err := self.request().
		SetQueryParams(params).
		SetResult(response).
		Get(path)
	return httpResult(result, err)
}

func (self *ApiClient) delete(path string) error {
	result, err := self.request().Delete(path)
	return httpResult(result, err)
}
