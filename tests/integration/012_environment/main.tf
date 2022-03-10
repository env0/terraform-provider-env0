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

resource "env0_template" "template_tg" {
  name               = "Template for environment resource - tg"
  type               = "terragrunt"
  repository         = "https://github.com/env0/templates"
  path               = "terragrunt/misc/null-resource"
  terraform_version  = "0.15.1"
  terragrunt_version = "0.35.0"
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
  revision                   = "master"
}

resource "env0_environment" "example_tg" {
  force_destroy = true
  name          = "environment-tg"
  project_id    = env0_project.test_project.id
  template_id   = env0_template.template_tg.id
  configuration {
    name  = "environment configuration variable"
    value = "value"
  }
  approve_plan_automatically = true
  revision                   = "master"
}

data "env0_environment" "test" {
  id = env0_environment.example.id
}

data "env0_environment" "test_tg" {
  id = env0_environment.example_tg.id
}

output "revision" {
  value = data.env0_environment.test.revision
}

output "active_tg" {
  value = data.env0_environment.test_tg.status
}


