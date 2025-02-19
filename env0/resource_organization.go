package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		UpdateContext: resourceOrganizationUpdate,
		DeleteContext: schema.NoopContext, // Organization deletion is not supported

		Description: "Manages an env0 Organization. Note: Organization deletion is not supported via the API.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of the organization",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "a description of the organization",
			},
			"photo_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL to the organization's photo or logo",
			},
		},
	}
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.OrganizationCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	organization, err := apiClient.OrganizationCreate(payload)
	if err != nil {
		return diag.Errorf("could not create organization: %v", err)
	}

	d.SetId(organization.Id)

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	organization, err := apiClient.OrganizationById(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "organization", d, err)
	}

	if err := writeResourceData(organization, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.OrganizationUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	organization, err := apiClient.OrganizationUpdate(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update organization: %v", err)
	}

	if err := writeResourceData(organization, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
