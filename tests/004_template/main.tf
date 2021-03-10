data "env0_template" "tested1" {
  name = "test"
}
output "tested1_template_id" {
  value = data.env0_template.tested1.id
}
output "tested1_template_type" {
  value = data.env0_template.tested1.type
}
