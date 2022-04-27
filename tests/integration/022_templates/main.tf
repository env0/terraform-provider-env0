data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_template" "github_template1" {
  name                                    = "Github Test-111"
  description                             = "Template description - GitHub"
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
  name                                    = "Github Test-222"
  description                             = "Template description - GitHub"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

data "env0_templates" "all_templates" {}

data "env0_template" "templates" {
  for_each = toset(data.env0_templates.all_templates.names)
  name     = each.value
}

output "template1_name" {
  value = var.second_run ? data.env0_template.templates["Github Test-111"].name : ""
}

output "template2_name" {
  value = var.second_run ? data.env0_template.templates["Github Test-222"].name : ""
}
