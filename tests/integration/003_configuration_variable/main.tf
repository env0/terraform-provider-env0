data "env0_configuration_variable" "region" {
  name = "AWS_DEFAULT_REGION"
}
output "region_value" {
  value     = data.env0_configuration_variable.region.value
  sensitive = true
}
output "region_id" {
  value = data.env0_configuration_variable.region.id
}


data "env0_project" "default" {
  name = "Default Organization Project"
}
data "env0_configuration_variable" "region_in_project" {
  name       = "AWS_DEFAULT_REGION"
  project_id = data.env0_project.default.id
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