terraform {
  backend "local" {
  }
  required_providers {
    env0 = {
      source = "terraform.env0.com/local/env0"
    }
  }
}

provider "env0" {}

variable "second_run" {
  default = false
}