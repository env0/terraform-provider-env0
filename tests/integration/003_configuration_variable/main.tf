resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name = "Project-for-003-configuration-variable-${random_string.random.result}"
}

resource "env0_configuration_variable" "region_in_project_resource" {
  name       = "AWS_DEFAULT_REGION"
  project_id = env0_project.test_project.id
  value      = "il-tnuvot-1"
}

resource "env0_configuration_variable" "tested1" {
  name  = "tested1-${random_string.random.result}"
  value = "fake value 1 ${var.second_run ? "after update" : ""}"
}

resource "env0_configuration_variable" "tested3" {
  name  = "tested3-${random_string.random.result}"
  value = "First"
  enum  = ["First", "Second"]
}

resource "env0_configuration_variable" "regex_var" {
  project_id = env0_project.test_project.id
  name       = "regex_var"
  regex      = "^test-\\d+$"
}

resource "time_sleep" "wait_5_seconds" {
  create_duration = "5s"

  depends_on = [env0_configuration_variable.region_in_project_resource, env0_configuration_variable.tested1, env0_configuration_variable.tested3, env0_configuration_variable.regex_var]
}

data "env0_configuration_variable" "region_in_project" {
  name       = "AWS_DEFAULT_REGION"
  project_id = env0_project.test_project.id

  depends_on = [time_sleep.wait_5_seconds]
}
output "region_in_project_value" {
  value     = data.env0_configuration_variable.region_in_project.value
  sensitive = true
}
output "region_in_project_id" {
  value = data.env0_configuration_variable.region_in_project.id
}


data "env0_configuration_variable" "tested1" {
  name = "tested1-${random_string.random.result}"

  depends_on = [time_sleep.wait_5_seconds, env0_configuration_variable.tested1]
}

output "tested1_value" {
  value     = replace(data.env0_configuration_variable.tested1.value, random_string.random.result, "")
  sensitive = true
}

data "env0_configuration_variable" "tested2" {
  id = env0_configuration_variable.tested1.id

  depends_on = [time_sleep.wait_5_seconds]
}


data "env0_configuration_variable" "tested3" {
  name = "tested3-${random_string.random.result}"

  depends_on = [time_sleep.wait_5_seconds]
}

output "tested3_enum_1" {
  value     = data.env0_configuration_variable.tested3.enum[0]
  sensitive = true
}
output "tested3_enum_2" {
  value     = data.env0_configuration_variable.tested3.enum[1]
  sensitive = true
}

data "env0_configuration_variable" "regex_var" {
  project_id = env0_project.test_project.id
  name       = "regex_var"

  depends_on = [time_sleep.wait_5_seconds]
}

output "regex" {
  value = data.env0_configuration_variable.regex_var.regex
}
