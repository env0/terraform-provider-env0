package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCodeVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeVariablesRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Description: "extracts source code terraform variables from the VCS configuration of this template",
				Required:    true,
			},

			"variables": {
				Type:        schema.TypeList,
				Description: "a list of terraform variables extracted from the source code",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "the name of the terraform variable",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "the value of the terraform variable",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "the description of the terraform variable",
							Computed:    true,
						},
						"format": {
							Type:        schema.TypeString,
							Description: "the format of the terraform variable (HCL or JSON)",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeVariablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	templateId := d.Get("template_id").(string)

	template, err := apiClient.Template(templateId)
	if err != nil {
		return diag.Errorf("could not get template: %v", err)
	}

	payload := &client.VariablesFromRepositoryPayload{
		BitbucketClientKey:   template.BitbucketClientKey,
		GithubInstallationId: template.GithubInstallationId,
		Path:                 template.Path,
		Revision:             template.Revision,
		TokenId:              template.TokenId,
		Repository:           template.Repository,
	}

	variables, err := apiClient.VariablesFromRepository(payload)
	if err != nil {
		return diag.Errorf("failed to extract variables from repository: %v", err)
	}

	ivalues, err := writeResourceDataGetSliceValues(variables, "variables", d)
	if err != nil {
		return diag.Errorf("schema slice resource data serialization failed: %v", err)
	}

	for i, ivalue := range ivalues {
		if variables[i].IsSensitive == nil || !*variables[i].IsSensitive {
			ivariable := ivalue.(map[string]interface{})
			ivariable["value"] = variables[i].Value
		}
	}

	d.Set("variables", ivalues)

	d.SetId(templateId)

	return nil
}
