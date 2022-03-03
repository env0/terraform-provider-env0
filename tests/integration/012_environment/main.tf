resource "env0_project" "test_project" {
  name = "Test-Project-for-environment"
}

resource "env0_template" "template" {
  name              = "Template for environment resource"
  type              = "terraform"
  repository        = "https://github.com/env0/templates"
  path              = "misc/null-resource"
  terraform_version = "0.15.1"
}

resource "env0_environment" "example" {
  force_destroy = true
  name          = "environment"
  project_id    = env0_project.test_project.id
  template_id   = env0_template.template.id
  configuration {
    name  = "environment configuration variable"
    value = "value"
  }
  approve_plan_automatically = true
}

data "env0_environment" "test" {
  id = env0_environment.example.id
}


data "env0_configuration_variable" "env_config_variable" {
  depends_on        = [data.env0_environment.test]
  deployment_log_id = data.env0_environment.test.deployment_id
  name              = "deployment configuration variable"
}

output "name" {
  value = data.env0_environment.test.name
}

output "configurationVariable" {
  value = data.env0_configuration_variable.env_config_variable.name
}

