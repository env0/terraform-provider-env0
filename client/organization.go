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

func (client *ApiClient) OrganizationId() (string, error) {
	if client.cachedOrganizationId != "" {
		return client.cachedOrganizationId, nil
	}
	organization, err := client.Organization()
	if err != nil {
		return "", err
	}
	client.cachedOrganizationId = organization.Id
	return client.cachedOrganizationId, nil
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
