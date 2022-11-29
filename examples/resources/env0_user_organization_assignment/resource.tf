data "env0_user" "user_example" {
  email = "example@email.com"
}

resource "env0_custom_role" "custom_role" {
  name = "custom-role"
  permissions = [
    "EDIT_PROJECT_SETTINGS"
  ]
}

resource "env0_user_organization_assignment" "user_org" {
  user_id        = data.env0_user.user_example.id
  custom_role_id = env0_custom_role.custom_role.id
}
