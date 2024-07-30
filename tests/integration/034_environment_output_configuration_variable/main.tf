provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name = "project-environment-output-${random_string.random.result}"
}

data "env0_environment" "output_environment" {
  name = "Output Environment Integration"
}

resource "env0_environment_output_configuration_variable" "to_project" {
  name                  = "name1-${random_string.random.result}"
  scope                 = "PROJECT"
  scope_id              = env0_project.project.id
  type                  = "terraform"
  description           = "description"
  output_environment_id = data.env0_environment.output_environment.id
  output_name           = "not_sensitive_output"
}

data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_template" "template" {
  repository             = data.env0_template.github_template.repository
  github_installation_id = data.env0_template.github_template.github_installation_id
  name                   = "github-template--${random_string.random.result}"
  type                   = "terraform"
  path                   = "misc/null-resource"
  terraform_version      = "0.15.1"
}

resource "env0_template_project_assignment" "template_project_assignment" {
  template_id = env0_template.template.id
  project_id  = env0_project.project.id
}

resource "env0_environment" "environment" {
  depends_on = [env0_template_project_assignment.template_project_assignment]

  name          = "environment-${random_string.random.result}"
  project_id    = env0_project.project.id
  template_id   = env0_template.template.id
  force_destroy = true
}

resource "env0_environment_output_configuration_variable" "to_environment" {
  name                  = "name2-${random_string.random.result}"
  scope_id              = env0_environment.environment.id
  type                  = "terraform"
  description           = "description"
  output_environment_id = data.env0_environment.output_environment.id
  output_name           = "not_sensitive_output"
}
