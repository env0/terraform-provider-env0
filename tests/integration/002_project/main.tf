data "env0_project" "default" {
  name = "Default Organization Project"
}

resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}

data "env0_project" "default2" {
  depends_on = [data.env0_project.default]
  id         = data.env0_project.default.id
}

output "default_project_id" {
  value = data.env0_project.default.id
}

output "default_project_name" {
  value = data.env0_project.default2.name
}

output "default_description" {
  value = env0_project.test_project.description
}
