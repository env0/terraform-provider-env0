package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCustomRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCustomRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the custom role",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the custom role",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataCustomRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error

	var role *client.Role

	id, ok := d.GetOk("id")
	if ok {
		role, err = getCustomRoleById(id.(string), meta)
	} else {
		role, err = getCustomRoleByName(d.Get("name").(string), meta)
	}

	if err != nil {
		return DataGetFailure("role", id, err)
	}

	if err := writeResourceData(role, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
