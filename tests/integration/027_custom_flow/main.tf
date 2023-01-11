provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
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
