provider "random" {}

resource "random_string" "random" {
  length    = 5
  special   = false
  min_lower = 5
}

resource "env0_provider" "test_provider" {
  type        = "aws-${random_string.random.result}"
  description = var.second_run ? "des1" : "des2"
}

data "env0_provider" "test_provider_data" {
  type = env0_provider.test_gpg_key.name
}

resource "env0_provider" "test_provider-type-change" {
  type        = var.second_run ? "aws2-${random_string.random.result}" : "aws1-${random_string.random.result}"
  description = var.second_run ? "des1" : "des2"
}
