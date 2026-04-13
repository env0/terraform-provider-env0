# Hosting deployment logs in your AWS account

data "aws_caller_identity" "this" {}
data "aws_region" "this" {}
data "env0_organization" "this" {}

resource "env0_agent_pool" "agent_with_logs_config" {
  name = "agent with logs config"
  logs {
    account_id  = data.aws_caller_identity.this.account_id
    region      = data.aws_region.this.region
    external_id = data.env0_organization.this.id
  }
}

module "agent_logs" {
  source = "github.com/env0/k8s-modules?ref=v1.1.0//log-storage/aws/dynamodb"

  agent_key   = env0_agent_pool.agent_with_logs_config.id
  external_id = data.env0_organization.this.id
}