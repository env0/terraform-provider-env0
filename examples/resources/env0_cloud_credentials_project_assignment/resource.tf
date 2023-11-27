resource "env0_aws_credentials" "credentials" {
  name = "example"
  arn  = "Example role ARN"
}

resource "env0_aws_oidc_credentials" "credentials" {
  name     = "example"
  role_arn = "Example role ARN"
}

data "env0_project" "project" {
  name = "Default Organization Project"
}

resource "env0_cloud_credentials_project_assignment" "example" {
  credential_id = env0_aws_credentials.credentials.id
  project_id    = data.env0_project.project.id
}

resource "env0_cloud_credentials_project_assignment" "example_oidc" {
  credential_id = env0_aws_oidc_credentials.credentials.id
  project_id    = data.env0_project.project.id
}
