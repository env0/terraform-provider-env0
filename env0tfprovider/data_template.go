package env0tfprovider

import (
	"context"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the template",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the template",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "template source code repository url",
				Computed:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "terraform / terrgrunt folder inside source code repository",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "'terraform' or 'terragrunt'",
				Computed:    true,
			},
		},
	}
}

func dataTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	name, nameSpecified := d.GetOk("name")
	var template env0apiclient.Template
	var err error
	if nameSpecified {
		templates, err := apiClient.Templates()
		if err != nil {
			return diag.Errorf("Could not query templates: %v", err)
		}
		for _, candidate := range templates {
			if candidate.Name == name {
				template = candidate
			}
		}
		if template.Name == "" {
			return diag.Errorf("Could not find an env0 template with name %s", name)
		}
	} else {
		template, err = apiClient.Template(d.Id())
		if err != nil {
			return diag.Errorf("Could not query template: %v", err)
		}
	}

	d.SetId(template.Id)
	d.Set("name", template.Name)
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("type", template.Type)
	//TODO: retry, sshkeys
	return nil
}
