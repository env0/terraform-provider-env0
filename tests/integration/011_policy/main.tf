data "env0_project" "default" {
  name = "Default Organization Project"
}

resource "env0_policy" "test_policy" {
  project_id                   = data.env0_project.default.id
  number_of_environments_total = 1
}

resource "env0_policy" "test_policy_2" {
  project_id                   = data.env0_project.default.id
  number_of_environments_total = 0
}
