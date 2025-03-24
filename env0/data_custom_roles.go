package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCustomRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCustomRolesRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Description: "list of all custom roles (by name)",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the role name",
				},
			},
		},
	}
}

func dataCustomRolesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	roles, err := apiClient.Roles()
	if err != nil {
		return diag.Errorf("Failed to get custom roles: %v", err)
	}

	data := []string{}

	for _, role := range roles {
		if !role.IsDefaultRole {
			data = append(data, role.Name)
		}
	}

	d.Set("names", data)

	// Not really needed. But required by Terraform SDK - https://github.com/hashicorp/terraform-plugin-sdk/issues/541
	d.SetId("all_roles_names")

	return nil
}
