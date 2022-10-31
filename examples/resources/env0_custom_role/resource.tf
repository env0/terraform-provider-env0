
resource "env0_custom_role" "custom_role_example" {
  name = "my custom role 1"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS"
  ]
}
