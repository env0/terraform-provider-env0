resource "env0_environment_output_configuration_variable" "example" {
  name                  = "name"
  scope                 = "PROJECT"
  scope_id              = "project_to_assign_to_id"
  is_read_only          = true
  is_required           = true
  type                  = "terraform"
  description           = "description"
  output_environment_id = "output_environment_id"
  output_name           = "my_output_name"
}
