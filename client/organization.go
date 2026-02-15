package client

import (
	"errors"
	"fmt"
)

type Organization struct {
	Id                                  string  `json:"id"`
	Name                                string  `json:"name"`
	MaxTtl                              *string `json:"maxTtl"`
	DefaultTtl                          *string `json:"defaultTtl"`
	DoNotReportSkippedStatusChecks      bool    `json:"doNotReportSkippedStatusChecks"`
	DoNotConsiderMergeCommitsForPrPlans bool    `json:"doNotConsiderMergeCommitsForPrPlans"`
	EnableOidc                          bool    `json:"enableOidc"`
	EnforcePrCommenterPermissions       bool    `json:"enforcePrCommenterPermissions"`
	Description                         string  `json:"description"`
	PhotoUrl                            string  `json:"photoUrl"`
	CreatedBy                           string  `json:"createdBy"`
	CreatedAt                           string  `json:"createdAt"`
	UpdatedAt                           string  `json:"updatedAt"`
	Role                                string  `json:"role"`
	IsSelfHostedK8s                     bool    `json:"isSelfHostedK8s"                     tfschema:"is_self_hosted"`
}

type OrganizationPolicyUpdatePayload struct {
	MaxTtl                              *string `json:"maxTtl"`
	DefaultTtl                          *string `json:"defaultTtl"`
	DoNotReportSkippedStatusChecks      *bool   `json:"doNotReportSkippedStatusChecks,omitempty"`
	DoNotConsiderMergeCommitsForPrPlans *bool   `json:"doNotConsiderMergeCommitsForPrPlans,omitempty"`
	EnableOidc                          *bool   `json:"enableOidc,omitempty"`
	EnforcePrCommenterPermissions       *bool   `json:"enforcePrCommenterPermissions,omitempty"`
}

func (client *ApiClient) Organization() (Organization, error) {
	var result []Organization

	err := client.http.Get("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}

	if len(result) != 1 {
		if client.defaultOrganizationId != "" {
			for _, organization := range result {
				if organization.Id == client.defaultOrganizationId {
					return organization, nil
				}
			}

			return Organization{}, fmt.Errorf("the api key is not assigned to organization id: %s", client.defaultOrganizationId)
		}

		return Organization{}, errors.New("server responded with too many organizations (set a default organization id in the provider settings)")
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

	if payload.DefaultTtl != nil && *payload.DefaultTtl == "" {
		payload.DefaultTtl = nil
	}

	if payload.MaxTtl != nil && *payload.MaxTtl == "" {
		payload.MaxTtl = nil
	}

	var result Organization
	if err := client.http.Post("/organizations/"+id+"/policies", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) OrganizationUserUpdateRole(userId string, roleId string) error {
	id, err := client.OrganizationId()
	if err != nil {
		return err
	}

	return client.http.Put("/organizations/"+id+"/users/"+userId+"/role", roleId, nil)
}
