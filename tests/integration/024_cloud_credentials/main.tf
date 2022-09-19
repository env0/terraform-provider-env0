resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_aws_credentials" "cred1" {
  name        = "Test Role arn1 ${random_string.random.result}"
  arn         = "Role ARN1"
  external_id = "External id1"
}

resource "env0_gcp_credentials" "cred2" {
  name                = "name ${random_string.random.result}"
  service_account_key = "your_account_key"
  project_id          = "your_project_id"
}

data "env0_cloud_credentials" "aws_credentials" {
  credential_type = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
}

data "env0_aws_credentials" "aws_credentials" {
  for_each = toset(data.env0_cloud_credentials.aws_credentials.names)
  name     = each.value
}

output "credentials_name" {
  value = var.second_run ? replace(data.env0_aws_credentials.aws_credentials["Test Role arn1 ${random_string.random.result}"].name, random_string.random.result, "") : ""
}
