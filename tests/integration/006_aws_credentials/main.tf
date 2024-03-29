resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_aws_credentials" "my_role_by_arn" {
  name     = "Test Role arn ${random_string.random.result}"
  arn      = "Role ARN"
  duration = 7200
}

data "env0_aws_credentials" "my_role_by_arn" {
  name       = "Test Role arn ${random_string.random.result}"
  depends_on = [env0_aws_credentials.my_role_by_arn]
}

resource "env0_aws_credentials" "my_role_by_access_key" {
  name              = "Test Role access key ${random_string.random.result}"
  access_key_id     = "Access id"
  secret_access_key = var.second_run ? "Secret Access id2" : "secret1"
}

data "env0_aws_credentials" "my_role_by_access_key" {
  name       = "Test Role access key ${random_string.random.result}"
  depends_on = [env0_aws_credentials.my_role_by_access_key]
}

resource "env0_aws_oidc_credentials" "oidc_credentials" {
  name     = "Test Oidc Credentials ${random_string.random.result}"
  role_arn = var.second_run ? "Role ARN2" : "Role ARN1"
  duration = 7200
}

data "env0_aws_oidc_credentials" "oidc_credentials" {
  name       = "Test Oidc Credentials ${random_string.random.result}"
  depends_on = [env0_aws_oidc_credentials.oidc_credentials]
}

output "name_by_arn" {
  value = replace(data.env0_aws_credentials.my_role_by_arn.name, random_string.random.result, "")
}

output "name_access_key" {
  value = replace(data.env0_aws_credentials.my_role_by_access_key.name, random_string.random.result, "")
}
