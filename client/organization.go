package client

import (
	"errors"
)

func (self *ApiClient) Organization() (Organization, error) {
	var result []Organization
	err := self.http.Get("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}
	if len(result) != 1 {
		return Organization{}, errors.New("Server responded with too many organizations")
	}
	return result[0], nil
}

func (self *ApiClient) organizationId() (string, error) {
	if self.cachedOrganizationId != "" {
		return self.cachedOrganizationId, nil
	}
	organization, err := self.Organization()
	if err != nil {
		return "", nil
	}
	self.cachedOrganizationId = organization.Id
	return self.cachedOrganizationId, nil
}
