data "env0_project" "default_project" {
  name = "Default Organization Project"
}

data "env0_project" "with_parent_name_filter" {
  name                = "Default Organization Project"
  parent_project_name = "parent projet name"
}

data "env0_project" "with_parent_id_filter" {
  name              = "Default Organization Project"
  parent_project_id = "parent-projet-id"
}
