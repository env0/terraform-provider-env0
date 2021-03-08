data "env0_configuration_variable" "region" {
  name = "AWS_DEFAULT_REGION"
}

output "region_value" {
  value = data.env0_configuration_variable.region.value
}
output "region_id" {
  value = data.env0_configuration_variable.region.id
}
