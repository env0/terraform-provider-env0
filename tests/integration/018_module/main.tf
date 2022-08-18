provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_module" "test_module" {
  module_name            = "test-module-${random_string.random.result}"
  module_provider        = "testprovider1"
  token_id               = var.second_run ? null : "37689e5a-5298-4555-b71e-92b80f736222"
  token_name             = var.second_run ? null : "johns_token"
  repository             = "https://gitlab.com/moooo/moooo-docs-aws-functions.git"
  github_installation_id = var.second_run ? 32112 : null
}
