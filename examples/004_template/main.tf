data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_template" "tested1" {
  name                                    = "tested1"
  description                             = "Tested 1 description"
  type                                    = "terraform"
  repository                              = "https://github.com/shlomimatichin/env0-template-jupyter-gpu"
  path                                    = var.second_run ? "second" : ""
  project_ids                             = [data.env0_project.default_project.id]
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
}

resource "env0_configuration_variable" "in_a_template" {
  name        = "fake_key"
  value       = "fake value"
  template_id = env0_template.tested1.id
}

resource "env0_configuration_variable" "in_a_template2" {
  name        = "fake_key_2"
  value       = "fake value 2"
  template_id = env0_template.tested1.id
  type        = "terraform"
}

data "env0_template" "tested2" {
  depends_on = [env0_template.tested1]
  name       = "tested1"
}
output "tested2_template_id" {
  value = data.env0_template.tested2.id
}
output "tested2_template_type" {
  value = data.env0_template.tested2.type
}
output "tested2_template_name" {
  value = data.env0_template.tested2.name
}
output "tested2_template_repository" {
  value = data.env0_template.tested2.repository
}
output "tested2_template_path" {
  value = data.env0_template.tested2.path
}

data "env0_template" "tested3" {
  id = env0_template.tested1.id
}
