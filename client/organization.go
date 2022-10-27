package client

import (
	"errors"
	"sync"
)

var cachedOrganizationLock sync.Mutex

type Organization struct {
	Id                                  string  `json:"id"`
	Name                                string  `json:"name"`
	MaxTtl                              *string `json:"maxTtl"`
	DefaultTtl                          *string `json:"defaultTtl"`
	DoNotReportSkippedStatusChecks      bool    `json:"doNotReportSkippedStatusChecks"`
	DoNotConsiderMergeCommitsForPrPlans bool    `json:"doNotConsiderMergeCommitsForPrPlans"`
	EnableOidc                          bool    `json:"enableOidc"`
	Description                         string  `json:"description"`
	PhotoUrl                            string  `json:"photoUrl"`
	CreatedBy                           string  `json:"createdBy"`
	CreatedAt                           string  `json:"createdAt"`
	UpdatedAt                           string  `json:"updatedAt"`
	Role                                string  `json:"role"`
	IsSelfHostedK8s                     bool    `json:"isSelfHostedK8s" tfschema:"is_self_hosted"`
}

type OrganizationPolicyUpdatePayload struct {
	MaxTtl                              *string `json:"maxTtl"`
	DefaultTtl                          *string `json:"defaultTtl"`
	DoNotReportSkippedStatusChecks      *bool   `json:"doNotReportSkippedStatusChecks,omitempty"`
	DoNotConsiderMergeCommitsForPrPlans *bool   `json:"doNotConsiderMergeCommitsForPrPlans,omitempty"`
	EnableOidc                          *bool   `json:"enableOidc,omitempty"`
}

func (client *ApiClient) Organization() (Organization, error) {
	var result []Organization
	err := client.http.Get("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}
	if len(result) != 1 {
		return Organization{}, errors.New("server responded with too many organizations")
	}
	return result[0], nil
}

func (client *ApiClient) getCachedOrganization() (*Organization, error) {
	cachedOrganizationLock.Lock()
	defer cachedOrganizationLock.Unlock()

	if client.cachedOrganization != nil {
		return client.cachedOrganization, nil
	}

	organization, err := client.Organization()
	if err != nil {
		return nil, err
	}

	client.cachedOrganization = &organization

	return client.cachedOrganization, nil
}

func (client *ApiClient) OrganizationId() (string, error) {
	organization, err := client.getCachedOrganization()
	if err != nil {
		return "", err
	}
	return organization.Id, nil
}

func (client *ApiClient) IsOrganizationSelfHostedAgent() (bool, error) {
	organization, err := client.getCachedOrganization()
	if err != nil {
		return false, err
	}

	return organization.IsSelfHostedK8s, nil
}

func (client *ApiClient) OrganizationPolicyUpdate(payload OrganizationPolicyUpdatePayload) (*Organization, error) {
	id, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result Organization
	if err := client.http.Post("/organizations/"+id+"/policies", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
