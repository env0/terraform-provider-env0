data "env0_cloud_credentials" "aws_credentials" {
  credential_type = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
}

data "env0_aws_credentials" "aws_credentials" {
  for_each = toset(data.env0_cloud_credentials.aws_credentials.names)
  name     = each.value
}

output "credentials_name" {
  value = var.second_run ? data.env0_aws_credentials.aws_credentials["Test Role arn1"].name : ""
}
