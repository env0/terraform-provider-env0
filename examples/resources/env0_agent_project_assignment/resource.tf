resource "env0_project" "example" {
  name        = "example"
  description = "Example project"
}

data "env0_agents" "agents" {}

resource "env0_agent_project_assignment" {
    agent_id = data.env0_agents.agents.0.agent_key
    project_id = env0_project.example.id
}
