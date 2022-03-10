resource "env0_aws_credentials" "my_role_by_arn" {
  name        = "Test Role arn"
  arn         = "Role ARN"
  external_id = "External id"
}

data "env0_aws_credentials" "my_role_by_arn" {
  name       = "Test Role"
  depends_on = [env0_aws_credentials.my_role_by_arn]
}

resource "env0_aws_credentials" "my_role_by_access_key" {
  name              = "Test Role access key"
  access_key_id     = "Access id"
  secret_access_key = "Secret Access id"
}

data "env0_aws_credentials" "my_role_by_access_key" {
  name       = "Test Role access key"
  depends_on = [env0_aws_credentials.my_role_by_access_key]
}


output "name_by_arn" {
  value = data.env0_aws_credentials.my_role_by_arn.name
}

output "name_access_key" {
  value = data.env0_aws_credentials.my_role_by_access_key.name
}