package env0

import (
	"context"
	"fmt"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProviderCreate,
		ReadContext:   resourceProviderRead,
		UpdateContext: resourceProviderUpdate,
		DeleteContext: resourceProviderDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceProviderImport},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:             schema.TypeString,
				Description:      `type of the provider registry (Match pattern: ^[0-9a-z](?:[0-9a-z-]{0,30}[0-9a-z])?$). Your provider’s type is essentially it’s name, and should match your provider’s files. For example, if your binaries look like terraform-provider-aws_1.1.1_linux_amd64.zip, than your provider’s type should be aws.`,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: NewRegexValidator(`^[0-9a-z](?:[0-9a-z-]{0,30}[0-9a-z])?$`),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the provider registry",
				Optional:    true,
			},
		},
	}
}

func resourceProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ProviderCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	provider, err := apiClient.ProviderCreate(payload)
	if err != nil {
		return diag.Errorf("could not create provider: %v", err)
	}

	d.SetId(provider.Id)

	return nil
}

func resourceProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	provider, err := apiClient.Provider(d.Id())
	if err != nil {
		return ResourceGetFailure("provider", d, err)
	}

	if err := writeResourceData(provider, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ProviderUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.ProviderUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update provider: %v", err)
	}

	return nil
}

func resourceProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ProviderDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete provider: %v", err)
	}

	return nil
}

func getProviderByName(name string, meta interface{}) (*client.Provider, error) {
	apiClient := meta.(client.ApiClientInterface)

	providers, err := apiClient.Providers()
	if err != nil {
		return nil, err
	}

	var foundProviders []client.Provider
	for _, provider := range providers {
		if provider.Type == name {
			foundProviders = append(foundProviders, provider)
		}
	}

	if len(foundProviders) == 0 {
		return nil, fmt.Errorf("provider with name/type %v not found", name)
	}

	if len(foundProviders) > 1 {
		return nil, fmt.Errorf("found multiple providers with name/type: %s. Use id instead or make sure provider names are unique %v", name, foundProviders)
	}

	return &foundProviders[0], nil
}

func getProvider(idOrName string, meta interface{}) (*client.Provider, error) {
	_, err := uuid.Parse(idOrName)
	if err == nil {
		log.Println("[INFO] Resolving provider by id: ", idOrName)
		return meta.(client.ApiClientInterface).Provider(idOrName)
	} else {
		log.Println("[INFO] Resolving provider by name: ", idOrName)
		return getProviderByName(idOrName, meta)
	}
}

func resourceProviderImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	provider, err := getProvider(d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(provider, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
