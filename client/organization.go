package client

import (
	"errors"
)

func (self *ApiClient) Organization() (Organization, error) {
	result, err := self.http.Get("/organizations", nil)
	if err != nil {
		return Organization{}, err
	}

	var organizations []Organization
	organizations = result.([]Organization)

	if len(organizations) != 1 {
		return Organization{}, errors.New("Server responded with too many organizations")
	}
	return organizations[0], nil
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
