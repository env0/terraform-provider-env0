data "env0_templates" "all_templates" {}

data "env0_template" "templates" {
  for_each = toset(data.env0_templates.all_templates.names)
  name     = each.value
}

output "template1_name" {
  value = data.env0_template.templates["Github Test-111"].name
}

output "template2_name" {
  value = data.env0_template.templates["Github Test-222"].name
}
