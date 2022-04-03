resource "env0_api_key" "api_key_sample" {
  name  = "name"
}

data "env0_api_key" "api_key_sample_by_id" {
  id = env0_api_key.api_key_sample.id
}

data "env0_api_key" "api_key_sample_by_name" {
  name = env0_api_key.api_key_sample.name
}
