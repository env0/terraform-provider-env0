provider "random" {}

resource "random_string" "random" {
  length    = 5
  special   = false
  min_lower = 5
}

resource "env0_gpg_key" "test_gpg_key" {
  name    = "gpg-key-${random_string.random.result}"
  key_id  = "ABCDABCDABCDABCD"
  content = "dasdasdasd"
}

data "env0_gpg_key" "test_gpg_key_data" {
  name       = env0_gpg_key.test_gpg_key.name
  depends_on = [env0_gpg_key.test_gpg_key]
}

resource "env0_gpg_key" "test_gpg_key_modify" {
  name    = "gpg-key-${random_string.random.result}-2"
  key_id  = "ABCDABCDABCDABCD"
  content = var.second_run ? "dasdasdasd" : "dsadasdsvcxvcx"
}
