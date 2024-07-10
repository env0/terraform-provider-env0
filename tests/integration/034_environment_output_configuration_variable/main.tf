provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name = "project-environment-output1-${random_string.random.result}"
}

resource "env0_environment_output_configuration_variable" "to_project" {
  name                  = "name1-${random_string.random.result}"
  scope                 = "PROJECT"
  scope_id              = env0_project.project.id
  type                  = "terraform"
  description           = "description"
  output_environment_id = "d73c4f0c-5569-44c6-936f-25a042b1aedf"
  output_name           = "not_sensitive_output"
}
