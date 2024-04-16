resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_template_project_assignment" "template_project_test" {
  template_id = env0_template.test_template.id
  project_id  = env0_project.test_project.id
}

resource "env0_template" "test_template" {
  name        = "Test-Template-${random_string.random.result}"
  description = "test template"
  type        = "terraform"
  repository  = "https://github.com/env0/templates"
  terraform_verison = "1.7.1"
}

resource "env0_project" "test_project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "Test Description"
}
