provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name        = "Test-Project-Custom-Flow-${random_string.random.result}"
  description = "Test Description"
}

data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_approval_policy" "test" {
  name                   = "approval-policy-PROJECT-${env0_project.project.id}"
  repository             = data.env0_template.github_template.repository
  github_installation_id = data.env0_template.github_template.github_installation_id
}

resource "env0_approval_policy_assignment" "assignment" {
  scope_id     = env0_project.project.id
  blueprint_id = env0_approval_policy.test.id
}
