resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}

resource "env0_team" "team_resource" {
  name        = "Test-Team"
  description = var.second_run ? "second description" : "first description"
}

resource "env0_team_project_assignment" "assignment" {
  depends_on = [env0_team.team_resource, env0_project.test_project]
  project_id = env0_project.test_project.id
  team_id    = env0_team.team_resource.id
  role       = "Admin"
}
