provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

# Github Integration must be done manually - so we expect an existing Github Template with this name -
# It must be for https://github.com/env0/templates - We validate that in the outputs
data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

# Gitlab Integration must be done manually - so we expect an existing Gitlab Template with this name
# It must be for https://gitlab.com/env0/gitlab-vcs-integration-tests - the gitlab_project_id is still static
data "env0_template" "gitlab_template" {
  name = "Gitlab Integrated Template"
}

resource "env0_template" "github_template" {
  name                                    = "Github Test-${random_string.random.result}"
  description                             = "Template description - GitHub"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = var.second_run ? "/second" : "/misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "gitlab_template" {
  name                                    = "GitLab Test-${random_string.random.result}"
  description                             = "Template description - Gitlab"
  type                                    = "terraform"
  repository                              = data.env0_template.gitlab_template.repository
  token_id                                = data.env0_template.gitlab_template.token_id
  gitlab_project_id                       = 32315446
  path                                    = var.second_run ? "second" : "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "template_tg" {
  name               = "Template for environment resource - tg-${random_string.random.result}"
  type               = "terragrunt"
  repository         = "https://github.com/env0/templates"
  path               = "terragrunt/misc/null-resource"
  terraform_version  = "0.15.1"
  terragrunt_version = "0.35.0"
}

resource "env0_configuration_variable" "in_a_template" {
  name        = "fake_key"
  value       = "fake value"
  template_id = env0_template.github_template.id
}

resource "env0_configuration_variable" "in_a_template2" {
  name        = "fake_key_2"
  value       = "fake value 2"
  template_id = env0_template.github_template.id
  type        = "terraform"
}

resource "env0_template" "github_template_source_code" {
  name                                    = "Github Test Source Code-${random_string.random.result}"
  description                             = "Template description - GitHub"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = "misc/custom-flow-tf-vars"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "helm_template" {
  name        = "helm-${random_string.random.result}-1"
  description = "Template description helm"
  repository  = "https://github.com/env0/templates"
  path        = "misc/helm/dummy"
  type        = "helm"
}

resource "env0_template" "helm_template_repo" {
  name               = "helm-${random_string.random.result}-2"
  description        = "Template description helm repo"
  repository         = "https://charts.bitnami.com/bitnami"
  type               = "helm"
  helm_chart_name    = "nginx"
  is_helm_repository = true
}

data "env0_source_code_variables" "variables" {
  template_id = env0_template.github_template_source_code.id
}

output "github_variables_name" {
  value = data.env0_source_code_variables.variables.variables.0.name
}

output "github_variables_value" {
  value = data.env0_source_code_variables.variables.variables.0.value
}

output "github_template_id" {
  value = env0_template.github_template.id
}
output "github_template_type" {
  value = env0_template.github_template.type
}
output "github_template_name" {
  value = replace(env0_template.github_template.name, random_string.random.result, "")
}
output "github_template_repository" {
  value = env0_template.github_template.repository
}
output "gitlab_template_repository" {
  value = env0_template.gitlab_template.repository
}
output "github_template_path" {
  value = env0_template.github_template.path
}
output "tg_tg_version" {
  value = env0_template.template_tg.terragrunt_version
}

output "data_github_template_type" {
  value = data.env0_template.github_template.type
}
