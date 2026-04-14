# GitHub - organization
data "env0_vcs_connection" "github_org" {
  access_scope    = "Organization:my-org"
  connection_type = "DeploymentPipeline"
}

# GitHub - user
data "env0_vcs_connection" "github_user" {
  access_scope    = "User:my-username"
  connection_type = "DeploymentPipeline"
}

# Bitbucket - workspace
data "env0_vcs_connection" "bitbucket" {
  access_scope    = "Workspace:my-workspace"
  connection_type = "DeploymentPipeline"
}

# GitLab - token
data "env0_vcs_connection" "gitlab" {
  access_scope    = "Token:my-token-name"
  connection_type = "DeploymentPipeline"
}

# Azure DevOps
data "env0_vcs_connection" "azure" {
  access_scope    = "User:my-display-name"
  connection_type = "DeploymentPipeline"
}

# Self-hosted (e.g. GitHub Enterprise)
data "env0_vcs_connection" "github_enterprise" {
  access_scope    = "url:https://github.mycompany.com"
  connection_type = "DeploymentPipeline"
}

# CodeWrite connection type
data "env0_vcs_connection" "github_codewrite" {
  access_scope    = "Organization:my-org"
  connection_type = "CodeWrite"
}

# Usage with a template
data "env0_vcs_connection" "github" {
  access_scope    = "Organization:my-org"
  connection_type = "DeploymentPipeline"
}

resource "env0_template" "example" {
  name              = "example"
  repository        = "https://github.com/my-org/my-repo"
  type              = "terraform"
  vcs_connection_id = data.env0_vcs_connection.github.id
}
