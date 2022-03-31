resource "env0_git_token" "git_token_sample" {
  name  = "name"
  value = "value"
}

data "env0_git_token" "git_token_sample_by_id" {
  id = env0_git_token.git_token_sample.id
}

data "env0_git_token" "git_token_sample_by_name" {
  name = env0_git_token.git_token_sample.name
}
