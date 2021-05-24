package env0

import (
	"context"
	"fmt"
	"github.com/env0/terraform-provider-env0/client"
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
	apiClient := meta.(*client.ApiClient)

	request := client.SshKeyCreatePayload{
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
	apiClient := meta.(*client.ApiClient)

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
	apiClient := meta.(*client.ApiClient)

	id := d.Id()
	err := apiClient.SshKeyDelete(id)
	if err != nil {
		return diag.Errorf("could not delete ssh key: %v", err)
	}
	return nil
}

func resourceSshKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(*client.ApiClient)

	id := d.Id()
	sshKeys, err := apiClient.SshKeys()
	if err != nil {
		return nil, err
	}

	for _, sshKey := range sshKeys {
		if sshKey.Id == id {
			d.Set("name", sshKey.Name)
			return []*schema.ResourceData{d}, nil
		}
	}
	return nil, fmt.Errorf("No ssh key for id %s", id)
}
