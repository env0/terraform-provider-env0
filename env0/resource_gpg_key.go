package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGpgKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGpgKeyCreate,
		ReadContext:   resourceGpgKeyRead,
		DeleteContext: resourceGpgKeyDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceGpgKeyImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the gpg key name",
				Required:    true,
				ForceNew:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "the gpg public key block",
				Required:    true,
				ForceNew:    true,
			},
			"key_id": {
				Type:             schema.TypeString,
				Description:      "the gpg key id",
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: NewRegexValidator(`[0-9A-F]{16}`),
			},
		},
	}
}

func resourceGpgKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.GpgKeyCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	gpgKey, err := apiClient.GpgKeyCreate(&payload)
	if err != nil {
		return diag.Errorf("could not create gpg key: %v", err)
	}

	d.SetId(gpgKey.Id)

	return nil
}

func resourceGpgKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	gpgKey, err := getGpgKeyById(d.Id(), meta)
	if err != nil {
		return ResourceGetFailure(ctx, "gpg key", d, err)
	}

	if err := writeResourceData(gpgKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceGpgKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.GpgKeyDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete gpg key: %v", err)
	}

	return nil
}

func getGpgKeyById(id string, meta interface{}) (*client.GpgKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	gpgKeys, err := apiClient.GpgKeys()
	if err != nil {
		return nil, err
	}

	for _, gpgKey := range gpgKeys {
		if gpgKey.Id == id {
			return &gpgKey, nil
		}
	}

	return nil, &client.NotFoundError{}
}

func getGpgKeyByName(name string, meta interface{}) (*client.GpgKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	gpgKeys, err := apiClient.GpgKeys()
	if err != nil {
		return nil, err
	}

	var foundGpgKeys []client.GpgKey

	for _, gpgKey := range gpgKeys {
		if gpgKey.Name == name {
			foundGpgKeys = append(foundGpgKeys, gpgKey)
		}
	}

	if len(foundGpgKeys) == 0 {
		return nil, fmt.Errorf("gpg key with name %v not found", name)
	}

	if len(foundGpgKeys) > 1 {
		return nil, fmt.Errorf("found multiple gpg keys with name: %s. Use id instead or make sure gpg key names are unique %v", name, foundGpgKeys)
	}

	return &foundGpgKeys[0], nil
}

func getGpgKey(ctx context.Context, idOrName string, meta interface{}) (*client.GpgKey, error) {
	_, err := uuid.Parse(idOrName)
	if err == nil {
		tflog.Info(ctx, "Resolving gpg key by id", map[string]interface{}{"id": idOrName})

		return getGpgKeyById(idOrName, meta)
	} else {
		tflog.Info(ctx, "Resolving gpg key by name", map[string]interface{}{"name": idOrName})

		return getGpgKeyByName(idOrName, meta)
	}
}

func resourceGpgKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	gpgKey, err := getGpgKey(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(gpgKey, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
