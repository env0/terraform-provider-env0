data "env0_project" "project" {
  name = "my project name"
}

data "env0_variable_set" "variable_set_project_scope" {
  name       = "variable set name"
  scope      = "PROJECT"
  project_id = data.env0_project.project.id
}

data "env0_variable_set" "variable_set_organization_scope" {
  name  = "variable set name"
  scope = "ORGANIZATION"
}
