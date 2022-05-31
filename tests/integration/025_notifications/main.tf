provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}


resource "env0_notification" "test_notification_1" {
  name  = "notification-${random_string.random.result}-1"
  type  = "Slack"
  value = "https://someurl1.com"
}

resource "env0_notification" "test_notification_2" {
  name  = "notification-${random_string.random.result}-2"
  type  = "Teams"
  value = "https://someurl2.com"
}

data "env0_notifications" "all_notifications" {}

data "env0_notification" "notification" {
  for_each = toset(data.env0_notifications.all_notifications.names)
  name     = each.value
}
