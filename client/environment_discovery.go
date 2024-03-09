package client

type EnvironmentDiscoveryPutPayload struct {
	GlobPattern            string           `json:"globPattern"`
	EnvironmentPlacement   string           `json:"environmentPlacement"`
	WorkspaceNaming        string           `json:"workspaceNaming"`
	AutoDeployByCustomGlob string           `json:"autoDeployByCustomGlob,omitempty"`
	Repository             string           `json:"repository"`
	TerraformVersion       string           `json:"terraformVersion,omitempty"`
	OpentofuVersion        string           `json:"opentofuVersion,omitempty"`
	TerragruntVersion      string           `json:"terragruntVersion,omitempty"`
	TerragruntTfBinary     string           `json:"terragruntTfBinary,omitempty"`
	Type                   string           `json:"type"`
	GitlabProjectId        int              `json:"gitlabProjectId,omitempty"`
	TokenId                string           `json:"tokenId,omitempty"`
	SshKeys                []TemplateSshKey `json:"sshKeys,omitempty"`
	GithubInstallationId   int              `json:"githubInstallationId,omitempty"`
	BitbucketClientKey     string           `json:"bitbucketClientKey,omitempty"`
	IsAzureDevops          bool             `json:"isAzureDevOps"`
	Retry                  TemplateRetry    `json:"retry"`
}

type EnvironmentDiscoveryPayload struct {
	Id                     string           `json:"id"`
	GlobPattern            string           `json:"globPattern"`
	EnvironmentPlacement   string           `json:"environmentPlacement"`
	WorkspaceNaming        string           `json:"workspaceNaming"`
	AutoDeployByCustomGlob string           `json:"autoDeployByCustomGlob" tfschema:",omitempty"`
	Repository             string           `json:"repository"`
	TerraformVersion       string           `json:"terraformVersion" tfschema:",omitempty"`
	OpentofuVersion        string           `json:"opentofuVersion" tfschema:",omitempty"`
	TerragruntVersion      string           `json:"terragruntVersion" tfschema:",omitempty"`
	TerragruntTfBinary     string           `json:"terragruntTfBinary" tfschema:",omitempty"`
	Type                   string           `json:"type"`
	GitlabProjectId        int              `json:"gitlabProjectId" tfschema:",omitempty"`
	TokenId                string           `json:"tokenId" tfschema:",omitempty"`
	SshKeys                []TemplateSshKey `json:"sshKeys" tfschema:"-"`
	GithubInstallationId   int              `json:"githubInstallationId" tfschema:",omitempty"`
	BitbucketClientKey     string           `json:"bitbucketClientKey" tfschema:",omitempty"`
	IsAzureDevops          bool             `json:"isAzureDevOps" tfschema:",omitempty"`
	Retry                  TemplateRetry    `json:"retry" tfschema:"-"`
}

func (client *ApiClient) EnableUpdateEnvironmentDiscovery(projectId string, payload *EnvironmentDiscoveryPutPayload) (*EnvironmentDiscoveryPayload, error) {
	var result EnvironmentDiscoveryPayload

	if err := client.http.Put("/environment-discovery/projects/"+projectId, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) GetEnvironmentDiscovery(projectId string) (*EnvironmentDiscoveryPayload, error) {
	var result EnvironmentDiscoveryPayload

	if err := client.http.Get("/environment-discovery/projects/"+projectId, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) DeleteEnvironmentDiscovery(projectId string) error {
	return client.http.Delete("/environment-discovery/projects/"+projectId, nil)
}
