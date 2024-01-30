provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}

resource "env0_project_budget" "test_project_budget" {
  project_id = env0_project.test_project.id
  amount     = var.second_run ? 10 : 20
  timeframe  = var.second_run ? "WEEKLY" : "MONTHLY"
}

resource "env0_project" "test_project2" {
  name        = "Test-Project-${random_string.random.result}2"
  description = "Test Description2 ${var.second_run ? "after update" : ""}"
}

resource "env0_project" "test_sub_project" {
  name              = "Test-Sub-Project-${random_string.random.result}"
  description       = "Test Description ${var.second_run ? "after update" : ""}"
  parent_project_id = var.second_run ? env0_project.test_project2.id : env0_project.test_project.id
}

resource "env0_project" "test_sub_project_to_project" {
  name              = "Test-Sub-Project-To-Project-${random_string.random.result}"
  description       = "Test Description ${var.second_run ? "after update" : ""}"
  parent_project_id = var.second_run ? "" : env0_project.test_project.id
}


data "env0_project" "data_by_name" {
  name = env0_project.test_project.name
}

data "env0_project" "data_by_id" {
  id = env0_project.test_project.id
}

resource "env0_project" "test_project_other" {
  name        = "Test-Project-${random_string.random.result}-other"
  description = "Test Description"
}

resource "env0_project" "test_sub_project_other" {
  name              = "Test-Sub-Project-${random_string.random.result}"
  description       = "Test Description"
  parent_project_id = env0_project.test_project_other.id
}

data "env0_project" "data_by_name_with_parent_name" {
  name                = env0_project.test_sub_project_other.name
  parent_project_name = env0_project.test_project_other.name
}

data "env0_project" "data_by_name_with_parent_id" {
  name              = env0_project.test_sub_project_other.name
  parent_project_id = env0_project.test_project_other.id
}

data "env0_projects" "list_of_projects" {}

output "test_project_name" {
  value = replace(env0_project.test_project.name, random_string.random.result, "")
}

output "test_project_description" {
  value = env0_project.test_project.description
}
