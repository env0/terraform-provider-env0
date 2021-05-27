data "env0_ssh_key" "my_key" {
  name = "Secret Key"
}

resource "env0_template" "example" {
  # ...
  ssh_keys = [data.env0_ssh_key.my_key]
}
