package client

type EnvironmentDiscoveryPutPayload struct {
	GlobPattern                           string           `json:"globPattern"`
	EnvironmentPlacement                  string           `json:"environmentPlacement"`
	WorkspaceNaming                       string           `json:"workspaceNaming"`
	AutoDeployByCustomGlob                string           `json:"autoDeployByCustomGlob,omitempty"`
	Repository                            string           `json:"repository"`
	TerraformVersion                      string           `json:"terraformVersion,omitempty"`
	OpentofuVersion                       string           `json:"opentofuVersion,omitempty"`
	TerragruntVersion                     string           `json:"terragruntVersion,omitempty"`
	TerragruntTfBinary                    string           `json:"terragruntTfBinary,omitempty"`
	IsTerragruntRunAll                    bool             `json:"isTerragruntRunAll"`
	Type                                  string           `json:"type"`
	TokenId                               string           `json:"tokenId,omitempty"`
	SshKeys                               []TemplateSshKey `json:"sshKeys,omitempty"`
	GithubInstallationId                  int              `json:"githubInstallationId,omitempty"`
	BitbucketClientKey                    string           `json:"bitbucketClientKey,omitempty"`
	IsAzureDevops                         bool             `json:"isAzureDevOps"`
	IsBitbucketServer                     bool             `json:"isBitbucketServer"`
	IsGitHubEnterprise                    bool             `json:"isGitHubEnterprise" tfschema:"is_github_enterprise"`
	IsGitLabEnterprise                    bool             `json:"isGitLabEnterprise" tfschema:"is_gitlab_enterprise"`
	Retry                                 TemplateRetry    `json:"retry"`
	RootPath                              string           `json:"rootPath"`
	CreateNewEnvironmentsFromPullRequests bool             `json:"createNewEnvironmentsFromPullRequests"`
}

type EnvironmentDiscoveryPayload struct {
	Id                                    string           `json:"id"`
	GlobPattern                           string           `json:"globPattern"`
	EnvironmentPlacement                  string           `json:"environmentPlacement"`
	WorkspaceNaming                       string           `json:"workspaceNaming"`
	AutoDeployByCustomGlob                string           `json:"autoDeployByCustomGlob"`
	Repository                            string           `json:"repository"`
	TerraformVersion                      string           `json:"terraformVersion"`
	OpentofuVersion                       string           `json:"opentofuVersion"`
	TerragruntVersion                     string           `json:"terragruntVersion"`
	TerragruntTfBinary                    string           `json:"terragruntTfBinary" tfschema:",omitempty"`
	IsTerragruntRunAll                    bool             `json:"isTerragruntRunAll"`
	Type                                  string           `json:"type"`
	TokenId                               string           `json:"tokenId"`
	SshKeys                               []TemplateSshKey `json:"sshKeys" tfschema:"-"`
	GithubInstallationId                  int              `json:"githubInstallationId"`
	BitbucketClientKey                    string           `json:"bitbucketClientKey"`
	IsAzureDevops                         bool             `json:"isAzureDevOps"`
	IsBitbucketServer                     bool             `json:"isBitbucketServer"`
	IsGitHubEnterprise                    bool             `json:"isGitHubEnterprise" tfschema:"is_github_enterprise"`
	IsGitLabEnterprise                    bool             `json:"isGitLabEnterprise" tfschema:"is_gitlab_enterprise"`
	Retry                                 TemplateRetry    `json:"retry" tfschema:"-"`
	RootPath                              string           `json:"rootPath"`
	CreateNewEnvironmentsFromPullRequests bool             `json:"createNewEnvironmentsFromPullRequests"`
}

func (client *ApiClient) PutEnvironmentDiscovery(projectId string, payload *EnvironmentDiscoveryPutPayload) (*EnvironmentDiscoveryPayload, error) {
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
