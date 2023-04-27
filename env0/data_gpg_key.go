package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGpgKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGpgKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the api key",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the api key",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"content": {
				Type:        schema.TypeString,
				Description: "the gpg public key block",
				Computed:    true,
			},
			"key_id": {
				Type:        schema.TypeString,
				Description: "the gpg key id",
				Computed:    true,
			},
		},
	}
}

func dataGpgKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var gpgKey *client.GpgKey
	var err error

	id, ok := d.GetOk("id")
	if ok {
		gpgKey, err = getGpgKeyById(id.(string), meta)
		if err != nil {
			return diag.Errorf("could not read gpg key: %v", err)
		}
	} else {
		gpgKey, err = getGpgKeyByName(d.Get("name").(string), meta)
		if err != nil {
			return diag.Errorf("could not read api key: %v", err)
		}
	}

	if err := writeResourceData(gpgKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
