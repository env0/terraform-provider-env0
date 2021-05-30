/*data "env0_project" "test_project" {
  name = "Default Organization Project"
}*/

/*resource "env0_template_project_assignment" "template_project_test" {
  template_id = env0_template.test_template.id
  project_id  = env0_project.test_project.id
}*/

resource "env0_template" "test_template123" {
  name        = "Test-Template123"
  description = "test template"
  type        = "terraform"
  repository  = "https://github.com/shlomimatichin/env0-template-jupyter-gpu"
}

resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description"
}

data "env0_template" "tested2" {
  depends_on = [env0_template.test_template123]
  name       = "test_template123"
}

output "env0_project_id" {
  value = data.env0_template.tested2.name
}