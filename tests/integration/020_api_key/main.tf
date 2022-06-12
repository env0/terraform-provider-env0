provider "random" {}

resource "random_string" "random" {
  length    = 5
  special   = false
  min_lower = 5
}

resource "env0_api_key" "test_api_key" {
  name = "my-little-api-key-${random_string.random.result}"
}

resource "env0_api_key" "test_user_api_key" {
  name = "my-little-user-api-key-${random_string.random.result}"
  organization_role = "User"
}

data "env0_api_key" "test_api_key1" {
  name       = env0_api_key.test_api_key.name
  depends_on = [env0_api_key.test_api_key]
}

data "env0_api_key" "test_api_key2" {
  id = env0_api_key.test_api_key.id
}
