resource "env0_configuration_variable" "example" {
  name  = "ENVIRONMENT_VARIABLE_NAME"
  value = "example value"
}

resource "env0_configuration_variable" "drop_down" {
  name  = "ENVIRONMENT_VARIABLE_DROP_DOWN"
  value = "first option"
  enum = [
    "first option",
    "second option"
  ]
}

