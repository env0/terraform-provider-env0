data "env0_custom_roles" "all_roles" {}

data "env0_custom_role" "roles" {
  for_each = toset(data.env0_custom_roles.all_roles.names)
  name     = each.value
}
