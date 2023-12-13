resource "env0_notification" "notification" {
  name  = "notification"
  type  = "Slack"
  value = "https://someurl.com"
}

resource "env0_project" "project" {
  name = "project"
}

resource "env0_notification_project_assignment" "notification_project_assignment" {
  project_id               = env0_project.project.id
  notification_endpoint_id = env0_notification.notification.id
  event_names              = ["environmentMarkedForAutoDestroy"]
}
