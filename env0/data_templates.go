package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTemplatesRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Description: "list of all templates (by name)",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the template name",
				},
			},
		},
	}
}

func dataTemplatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	templates, err := apiClient.Templates()
	if err != nil {
		return diag.Errorf("Could not get templates: %v", err)
	}

	data := []string{}

	for _, template := range templates {
		if !template.IsDeleted {
			data = append(data, template.Name)
		}
	}

	d.Set("names", data)

	// Not really needed. But required by Terraform SDK - https://github.com/hashicorp/terraform-plugin-sdk/issues/541
	d.SetId("all_templates_names")

	return nil
}
