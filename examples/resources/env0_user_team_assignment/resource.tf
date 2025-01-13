data "env0_user" "user_example" {
  email = "example@email.com"
}

resource "env0_team" "team_example" {
  name = "team-example"
}

resource "env0_user_team_assignment" "assignment_example" {
  user_id = data.env0_user.user_example.id
  team_id = env0_team.team_example.id
}
