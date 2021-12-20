data "env0_project" "default" {
  name = "Default Organization Project"
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
  project_id    = data.env0_project.default.id
  template_id   = env0_template.template.id
  configuration {
    name  = "environment configuration variable"
    value = "value"
  }
}

data "env0_environment" "test" {
  id = env0_environment.example.id
}

resource "time_sleep" "wait_30_seconds" {
  depends_on      = [data.env0_environment.test]
  create_duration = "30s"
}


data "env0_configuration_variable" "env_config_variable" {
  depends_on     = [time_sleep.wait_30_seconds] // configuration variable scope update to "ENVIRONMENT" only after the deployment approved
  environment_id = data.env0_environment.test.id
  name           = "environment configuration variable"
}

output "name" {
  value = data.env0_environment.test.name
}

output "configurationVariable" {
  value = data.env0_configuration_variable.env_config_variable.name
}

