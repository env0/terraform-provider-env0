provider "random" {}

resource "random_string" "random" {
  length = 5
  special = false
  min_lower = 5
}

resource "env0_project" "test_project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}
data "env0_project" "data_by_name" {
  name = env0_project.test_project.name
}

data "env0_project" "data_by_id" {
  id = env0_project.test_project.id
}

output "test_project_name" {
  value = replace(env0_project.test_project.name, random_string.random.result, "")
}

output "test_project_description" {
  value = env0_project.test_project.description
}
