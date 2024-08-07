resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "tls_private_key" "throwaway" {
  algorithm = "RSA"
}
output "public_key_you_need_to_add_to_github_ssh_keys" {
  value = tls_private_key.throwaway.public_key_openssh
}

resource "env0_ssh_key" "tested" {
  name  = "test-key-${random_string.random.result}"
  value = tls_private_key.throwaway.private_key_pem
}

data "env0_ssh_key" "tested" {
  name       = "test-key-${random_string.random.result}"
  depends_on = [env0_ssh_key.tested]
}

data "env0_ssh_key" "tested2" {
  id = env0_ssh_key.tested.id
}

output "name" {
  value = replace(env0_ssh_key.tested.name, random_string.random.result, "")
}

resource "env0_template" "usage" {
  name        = "use-a-ssh-key-${random_string.random.result}"
  description = "use a ssh key"
  type        = "terraform"
  repository  = "https://github.com/env0/templates"
  ssh_keys    = [env0_ssh_key.tested]
  terraform_version = "1.3.1"
}
