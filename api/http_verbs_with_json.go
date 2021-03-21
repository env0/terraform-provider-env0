package api

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

func (self *ApiClient) postJSON(path string, request interface{}, response interface{}) error {
	return self.executeJSONInJSONOut(resty.MethodPost, path, request, response)
}

func (self *ApiClient) putJSON(path string, request interface{}, response interface{}) error {
	return self.executeJSONInJSONOut(resty.MethodPut, path, request, response)
}

func (self *ApiClient) executeJSONInJSONOut(method string, path string, request interface{}, response interface{}) error {
	self.normalizeEndpoint()
	result, err := self.client.R().
		SetBasicAuth(self.ApiKey, self.ApiSecret).
		SetBody(request).
		SetResult(response).
		Execute(method, self.Endpoint+path)
	if err != nil {
		return err
	}
	if result.StatusCode() < 200 || result.StatusCode() > 299 {
		return errors.New(result.Status() + ": " + string(result.Body()))
	}
	return nil
}

func (self *ApiClient) getJSON(path string, params map[string]string, response interface{}) error {
	self.normalizeEndpoint()
	result, err := self.client.R().
		SetBasicAuth(self.ApiKey, self.ApiSecret).
		SetQueryParams(params).
		SetResult(response).
		Get(self.Endpoint + path)
	if err != nil {
		return err
	}
	if result.StatusCode() < 200 || result.StatusCode() > 299 {
		return errors.New(result.Status() + ": " + string(result.Body()))
	}
	return nil
}

func (self *ApiClient) delete(path string) error {
	self.normalizeEndpoint()
	result, err := self.client.R().
		SetBasicAuth(self.ApiKey, self.ApiSecret).
		Delete(self.Endpoint + path)
	if err != nil {
		return err
	}
	if result.StatusCode() < 200 || result.StatusCode() > 299 {
		return errors.New(result.Status() + ": " + string(result.Body()))
	}
	return nil
}
