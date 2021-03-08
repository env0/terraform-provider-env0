package env0apiclient

import (
	"encoding/json"
	"net/url"
)

func (self *ApiClient) postJSON(path string, request interface{}, response interface{}) error {
	serialized, err := json.Marshal(request)
	if err != nil {
		return err
	}
	body, err := self.post(path, serialized)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}

func (self *ApiClient) getJSON(path string, params url.Values, response interface{}) error {
	body, err := self.get(path, params)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}
