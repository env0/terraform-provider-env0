resource "env0_template_project_assignment" "template_project_test" {
  template_id = env0_template.test_template.id
  project_id  = env0_project.test_project.id
}

resource "env0_template" "test_template" {
  name        = "Test-Template"
  description = "test template"
  type        = "terraform"
  repository  = "https://github.com/env0/templates"
}

resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description"
}
