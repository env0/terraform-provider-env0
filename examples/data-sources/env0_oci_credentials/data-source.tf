data "env0_oci_credentials" "example" {
  name = "my-oci-credentials"
}

output "oci_credentials_id" {
  value = data.env0_oci_credentials.example.id
}
