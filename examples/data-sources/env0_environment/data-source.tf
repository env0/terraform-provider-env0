data "env0_environment" "by_name" {
  name   =  "best-env"
}

output "environment_project_id" {
  value = data.env0_environment.by_name.project_id
}

data "env0_environment" "by_id" {
  id   =  "some_id"
}

output "environment_name" {
  value = data.env0_environment.by_id.name
}
