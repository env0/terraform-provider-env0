resource "env0_aws_cloud_configuration" "example" {
  name        = "example"
  account_id  = "242345678901"
  bucket_name = "a_bucket_name"
  regions     = ["us-east-1", "us-west-2"]
}
