data "env0_user" "user_by_email_example" {
  email = "john.doe@email.com"
}

output "user_id" {
  value = data.env0_user.user_by_email_exmple.id
}
