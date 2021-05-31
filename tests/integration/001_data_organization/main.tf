data "env0_organization" "my_organization" {}

output "organization_name" {
  value = data.env0_organization.my_organization.name
}
