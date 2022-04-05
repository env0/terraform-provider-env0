package client

import (
	"errors"
)

type Organization struct {
	Id                                  string  `json:"id"`
	Name                                string  `json:"name"`
	MaxTtl                              *string `json:"maxTtl"`
	DefaultTtl                          *string `json:"defaultTtl"`
	DoNotReportSkippedStatusChecks      bool    `json:"doNotReportSkippedStatusChecks"`
	DoNotConsiderMergeCommitsForPrPlans bool    `json:"doNotConsiderMergeCommitsForPrPlans"`
	Description                         string  `json:"description"`
	PhotoUrl                            string  `json:"photoUrl"`
	CreatedBy                           string  `json:"createdBy"`
	CreatedAt                           string  `json:"createdAt"`
	UpdatedAt                           string  `json:"updatedAt"`
	Role                                string  `json:"role"`
	IsSelfHosted                        bool    `json:"isSelfHosted"`
	IsSelfHostedK8s                     bool    `json:"isSelfHostedK8s"`
}

type OrganizationPolicyUpdatePayload struct {
	MaxTtl                              *string `json:"maxTtl"`
	DefaultTtl                          *string `json:"defaultTtl"`
	DoNotReportSkippedStatusChecks      *bool   `json:"doNotReportSkippedStatusChecks,omitempty"`
	DoNotConsiderMergeCommitsForPrPlans *bool   `json:"doNotConsiderMergeCommitsForPrPlans,omitempty"`
}

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
