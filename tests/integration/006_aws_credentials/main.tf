resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_aws_credentials" "my_role_by_arn" {
  name        = "Test Role arn ${random_string.random.result}"
  arn         = "Role ARN"
  external_id = "External-id"
}

data "env0_aws_credentials" "my_role_by_arn" {
  name       = "Test Role arn ${random_string.random.result}"
  depends_on = [env0_aws_credentials.my_role_by_arn]
}

resource "env0_aws_credentials" "my_role_by_access_key" {
  name              = "Test Role access key ${random_string.random.result}"
  access_key_id     = "Access id"
  secret_access_key = "Secret Access id"
}

data "env0_aws_credentials" "my_role_by_access_key" {
  name       = "Test Role access key ${random_string.random.result}"
  depends_on = [env0_aws_credentials.my_role_by_access_key]
}


output "name_by_arn" {
  value = replace(data.env0_aws_credentials.my_role_by_arn.name, random_string.random.result, "")
}

output "name_access_key" {
  value = replace(data.env0_aws_credentials.my_role_by_access_key.name, random_string.random.result, "")
}
