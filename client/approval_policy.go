package client

import "fmt"

type ApprovalPolicy struct {
	Id                   string           `json:"id"`
	Name                 string           `json:"name"`
	Repository           string           `json:"repository"`
	Path                 string           `json:"path" tfschema:",omitempty"`
	Revision             string           `json:"revision" tfschema:",omitempty"`
	TokenId              string           `json:"tokenId" tfschema:",omitempty"`
	SshKeys              []TemplateSshKey `json:"sshKeys"`
	GithubInstallationId int              `json:"githubInstallationId" tfschema:",omitempty"`
	BitbucketClientKey   string           `json:"bitbucketClientKey" tfschema:",omitempty"`
	IsBitbucketServer    bool             `json:"isBitbucketServer"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	IsGithubEnterprise   bool             `json:"isGitHubEnterprise"`
	IsGitLab             bool             `json:"isGitLab" tfschema:"is_gitlab"`
	IsAzureDevOps        bool             `json:"isAzureDevOps" tfschema:"is_azure_devops"`
	IsTerragruntRunAll   bool             `json:"isTerragruntRunAll"`
}

type ApprovalPolicyByScope struct {
	Scope          string          `json:"scope"`
	ScopeId        string          `json:"scopeId"`
	ApprovalPolicy *ApprovalPolicy `json:"blueprint"`
}

type ApprovalPolicyCreatePayload struct {
	Name                 string `json:"name"`
	Repository           string `json:"repository"`
	Revision             string `json:"revision,omitempty"`
	Path                 string `json:"path,omitempty"`
	TokenId              string `json:"tokenId,omitempty"`
	GithubInstallationId int    `json:"githubInstallationId,omitempty"`
	BitbucketClientKey   string `json:"bitbucketClientKey,omitempty"`
	IsBitbucketServer    bool   `json:"isBitbucketServer,omitempty"`
	IsGitlabEnterprise   bool   `json:"isGitLabEnterprise"`
	IsGithubEnterprise   bool   `json:"isGitHubEnterprise"`
	IsGitLab             bool   `json:"isGitLab" tfschema:"is_gitlab"`
	IsAzureDevOps        bool   `json:"isAzureDevOps" tfschema:"is_azure_devops"`
}

type ApprovalPolicyUpdatePayload struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Repository           string `json:"repository"`
	Revision             string `json:"revision,omitempty"`
	Path                 string `json:"path,omitempty"`
	TokenId              string `json:"tokenId,omitempty"`
	GithubInstallationId int    `json:"githubInstallationId,omitempty"`
	BitbucketClientKey   string `json:"bitbucketClientKey,omitempty"`
	IsBitbucketServer    bool   `json:"isBitbucketServer,omitempty"`
	IsGitlabEnterprise   bool   `json:"isGitLabEnterprise"`
	IsGithubEnterprise   bool   `json:"isGitHubEnterprise"`
	IsGitLab             bool   `json:"isGitLab" tfschema:"is_gitlab"`
	IsAzureDevOps        bool   `json:"isAzureDevOps" tfschema:"is_azure_devops"`
}

const (
	ApprovalPolicyProjectScope  string = "PROJECT"
	ApprovalPolicyTemplateScope string = "BLUEPRINT"
)

type ApprovalPolicyAssignment struct {
	Scope       string `json:"scope"`
	ScopeId     string `json:"scopeId"`
	BlueprintId string `json:"blueprintId"`
}

func (client *ApiClient) ApprovalPolicyCreate(payload *ApprovalPolicyCreatePayload) (*ApprovalPolicy, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := struct {
		ApprovalPolicyCreatePayload
		OrganizationId string `json:"organizationId"`
	}{
		*payload,
		organizationId,
	}

	var result ApprovalPolicy
	if err := client.http.Post("/approval-policy", &payloadWithOrganizationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ApprovalPolicyUpdate(payload *ApprovalPolicyUpdatePayload) (*ApprovalPolicy, error) {
	var result ApprovalPolicy
	if err := client.http.Put("/approval-policy", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ApprovalPolicies(name string) ([]ApprovalPolicy, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []ApprovalPolicy
	if err := client.http.Get("/approval-policy", map[string]string{"organizationId": organizationId, "name": name}, &result); err != nil {
		return nil, err
	}

	return result, err
}

func (client *ApiClient) ApprovalPolicyAssign(assignment *ApprovalPolicyAssignment) (*ApprovalPolicyAssignment, error) {
	var result ApprovalPolicyAssignment

	if err := client.http.Post("/approval-policy/assignment", assignment, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ApprovalPolicyUnassign(id string) error {
	return client.http.Delete("/approval-policy/assignment", map[string]string{"id": id})
}

func (client *ApiClient) ApprovalPolicyByScope(scope string, scopeId string) ([]ApprovalPolicyByScope, error) {
	var result []ApprovalPolicyByScope

	if err := client.http.Get(fmt.Sprintf("/approval-policy/%s/%s", scope, scopeId), nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
