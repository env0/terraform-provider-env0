resource "env0_gpg_key" "example" {
  name    = "gpg-key-example"
  key_id  = "ABCDABCDABCDABCD"
  content = "key block"
}
