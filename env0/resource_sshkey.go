package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeyCreate,
		ReadContext:   resourceSshKeyRead,
		UpdateContext: resourceSshKeyUpdate,
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
				Sensitive:   true,
			},
		},
	}
}

func resourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.SshKeyCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	sshKey, err := apiClient.SshKeyCreate(payload)
	if err != nil {
		return diag.Errorf("could not create ssh key: %v", err)
	}

	d.SetId(sshKey.Id)

	return nil
}

func resourceSshKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.SshKeyUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.SshKeyUpdate(d.Id(), &payload); err != nil {
		return diag.Errorf("could not update ssh key: %v", err)
	}

	return nil
}

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sshKey, err := getSshKeyById(d.Id(), meta)
	if err != nil {
		return diag.Errorf("could not get ssh key: %v", err)
	}

	if sshKey == nil {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
	}

	return nil
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	err := apiClient.SshKeyDelete(id)
	if err != nil {
		return diag.Errorf("could not delete ssh key: %v", err)
	}

	return nil
}

func resourceSshKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	var getErr error

	_, uuidErr := uuid.Parse(id)

	if uuidErr == nil {
		tflog.Info(ctx, "Resolving SSH key by id", map[string]interface{}{"id": id})
		_, getErr = getSshKeyById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving SSH key by name", map[string]interface{}{"name": id})

		var sshKey *client.SshKey

		sshKey, getErr = getSshKeyByName(id, meta)
		d.SetId(sshKey.Id)
	}

	if getErr != nil {
		return nil, getErr
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
