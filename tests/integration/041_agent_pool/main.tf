provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_agent_pool" "test_pool" {
  name        = "Test-Agent-Pool-041-${random_string.random.result}"
  description = var.second_run ? "updated description" : "initial description"
}

data "env0_agent_pool" "by_name" {
  name       = env0_agent_pool.test_pool.name
  depends_on = [env0_agent_pool.test_pool]
}

data "env0_agent_pool" "by_id" {
  id         = env0_agent_pool.test_pool.id
  depends_on = [env0_agent_pool.test_pool]
}

resource "env0_agent_secret" "test_secret" {
  agent_id    = env0_agent_pool.test_pool.id
  description = "integration test secret"
}

output "agent_pool_name" {
  value = env0_agent_pool.test_pool.name
}

output "agent_pool_description" {
  value = env0_agent_pool.test_pool.description
}

output "data_by_name_description" {
  value = data.env0_agent_pool.by_name.description
}

output "data_by_id_description" {
  value = data.env0_agent_pool.by_id.description
}

output "secret_agent_id" {
  value = env0_agent_secret.test_secret.agent_id
}
