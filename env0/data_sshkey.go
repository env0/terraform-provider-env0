package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
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
	name, nameSpecified := d.GetOk("name")
	var sshKey client.SshKey
	var err diag.Diagnostics
	if nameSpecified {
		sshKey, err = getSshKeyByName(name, meta)
		if err != nil {
			return err
		}
	} else {
		id, idSpecified := d.GetOk("id")
		if !idSpecified {
			return diag.Errorf("At lease one of 'id', 'name' must be specified")
		}
		sshKey, err = getSshKeyById(id, meta)
		if err != nil {
			return err
		}

	}

	d.SetId(sshKey.Id)
	d.Set("name", sshKey.Name)

	return nil
}

func getSshKeyByName(name interface{}, meta interface{}) (client.SshKey, diag.Diagnostics) {
	apiClient := meta.(*client.ApiClient)

	sshKeys, err := apiClient.SshKeys()
	if err != nil {
		return client.SshKey{}, diag.Errorf("Could not query ssh keys: %v", err)
	}

	var sshKeysByName []client.SshKey
	for _, candidate := range sshKeys {
		if candidate.Name == name {
			sshKeysByName = append(sshKeysByName, candidate)
		}
	}

	if len(sshKeysByName) > 1 {
		return client.SshKey{}, diag.Errorf("Found multiple SSH Keys for name: %s. Use ID instead or make sure SSH Keys names are unique %v", name, sshKeysByName)
	}
	if len(sshKeysByName) == 0 {
		return client.SshKey{}, diag.Errorf("Could not find an env0 ssh key with name %s", name)
	}

	return sshKeysByName[0], nil
}

func getSshKeyById(id interface{}, meta interface{}) (client.SshKey, diag.Diagnostics) {
	apiClient := meta.(*client.ApiClient)

	sshKeys, err := apiClient.SshKeys()
	var sshKey client.SshKey
	if err != nil {
		return client.SshKey{}, diag.Errorf("Could not query ssh keys: %v", err)
	}

	for _, candidate := range sshKeys {
		if candidate.Id == id.(string) {
			sshKey = candidate
		}
	}
	if sshKey.Name == "" {
		return client.SshKey{}, diag.Errorf("Could not find an env0 ssh key with id %s", id)
	}
	return sshKey, nil
}
