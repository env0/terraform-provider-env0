data "env0_configuration_variable" "region" {
  name = "AWS_DEFAULT_REGION"
}
output "region_value" {
  value = data.env0_configuration_variable.region.value
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
  value = data.env0_configuration_variable.region_in_project.value
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
  value = data.env0_configuration_variable.tested1.value
}
