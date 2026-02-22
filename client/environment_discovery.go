package client

import "errors"

type DiscoveryFileConfiguration struct {
	RepositoryRegex string `json:"repositoryRegex"`
}

type EnvironmentDiscoveryPutPayload struct {
	GlobPattern                           string                      `json:"globPattern,omitempty"`
	EnvironmentPlacement                  string                      `json:"environmentPlacement,omitempty"`
	WorkspaceNaming                       string                      `json:"workspaceNaming,omitempty"`
	AutoDeployByCustomGlob                string                      `json:"autoDeployByCustomGlob,omitempty"`
	Repository                            string                      `json:"repository,omitempty"`
	TerraformVersion                      string                      `json:"terraformVersion,omitempty"`
	OpentofuVersion                       string                      `json:"opentofuVersion,omitempty"`
	TerragruntVersion                     string                      `json:"terragruntVersion,omitempty"`
	TerragruntTfBinary                    string                      `json:"terragruntTfBinary,omitempty"`
	IsTerragruntRunAll                    bool                        `json:"isTerragruntRunAll"`
	Type                                  string                      `json:"type,omitempty"`
	TokenId                               string                      `json:"tokenId,omitempty"`
	SshKeys                               []TemplateSshKey            `json:"sshKeys,omitempty"`
	GithubInstallationId                  int                         `json:"githubInstallationId,omitempty"`
	VcsConnectionId                       string                      `json:"vcsConnectionId,omitempty"`
	BitbucketClientKey                    string                      `json:"bitbucketClientKey,omitempty"`
	IsAzureDevops                         bool                        `json:"isAzureDevOps,omitempty"`
	IsBitbucketServer                     bool                        `json:"isBitbucketServer,omitempty"`
	IsGitHubEnterprise                    bool                        `json:"isGitHubEnterprise,omitempty"                    tfschema:"is_github_enterprise"`
	IsGitLabEnterprise                    bool                        `json:"isGitLabEnterprise,omitempty"                    tfschema:"is_gitlab_enterprise"`
	Retry                                 TemplateRetry               `json:"retry,omitempty"`
	RootPath                              string                      `json:"rootPath,omitempty"`
	CreateNewEnvironmentsFromPullRequests bool                        `json:"createNewEnvironmentsFromPullRequests,omitempty"`
	DiscoveryFileConfiguration            *DiscoveryFileConfiguration `json:"discoveryFileConfiguration,omitempty"`
}

func (payload *EnvironmentDiscoveryPutPayload) Invalidate() error {
	if payload.GithubInstallationId != 0 && payload.VcsConnectionId != "" {
		return errors.New("github_installation_id and vcs_connection_id are mutually exclusive")
	}

	return nil
}

type EnvironmentDiscoveryPayload struct {
	Id                                    string                      `json:"id"`
	GlobPattern                           string                      `json:"globPattern"                           tfschema:",omitempty"`
	EnvironmentPlacement                  string                      `json:"environmentPlacement"`
	WorkspaceNaming                       string                      `json:"workspaceNaming"`
	AutoDeployByCustomGlob                string                      `json:"autoDeployByCustomGlob"                tfschema:",omitempty"`
	Repository                            string                      `json:"repository"                            tfschema:",omitempty"`
	TerraformVersion                      string                      `json:"terraformVersion"                      tfschema:",omitempty"`
	OpentofuVersion                       string                      `json:"opentofuVersion"                       tfschema:",omitempty"`
	TerragruntVersion                     string                      `json:"terragruntVersion"                     tfschema:",omitempty"`
	TerragruntTfBinary                    string                      `json:"terragruntTfBinary"                    tfschema:",omitempty"`
	IsTerragruntRunAll                    bool                        `json:"isTerragruntRunAll"`
	Type                                  string                      `json:"type"`
	TokenId                               string                      `json:"tokenId"                               tfschema:",omitempty"`
	SshKeys                               []TemplateSshKey            `json:"sshKeys"                               tfschema:"-"`
	GithubInstallationId                  int                         `json:"githubInstallationId"                  tfschema:",omitempty"`
	VcsConnectionId                       string                      `json:"vcsConnectionId"                       tfschema:",omitempty"`
	BitbucketClientKey                    string                      `json:"bitbucketClientKey"                    tfschema:",omitempty"`
	IsAzureDevops                         bool                        `json:"isAzureDevOps"`
	IsBitbucketServer                     bool                        `json:"isBitbucketServer"`
	IsGitHubEnterprise                    bool                        `json:"isGitHubEnterprise"                    tfschema:"is_github_enterprise"`
	IsGitLabEnterprise                    bool                        `json:"isGitLabEnterprise"                    tfschema:"is_gitlab_enterprise"`
	Retry                                 TemplateRetry               `json:"retry"                                 tfschema:"-"`
	RootPath                              string                      `json:"rootPath"                              tfschema:",omitempty"`
	CreateNewEnvironmentsFromPullRequests bool                        `json:"createNewEnvironmentsFromPullRequests"`
	DiscoveryFileConfiguration            *DiscoveryFileConfiguration `json:"discoveryFileConfiguration"`
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
