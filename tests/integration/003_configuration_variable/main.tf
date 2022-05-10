resource "env0_project" "test_project" {
  name = "Project-for-003-configuration-variable"
}

resource "env0_configuration_variable" "region_in_project_resource" {
  name       = "AWS_DEFAULT_REGION"
  project_id = env0_project.test_project.id
  value      = "il-tnuvot-1"
}

data "env0_configuration_variable" "region_in_project" {
  name       = "AWS_DEFAULT_REGION"
  project_id = env0_project.test_project.id
  depends_on = [env0_configuration_variable.region_in_project_resource]
}
output "region_in_project_value" {
  value     = data.env0_configuration_variable.region_in_project.value
  sensitive = true
}
output "region_in_project_id" {
  value = data.env0_configuration_variable.region_in_project.id
}

resource "env0_configuration_variable" "tested1" {
  name  = "tested1"
  value = "fake value 1 ${var.second_run ? "after update" : ""}"
}
data "env0_configuration_variable" "tested1" {
  name       = "tested1"
  depends_on = [env0_configuration_variable.tested1]
}

output "tested1_value" {
  value     = data.env0_configuration_variable.tested1.value
  sensitive = true
}

data "env0_configuration_variable" "tested2" {
  id = env0_configuration_variable.tested1.id
}

resource "env0_configuration_variable" "tested3" {
  name  = "tested3"
  value = "First"
  enum  = ["First", "Second"]
}
data "env0_configuration_variable" "tested3" {
  name       = "tested3"
  depends_on = [env0_configuration_variable.tested3]
}


output "tested3_enum_1" {
  value     = data.env0_configuration_variable.tested3.enum[0]
  sensitive = true
}
output "tested3_enum_2" {
  value     = data.env0_configuration_variable.tested3.enum[1]
  sensitive = true
}

resource "env0_configuration_variable" "regex_var" {
  project_id = env0_project.test_project.id
  name       = "regex_var"
  regex      = "^test-\\d+$"
}

data "env0_configuration_variable" "regex_var" {
  project_id = env0_project.test_project.id
  name       = "regex_var"
  depends_on = [env0_configuration_variable.regex_var]
}

output "regex" {
  value = data.env0_configuration_variable.regex_var.regex
}
