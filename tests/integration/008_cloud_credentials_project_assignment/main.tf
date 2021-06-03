resource "env0_cloud_credentials_project_assignment" "example" {
  credential_id = env0_aws_credentials.credentials.id
  project_id    = env0_project.test_project.id
}

resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description"
}

resource "env0_aws_credentials" "credentials" {
  name        = "example"
  arn         = "Example role ARN"
  external_id = "Example external id"
}
