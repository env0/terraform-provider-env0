resource "tls_private_key" "throwaway" {
  algorithm = "RSA"
}
output "public_key_you_need_to_add_to_github_ssh_keys" {
  value = tls_private_key.throwaway.public_key_openssh
}

resource "env0_ssh_key" "tested" {
  name  = "test key"
  value = tls_private_key.throwaway.private_key_pem
}

data "env0_ssh_key" "tested" {
  name       = "test key"
  depends_on = [env0_ssh_key.tested]
}

data "env0_ssh_key" "tested2" {
  id = env0_ssh_key.tested.id
}

output "name" {
  value = data.env0_ssh_key.tested2.name
}

resource "env0_template" "usage" {
  name        = "use a ssh key"
  description = "use a ssh key"
  type        = "terraform"
  repository  = "https://github.com/env0/templates"
  ssh_keys    = [env0_ssh_key.tested]
}
