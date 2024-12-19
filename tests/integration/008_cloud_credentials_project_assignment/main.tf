resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_cloud_credentials_project_assignment" "example" {
  credential_id = env0_aws_credentials.credentials.id
  project_id    = env0_project.test_project.id
}

resource "env0_project" "test_project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "Test Description"
}

resource "env0_aws_credentials" "credentials" {
  name = "example-${random_string.random.result}"
  arn  = "Example role ARN"
}

data "env0_project_cloud_credentials" "project_cloud_credentials" {
  project_id = env0_project.test_project.id
}

output "validate" {
  value = var.second_run ? "${data.env0_project_cloud_credentials.project_cloud_credentials.ids.0}" == "${env0_aws_credentials.credentials.id}" ? "1" : "0" : "0"
}
