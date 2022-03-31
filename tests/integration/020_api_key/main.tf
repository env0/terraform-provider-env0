resource "env0_api_key" "test_api_key" {
  name = "name"
}

data "env0_api_key" "test_api_key1" {
  name       = env0_api_key.test_api_key.name
  depends_on = [env0_api_key.test_api_key]
}

data "env0_api_key" "test_api_key2" {
  id = env0_api_key.test_api_key.id
}
