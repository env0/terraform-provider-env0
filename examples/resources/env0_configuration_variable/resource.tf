resource "env0_configuration_variable" "example" {
  name        = "ENVIRONMENT_VARIABLE_NAME"
  value       = "example value"
  description = "Here you can fill description for this variable, note this field have limit of 255 chars"
}

resource "env0_configuration_variable" "drop_down" {
  name  = "ENVIRONMENT_VARIABLE_DROP_DOWN"
  value = "first option"
  enum = [
    "first option",
    "second option"
  ]
}

resource "env0_configuration_variable" "json_variable" {
  name   = "organization_tf_json_var"
  type   = "terraform"
  value  = "{ \"a\": 1234 }"
  format = "JSON"
}

# Workflow scope: set template_id to the workflow template and omit sub_environment_alias.
# The variable is applied to every sub environment of the workflow.
resource "env0_configuration_variable" "workflow_example" {
  name        = "WORKFLOW_VARIABLE"
  value       = "workflow value"
  template_id = "example-workflow-template-id"
  description = "Variable for the entire workflow scope"
}

# Sub environment scope: set both template_id and sub_environment_alias.
# The variable is applied only to the matching sub environment.
resource "env0_configuration_variable" "sub_environment_example" {
  name                  = "SUB_ENVIRONMENT_VARIABLE"
  value                 = "sub env value"
  template_id           = "example-template-id"
  sub_environment_alias = "example-alias"
  description           = "Variable for a sub environment scope"
}
