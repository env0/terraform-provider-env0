resource "env0_project" "test_project" {
  name          = "Test-Project-for-environment"
  force_destroy = true
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
  wait_for      = "FULLY_DEPLOYED"
  configuration {
    name  = "environment configuration variable"
    value = "value"
  }
  approve_plan_automatically = true
}

data "env0_configuration_variable" "env_config_variable" {
  environment_id = env0_environment.example.id
  name           = "environment configuration variable"
}

resource "env0_template" "terragrunt_template" {
  name               = "Terragrunt template for environment resource"
  type               = "terragrunt"
  repository         = "https://github.com/env0/templates"
  path               = "misc/null-resource"
  terraform_version  = "0.15.1"
  terragrunt_version = "0.35.0"
}

resource "env0_environment" "terragrunt_environment" {
  force_destroy                    = true
  name                             = "environment"
  project_id                       = env0_project.test_project.id
  template_id                      = env0_template.terragrunt_template.id
  approve_plan_automatically       = true
  revision                         = "master"
  terragrunt_working_directory     = var.second_run ? "/second-dir" : "/first-dir"
  auto_deploy_on_path_changes_only = false
}

data "env0_environment" "test" {
  id = env0_environment.example.id
}

output "name" {
  value = data.env0_environment.test.name
}

output "configurationVariable" {
  value = data.env0_configuration_variable.env_config_variable.name
}

output "terragrunt_working_directory" {
  value = env0_environment.terragrunt_environment.terragrunt_working_directory
}
