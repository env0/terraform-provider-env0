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
}

data "env0_environment" "test" {
  id = env0_environment.example.id
}

output "name" {
  value = data.env0_environment.test.name
}

