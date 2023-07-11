provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_template" "github_template1" {
  name                                    = "Github Test Templates ${random_string.random.result}-1"
  description                             = "template-1"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "github_template2" {
  name                                    = "Github Test Templates ${random_string.random.result}-2"
  description                             = "template-2"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

data "env0_templates" "all_templates" {
  depends_on = [env0_template.github_template1, env0_template.github_template2]
}

data "env0_template" "github_template1" {
  name = data.env0_templates.all_templates.names[index(data.env0_templates.all_templates.names, env0_template.github_template1.name)]
}

data "env0_template" "github_template2" {
  name = data.env0_templates.all_templates.names[index(data.env0_templates.all_templates.names, env0_template.github_template2.name)]
}

output "template1_description" {
  value = var.second_run ? data.env0_template.github_template1.description : ""
}

output "template2_description" {
  value = var.second_run ? data.env0_template.github_template2.description : ""
}
