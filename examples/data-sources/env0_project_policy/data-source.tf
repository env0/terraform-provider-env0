data "env0_project" "project" {
  name = "project"
}

data "env0_project_policy" "project_policy" {
  project_id = data.env0_project.project.id
}
