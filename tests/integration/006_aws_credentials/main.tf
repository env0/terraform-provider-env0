resource "env0_aws_credentials" "my_role" {
  name        = "Test Role"
  arn         = "Role ARN"
  external_id = "External id"
}

data "env0_aws_credentials" "my_role" {
  name       = "Test Role"
  depends_on = [env0_aws_credentials.my_role]
}

output "name" {
  value = data.env0_aws_credentials.my_role.name
}