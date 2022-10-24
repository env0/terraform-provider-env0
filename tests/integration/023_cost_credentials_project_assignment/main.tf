provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "description"
}

resource "env0_aws_cost_credentials" "cost" {
  name        = "cost-${random_string.random.result}"
  arn         = "arn"
  external_id = "external-id"
}

resource "env0_cost_credentials_project_assignment" "cost_project_assignment" {
  credential_id = env0_aws_cost_credentials.cost.id
  project_id    = env0_project.project.id
}
