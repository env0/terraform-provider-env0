data "env0_template" "example" {
  name = "Template Name"
}

data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_environment" "example" {
  name        = "environment"
  project_id  = data.env0_project.default_project.id
  template_id = data.env0_template.example.id
}

