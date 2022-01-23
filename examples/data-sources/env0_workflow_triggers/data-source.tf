data "env0_workflow_triggers" "downstream_envs" {
  environment_id = "environment_id"
}

data "env0_environment" "by_name" {
  id = data.env0_workflow_triggers.downstream_envs.0.id
}

output "downstream_env_name" {
  value = data.env0_environment.by_name.name
}
