data "env0_user" "user_example" {
  email = "example@email.com"
}

resource "env0_project" "project_example" {
  name = "project-example"
}

resource "env0_user_project_assignment" "project_assignment_example" {
  user_id    = data.env0_user.user_example.id
  project_id = env0_project.project_example.id
  role       = "Viewer"
}
