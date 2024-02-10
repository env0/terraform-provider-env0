data "env0_module_testing_project" "project" {}

data "env0_agents" "agents" {}

resource "env0_agent_project_assignment" "example" {
  agent_id   = data.env0_agents.agents.0.agent_key
  project_id = env0_module_testing_project.project.id
}
