data "env0_template" "template" {
  name = "Template Name"
}

data "env0_source_code_variables" "variables" {
  template_id = data.env0_template.template.id
}

output "variable_0_value" {
  value = data.env0_source_code_variables.variables.0.value
}
