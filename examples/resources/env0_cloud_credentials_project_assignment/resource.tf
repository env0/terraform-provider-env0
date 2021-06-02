data "env0_aws_credentials" "credentials" {
  name = "aws"
}

data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_cloud_credentials_project_assignment" "example" {
  credential_id = data.env0_aws_credentials.credentials.id
  project_id = data.env0_project.default_project.id
}
