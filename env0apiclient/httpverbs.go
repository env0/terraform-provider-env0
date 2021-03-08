package env0apiclient

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

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return body, nil
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

	if resp.StatusCode != 200 {
		////////////
		ioutil.WriteFile("/tmp/log1", body, 0644)
		///////////
		return nil, errors.New(resp.Status)
	}

	return body, nil
}
