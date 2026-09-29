# Create a GitHub Enterprise VCS connection with default agent
resource "env0_vcs_connection" "github_enterprise" {
  name          = "github-enterprise"
  type          = "GitHubEnterprise"
  url           = "https://github.example.com"
  vcs_agent_key = "ENV0_DEFAULT"
}

# Create a GitLab Enterprise VCS connection with custom agent
resource "env0_vcs_connection" "gitlab_enterprise" {
  name          = "gitlab-enterprise"
  type          = "GitLabEnterprise"
  url           = "https://gitlab.example.com"
  vcs_agent_key = "my-custom-agent"
}

# Create a BitBucket Server VCS connection
resource "env0_vcs_connection" "bitbucket_server" {
  name = "bitbucket-server"
  type = "BitBucketServer"
  url  = "https://bitbucket.example.com"
}

# Create an Azure DevOps VCS connection from a token authorized in the env0 UI
resource "env0_vcs_connection" "azure_devops" {
  name     = "azure-devops"
  type     = "AzureDevOps"
  token_id = "11111111-1111-1111-1111-111111111111"
}
