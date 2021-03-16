data "env0_project" "default" {
  name = "Default Organization Project"
}

data "env0_project" "default2" {
  depends_on = [data.env0_project.default]
  id         = data.env0_project.default.id
}

output "default_project_id" {
  value = data.env0_project.default.id
}

output "default_project_name" {
  value = data.env0_project.default2.name
}
