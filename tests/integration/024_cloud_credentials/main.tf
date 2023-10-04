resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_aws_credentials" "aws_cred1" {
  name = "Test Role arn1 ${random_string.random.result}"
  arn  = "Role ARN1"
}

resource "env0_aws_credentials" "aws_cred2" {
  name = "Test Role arn2 ${random_string.random.result}"
  arn  = "Role ARN2"
}

resource "env0_gcp_credentials" "gcp_cred" {
  name                = "name ${random_string.random.result}"
  service_account_key = "your_account_key"
  project_id          = "your_project_id"
}

data "env0_cloud_credentials" "all_aws_credentials" {
  depends_on = [env0_aws_credentials.aws_cred1, aws_cred2,env0_gcp_credentials.gcp_cred]
  credential_type = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
}

data "env0_aws_credentials" "aws_credentials1" {
  name     = data.env0_cloud_credentials.all_aws_credentials.names[index(data.env0_cloud_credentials.all_aws_credentials.names, env0_aws_credentials.aws_cred1.name)]
}

data "env0_aws_credentials" "aws_credentials2" {
  name     = data.env0_cloud_credentials.all_aws_credentials.names[index(data.env0_cloud_credentials.all_aws_credentials.names, env0_aws_credentials.aws_cred2.name)]
}

output "credentials_name" {
  value = var.second_run ? replace(data.env0_aws_credentials.aws_credentials1.name, random_string.random.result, "") : ""
}
