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

resource "env0_custom_flow" "test" {
  name                   = "Custom Flow Github Test ${random_string.random.result}"
  repository             = data.env0_template.github_template.repository
  github_installation_id = data.env0_template.github_template.github_installation_id
  path                   = "custom-flows/opa.yaml"
}

data "env0_custom_flow" "test" {
  name       = env0_custom_flow.test.name
  depends_on = [env0_custom_flow.test]
}

resource "env0_custom_flow_assignment" "assignment" {
  scope_id    = env0_project.project.id
  template_id = data.env0_custom_flow.test.id
}
