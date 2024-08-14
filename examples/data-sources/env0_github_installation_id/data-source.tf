data "env0_github_installation_id" "example" {
  repository = "https://github.com/env0/templates"
}

output "github_installation_id" {
  value = data.env0_github_installation_id.example.github_installation_id
}
