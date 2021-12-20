terraform {
  backend "local" {
  }
  required_providers {
    env0 = {
      source = "terraform-registry.env0.com/env0/env0"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.7.2"
    }
  }
}

provider "env0" {}

variable "second_run" {
  default = false
}
provider "time" {}
