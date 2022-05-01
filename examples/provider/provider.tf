# Terraform 0.13+ uses the Terraform Registry:

terraform {
  required_providers {
    env0 = {
      source = "env0/env0"
    }
  }
}

# Configure the env0 provider

provider "env0" {
  api_key    = var.env0_api_key
  api_secret = var.env0_api_secret
}
