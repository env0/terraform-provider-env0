resource "env0_oci_credentials" "example" {
  name         = "my-oci-credentials"
  tenancy_ocid = "ocid1.tenancy.oc1..exampleuniqueID"
  user_ocid    = "ocid1.user.oc1..exampleuniqueID"
  fingerprint  = "12:34:56:78:90:ab:cd:ef:12:34:56:78:90:ab:cd:ef"
  private_key  = <<EOF
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASC...
-----END PRIVATE KEY-----
EOF
  region       = "us-phoenix-1"
}
