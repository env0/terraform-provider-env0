data "env0_configuration_variable" "aws_default_region" {
  name = "AWS_DEFAULT_REGION"
}

output "aws_default_region" {
  value = data.env0_configuration_variable.aws_default_region.value
}
