resource "env0_git_token" "git_token_sample" {
  name  = "name"
  value = "value"
}

data "env0_git_token" "git_token_sample" {
  id = env0_git_token.git_token_sample.id
}

output "value" {
  value = data.env0_git_token.git_token_sample.value
}
