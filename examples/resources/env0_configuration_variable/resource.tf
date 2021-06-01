resource "env0_configuration_variable" "example" {
  name  = "ENVIRONMENT_VARIABLE_NAME"
  value = "example value anoter"
  project_id = "51483b68-189b-428e-ab63-3476c53f0050"
  is_sensitive = true
  type = "terraform"
}

data "env0_configuration_variable" "example" {
  id = "8113c2db-1454-4732-b770-ea7b7e6e6408"
}

output "some" {
  value = data.env0_configuration_variable.example.value
}

data "env0_configuration_variable" "exampleGlobal" {
  id = "1937a457-9930-4591-bb47-1d0aa3a60509"
}

output "someGlobal" {
  value = data.env0_configuration_variable.exampleGlobal.value
}