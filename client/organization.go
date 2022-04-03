package client

import (
	"errors"
)

func (ac *ApiClient) Organization() (Organization, error) {
	var result []Organization
	err := ac.http.Get("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}
	if len(result) != 1 {
		return Organization{}, errors.New("server responded with too many organizations")
	}
	return result[0], nil
}

func (ac *ApiClient) organizationId() (string, error) {
	if ac.cachedOrganizationId != "" {
		return ac.cachedOrganizationId, nil
	}
	organization, err := ac.Organization()
	if err != nil {
		return "", nil
	}
	ac.cachedOrganizationId = organization.Id
	return ac.cachedOrganizationId, nil
}

func (ac *ApiClient) OrganizationPolicyUpdate(payload OrganizationPolicyUpdatePayload) (*Organization, error) {
	id, err := ac.organizationId()
	if err != nil {
		return nil, err
	}

	var result Organization
	if err := ac.http.Post("/organizations/"+id+"/policies", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
