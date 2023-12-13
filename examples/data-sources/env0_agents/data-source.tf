data "env0_agent_values" "agent_values" {
  agent_key = "pr12"
}

output "values" {
  value = data.env0_agent_values.agent_values.values
}
