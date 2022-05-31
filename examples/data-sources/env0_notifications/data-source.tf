data "env0_notifications" "all_notifications" {}

data "env0_notification" "notifications" {
  for_each = toset(data.env0_notifications.all_notifications.names)
  name     = each.value
}

output "notification1_name" {
  value = data.env0_notification.notifications["my notification 345"].name
}

output "notification2_name" {
  value = data.env0_notification.notifications["my notification 123"].name
}
