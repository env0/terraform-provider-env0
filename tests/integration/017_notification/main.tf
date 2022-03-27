resource "env0_notification" "test_notification" {
  name  = var.second_run ? "Test-some-other-name" : "Test-Notification"
  type  = var.second_run ? "Slack" : "Teams"
  value = var.second_run ? "https://someotherurl.com" : "https://someurl.com"
}

resource "env0_project" "test_project" {
  name        = "Test-Project-For-Notification"
  description = "Test Description"
}

resource "env0_notification_project_assignment" "test_assignment" {
  project_id               = env0_project.test_project.id
  notification_endpoint_id = env0_notification.test_notification.id
  event_names              = var.second_run ? ["deploymentCancelled"] : ["environmentMarkedForAutoDestroy"]
}
