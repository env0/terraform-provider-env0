resource "env0_aws_credentials" "credentials" {
  name        = "example"
  arn         = "Example role ARN"
  external_id = "Example external id"
}

data "env0_project" "project" {
  name = "Default Organization Project"
}

resource "env0_cloud_credentials_project_assignment" "example" {
  credential_id = env0_aws_credentials.credentials.id
  project_id = data.env0_project.project.id
}