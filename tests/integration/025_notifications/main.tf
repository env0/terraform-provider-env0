provider "random" {}

resource "random_string" "random" {
  length    = 20
  special   = false
  min_lower = 20
}

locals {
  notification_name_prefix = "integration-test-025-notification"
}

resource "env0_notification" "test_notification_1" {
  name           = "${local.notification_name_prefix}-1-${random_string.random.result}"
  type           = "Webhook"
  value          = "https://someurl1.com"
  webhook_secret = "my_little_secret"
}

resource "env0_notification" "test_notification_2" {
  name  = "${local.notification_name_prefix}-2-${random_string.random.result}"
  type  = "Teams"
  value = "https://someurl2.com"
}

data "env0_notifications" "all_notifications" {
  depends_on = [env0_notification.test_notification_1, env0_notification.test_notification_2]
}

data "env0_notification" "test_notification_1" {
  depends_on = [env0_notification.test_notification_1]
  name       = "${local.notification_name_prefix}-1-${random_string.random.result}"
}

data "env0_notification" "test_notification_2" {
  depends_on = [env0_notification.test_notification_2]
  name       = "${local.notification_name_prefix}-2-${random_string.random.result}"
}

output "notification_1_from_all_notifications" {
  value = replace(
    data.env0_notifications.all_notifications.names[
      index(data.env0_notifications.all_notifications.names,
      env0_notification.test_notification_1.name)
    ]
    , random_string.random.result
    , ""
  )
}
