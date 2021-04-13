package api

import (
	"errors"
)

func (self *ApiClient) Organization() (Organization, error) {
	var result []Organization
	err := self.client.Get("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}
	if len(result) != 1 {
		return Organization{}, errors.New("Server responded with a too many organizations")
	}
	return result[0], nil
}
