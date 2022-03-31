resource "env0_project" "example_project" {
  name = "project-example"
}

resource "env0_notification" "example_notification" {
  name  = "notification-example"
  type  = "Slack"
  value = "https://www.slack.com/example/webhook"
}

resource "env0_notification_project_assignment" "test_assignment" {
  project_id               = env0_project.example_project.id
  notification_endpoint_id = env0_notification.example_notification.id
  event_names              = ["environmentMarkedForAutoDestroy", "deploymentCancelled"]
}
