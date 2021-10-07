package client

import "errors"

func (self *ApiClient) Policy() (Policy, error) {
	var result []Policy
	err := self.http.Get("/policies", nil, &result)
	if err != nil {
		return Policy{}, err
	}
	if len(result) != 1 {
		return Policy{}, errors.New("Server responded with too many policies")
	}
	return result[0], nil
}
