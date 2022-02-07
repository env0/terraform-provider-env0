terraform {
  backend "local" {
  }
  required_providers {
    env0 = {
      source = "terraform-registry.env0.com/env0/env0"
    }
  }
}

provider "env0" {
  api_key = "tabarvanrk3p0xgd"
  api_secret = "jDYuSMFQIVrbYjw7nx9Gh9-1WkGSDaQt"
}

variable "second_run" {
  default = false
}
