provider "random" {}

resource "random_string" "random" {
  length    = 5
  special   = false
  min_lower = 5
}

resource "random_id" "random_key_id" {
  byte_length = 7
}

resource "env0_gpg_key" "test_gpg_key" {
  name    = "gpg-key-${random_string.random.result}"
  key_id  = upper("${random_id.random_key_id.hex}CD")
  content = "dasdasdasd"
}

data "env0_gpg_key" "test_gpg_key_data" {
  name = env0_gpg_key.test_gpg_key.name
}

resource "env0_gpg_key" "test_gpg_key_modify" {
  name    = "gpg-key-${random_string.random.result}-2"
  key_id  = upper("${random_id.random_key_id.hex}AB")
  content = var.second_run ? "dasdasdasd" : "dsadasdsvcxvcx"
}
