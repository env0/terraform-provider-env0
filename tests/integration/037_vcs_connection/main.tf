provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

# resource "env0_vcs_connection" "github" {
#   name          = "github-enterprise-${random_string.random.result}"
#   type          = "GitHubEnterprise"
#   url           = "https://github.example.com"
#   vcs_agent_key = var.second_run ? "ENV0_DEFAULT" : "custom-agent-key"
# }


