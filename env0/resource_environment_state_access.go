package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironmentStateAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentStateAccessCreate,
		ReadContext:   resourceEnvironmentStateAccessRead,
		DeleteContext: resourceEnvironmentStateAccessDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Description: "id of the environment",
				Required:    true,
				ForceNew:    true,
			},
			"accessible_from_entire_organization": {
				Type:        schema.TypeBool,
				Description: "when this parameter is 'false', allowed_project_ids should be provided. Defaults to 'false'",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"allowed_project_ids": {
				Type:        schema.TypeList,
				Description: "list of allowed project_ids. Used when 'accessible_from_entire_organization' is 'false'",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceEnvironmentStateAccessCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	environmentId := d.Get("environment_id").(string)

	var payload client.RemoteStateAccessConfigurationCreate
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if payload.AccessibleFromEntireOrganization && payload.AllowedProjectIds != nil {
		return diag.Errorf("'allowed_project_ids' should not be set when 'accessible_from_entire_organization' is set to 'true'")
	}

	apiClient := meta.(client.ApiClientInterface)

	remoteStateAccess, err := apiClient.RemoteStateAccessConfigurationCreate(environmentId, payload)
	if err != nil {
		return diag.Errorf("could not create a remote state access configation: %v", err)
	}

	d.SetId(remoteStateAccess.EnvironmentId)

	return nil
}

func resourceEnvironmentStateAccessRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	environmentId := d.Get("environment_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	remoteStateAccess, err := apiClient.RemoteStateAccessConfiguration(environmentId)
	if err != nil {
		return ResourceGetFailure(ctx, "remote state access configation", d, err)
	}

	if err := writeResourceData(remoteStateAccess, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceEnvironmentStateAccessDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	environmentId := d.Get("environment_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.RemoteStateAccessConfigurationDelete(environmentId); err != nil {
		return diag.Errorf("could not delete remote state access configation: %v", err)
	}

	return nil
}
