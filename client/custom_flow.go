package client

type CustomFlowCreatePayload struct {
	Name                 string           `json:"name"`
	Repository           string           `json:"repository"`
	Path                 string           `json:"path,omitempty"`
	Revision             string           `json:"revision,omitempty"`
	TokenId              string           `json:"tokenId,omitempty"`
	SshKeys              []TemplateSshKey `json:"sshKeys,omitempty"`
	GithubInstallationId int              `json:"githubInstallationId,omitempty"`
	BitbucketClientKey   string           `json:"bitbucketClientKey,omitempty"`
	IsBitbucketServer    bool             `json:"isBitbucketServer"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	IsGithubEnterprise   bool             `json:"isGitHubEnterprise"`
	IsGitLab             bool             `json:"isGitLab"`
	IsAzureDevOps        bool             `json:"isAzureDevOps" tfschema:"is_azure_devops"`
	IsTerragruntRunAll   bool             `json:"isTerragruntRunAll"`
}

type CustomFlow struct {
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
	IsGitLab             bool             `json:"isGitLab"`
	IsAzureDevOps        bool             `json:"isAzureDevOps" tfschema:"is_azure_devops"`
	IsTerragruntRunAll   bool             `json:"isTerragruntRunAll"`
}

func (client *ApiClient) CustomFlowCreate(payload CustomFlowCreatePayload) (*CustomFlow, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganzationId := struct {
		OrganizationId string `json:"organizationId"`
		CustomFlowCreatePayload
	}{
		organizationId,
		payload,
	}

	var result CustomFlow
	if err := client.http.Post("/custom-flow", &payloadWithOrganzationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) CustomFlow(id string) (*CustomFlow, error) {
	var result CustomFlow

	if err := client.http.Get("/custom-flow/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) CustomFlowDelete(id string) error {
	return client.http.Delete("/custom-flow/" + id)
}

func (client *ApiClient) CustomFlowUpdate(id string, payload CustomFlowCreatePayload) (*CustomFlow, error) {
	payloadWithId := struct {
		Id string `json:"id"`
		CustomFlowCreatePayload
	}{
		id,
		payload,
	}

	var result CustomFlow
	if err := client.http.Put("/custom-flow", &payloadWithId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) CustomFlows(name string) ([]CustomFlow, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []CustomFlow
	if err := client.http.Get("/custom-flows", map[string]string{"organizationId": organizationId, "name": name}, &result); err != nil {
		return nil, err
	}

	return result, err
}
