# Audit logs only: the harness runs tests in parallel in one organization, and deployment log
# forwarding would send the other tests' deployment logs to this fake bucket.
resource "env0_log_forwarding" "test_s3" {
  type = "S3"
  value = jsonencode({
    bucketName = "tf-provider-integration-tests-042"
    awsRegion  = "us-east-1"
    bucketPath = var.second_run ? "second-run" : "first-run"
    roleArn    = "arn:aws:iam::123456789012:role/tf-provider-integration-tests-042"
  })
  audit_log_forwarding = true
}

output "log_forwarding_type" {
  value = env0_log_forwarding.test_s3.type
}

output "audit_log_forwarding" {
  value = tostring(env0_log_forwarding.test_s3.audit_log_forwarding)
}

output "deployment_log_forwarding" {
  value = tostring(env0_log_forwarding.test_s3.deployment_log_forwarding)
}
