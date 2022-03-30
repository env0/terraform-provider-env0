resource "env0_git_token" "test_git_token" {
  name  = "name"
  value = "value"
}

data "env0_git_token" "test_git_token1" {
  name       = env0_git_token.test_git_token.name
  depends_on = [env0_git_token.test_git_token]
}

data "env0_git_token" "test_git_token2" {
  id = env0_git_token.test_git_token.id
}

output "value1" {
  value = data.env0_git_token.test_git_token1.value
}

output "value2" {
  value = data.env0_git_token.test_git_token2.value
}
