package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (self *ApiClient) normalizeEndpoint() {
	for strings.HasSuffix(self.Endpoint, "/") {
		self.Endpoint = self.Endpoint[:len(self.Endpoint)-1]
	}
}

func (self *ApiClient) post(path string, payload []byte) ([]byte, error) {
	self.normalizeEndpoint()
	req, err := http.NewRequest(http.MethodPost, self.Endpoint+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return self.do(req)
}

func (self *ApiClient) put(path string, payload []byte) ([]byte, error) {
	self.normalizeEndpoint()
	req, err := http.NewRequest(http.MethodPut, self.Endpoint+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return self.do(req)
}

func (self *ApiClient) get(path string, params url.Values) ([]byte, error) {
	self.normalizeEndpoint()
	if params != nil {
		path += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, self.Endpoint+path, nil)
	if err != nil {
		return nil, err
	}
	return self.do(req)
}

func (self *ApiClient) delete(path string) error {
	self.normalizeEndpoint()
	req, err := http.NewRequest(http.MethodDelete, self.Endpoint+path, nil)
	if err != nil {
		return err
	}
	_, err = self.do(req)
	return err
}

func (self *ApiClient) do(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(self.ApiKey, self.ApiSecret)
	req.Header.Add("Accept", "application/json")

	if self.client == nil {
		self.client = &http.Client{}
	}
	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 204 {
		return nil, nil
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New(resp.Status + ": " + string(body))
	}

	return body, nil
}
