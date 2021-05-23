package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataOrganizationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the name of the organization",
				Computed:    true,
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "textual description of the entity who created the organization",
				Computed:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "role of the authenticated user (through api key) in the organization",
				Computed:    true,
			},
			"is_self_hosted": {
				Type:        schema.TypeBool,
				Description: "is the organization self hosted",
				Computed:    true,
			},
		},
	}
}

func dataOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	organization, err := apiClient.Organization()
	if err != nil {
		return diag.Errorf("Could not query organization: %v", err)
	}

	d.SetId(organization.Id)
	d.Set("name", organization.Name)
	d.Set("created_by", organization.CreatedBy)
	d.Set("role", organization.Role)
	d.Set("is_self_hosted", organization.IsSelfHosted)

	return nil
}
