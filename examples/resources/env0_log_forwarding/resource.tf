# Deployment logs to Datadog. The API key is a secret reference: `$${...}` escapes Terraform's own interpolation.
resource "env0_log_forwarding" "datadog" {
  type = "DATADOG"
  value = jsonencode({
    host   = "http-intake.logs.datadoghq.com"
    port   = 443
    apiKey = "$${ssm:datadog-api-key}"
  })
  deployment_log_forwarding = true
}

# Audit logs to an S3 bucket, written through an IAM role env0 assumes.
resource "env0_log_forwarding" "s3" {
  type = "S3"
  value = jsonencode({
    bucketName = "my-env0-audit-logs"
    awsRegion  = "us-east-1"
    bucketPath = "env0/audit"
    roleArn    = "arn:aws:iam::123456789012:role/env0-log-forwarding"
  })
  audit_log_forwarding = true
}
