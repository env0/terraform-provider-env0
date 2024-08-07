provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_aws_cloud_configuration" "aws_cloud_configuration" {
  name        = "aws-${random_string.random.result}"
  account_id  = var.second_run ? "012345678901" : "242345678901"
  bucket_name = var.second_run ? "my_bucket_name2" : "my_bucket_name1"
  regions     = var.second_run ? ["us-west-2"] : ["us-east-1", "us-west-2"]
}
