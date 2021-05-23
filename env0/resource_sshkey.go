package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeyCreate,
		ReadContext:   resourceSshKeyRead,
		DeleteContext: resourceSshKeyDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceSshKeyImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the ssh key",
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "value is a private key in PEM format (first line usually looks like -----BEGIN OPENSSH PRIVATE KEY-----)",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	request := env0apiclient.SshKeyCreatePayload{
		Name:  d.Get("name").(string),
		Value: d.Get("value").(string),
	}
	sshKey, err := apiClient.SshKeyCreate(request)
	if err != nil {
		return diag.Errorf("could not create ssh key: %v", err)
	}

	d.SetId(sshKey.Id)

	return nil
}

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	sshKeys, err := apiClient.SshKeys()
	if err != nil {
		return diag.Errorf("could not query ssh keys: %v", err)
	}
	found := false
	for _, candidate := range sshKeys {
		if candidate.Id == d.Id() {
			found = true
		}
	}
	if !found {
		return diag.Errorf("ssh key %s not found", d.Id())
	}

	return nil
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)

	id := d.Id()
	err := apiClient.SshKeyDelete(id)
	if err != nil {
		return diag.Errorf("could not delete ssh key: %v", err)
	}
	return nil
}

func resourceSshKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, errors.New("Not implemented")
	// apiClient := meta.(*env0apiclient.ApiClient)

	// id := d.Id()
	// ssh key, err := apiClient.SshKey(id)
	// if err != nil {
	// 	return nil, err
	// }

	// d.Set("name", ssh key.Name)

	// return []*schema.ResourceData{d}, nil
}
