resource "env0_aws_oidc_credentials" "credentials" {
  name     = "example"
  role_arn = "arn::role::34"
  duration = 3600
}

resource "env0_gcp_oidc_credentials" "credentials" {
  name = "example"
  credential_configuration_file_content = jsonencode({
    "key" : "value"
  })
}
