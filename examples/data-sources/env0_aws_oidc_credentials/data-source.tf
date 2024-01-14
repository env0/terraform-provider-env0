resource "env0_aws_oidc_credentials" "example" {
  name     = "name"
  role_arn = "role_arn"
}

data "env0_aws_oidc_credentials" "by_id" {
  id = env0_aws_oidc_credentials.example.id
}

data "env0_aws_oidc_credentials" "by_name" {
  name = env0_aws_oidc_credentials.example.name
}
