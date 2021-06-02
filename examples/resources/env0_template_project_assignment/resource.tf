data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_template" "example" {
  name        = "example"
  description = "Example template"
  repository  = "https://github.com/env0/templates"
  path        = "aws/hello-world"
}

resource "env0_template_project_assignment" "assignment" {
  template_id = env0_template.example.id
  project_id  = data.env0_project.default_project.id
}