resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_git_token" "test_git_token" {
  name  = "name-${random_string.random.result}"
  value = "value"
}

data "env0_git_token" "test_git_token1" {
  name       = env0_git_token.test_git_token.name
  depends_on = [env0_git_token.test_git_token]
}

data "env0_git_token" "test_git_token2" {
  id = env0_git_token.test_git_token.id
}
