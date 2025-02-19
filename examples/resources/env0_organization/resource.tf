# Create a basic organization
resource "env0_organization" "example" {
  name        = "example-org"
  description = "Example organization created via Terraform"
}