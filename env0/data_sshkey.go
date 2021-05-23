package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSshKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSshKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the ssh key",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the ssh key",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	name, nameSpecified := d.GetOk("name")
	var sshKey env0apiclient.SshKey
	if nameSpecified {
		sshKeys, err := apiClient.SshKeys()
		if err != nil {
			return diag.Errorf("Could not query ssh keys: %v", err)
		}
		for _, candidate := range sshKeys {
			if candidate.Name == name {
				sshKey = candidate
			}
		}
		if sshKey.Name == "" {
			return diag.Errorf("Could not find an env0 ssh key with name %s", name)
		}
	} else {
		id, idSpecified := d.GetOk("id")
		if !idSpecified {
			return diag.Errorf("At lease one of 'id', 'name' must be specified")
		}
		sshKeys, err := apiClient.SshKeys()
		if err != nil {
			return diag.Errorf("Could not query ssh keys: %v", err)
		}
		for _, candidate := range sshKeys {
			if candidate.Id == id.(string) {
				sshKey = candidate
			}
		}
		if sshKey.Name == "" {
			return diag.Errorf("Could not find an env0 ssh key with id %s", id)
		}
	}

	d.SetId(sshKey.Id)
	d.Set("name", sshKey.Name)

	return nil
}
