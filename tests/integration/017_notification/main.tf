resource "env0_notification" "test_notification" {
  name  = var.second_run ? "Test-some-other-name" : "Test-Notification"
  type  = var.second_run ? "Slack" : "Teams"
  value = var.second_run ? "https://someotherurl.com" : "https://someurl.com"
}
