resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}
data "env0_project" "data_by_name" {
  name = env0_project.test_project.name
}

data "env0_project" "data_by_id" {
  id = env0_project.test_project.id
}

output "test_project_name" {
  value = env0_project.test_project.name
}

output "test_project_description" {
  value = env0_project.test_project.description
}
