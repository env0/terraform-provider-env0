resource "env0_aws_oidc_credentials" "credentials" {
  name     = "example"
  role_arn = "arn::role::34"
  duration = 3600
}
