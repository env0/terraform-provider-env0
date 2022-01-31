package env0

import (
	"context"
	"errors"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
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
				Sensitive:   true,
			},
		},
	}
}

func resourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

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
	_, err := getSshKeyById(d.Id(), meta)
	if err != nil {
		if !d.IsNewResource() {
			d.SetId("")
			return nil
		}
		return err
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
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving SSH Key by id: ", id)
		_, getErr = getSshKeyById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving SSH Key by name: ", id)
		var sshKey client.SshKey
		sshKey, getErr = getSshKeyByName(id, meta)
		d.SetId(sshKey.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
