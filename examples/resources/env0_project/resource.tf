terraform {
  required_providers {
    env0 = {
      source = "env0/env0"
      version = "~> 0.0.13"
    }
  }
}

provider "env0" {}

resource "env0_project" "by_id" {
  name        = "example"
  description = "Example project"
}