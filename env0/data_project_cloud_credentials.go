package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjectCloudCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectCloudCredentialsRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "the project id for listing the cloud credentials",
				Required:    true,
			},

			"ids": {
				Type:        schema.TypeList,
				Description: "a list of cloud credentials (ids) associated with the project",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the cloud credential's id",
				},
			},
		},
	}
}

func dataProjectCloudCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)

	credentialIds, err := apiClient.CloudCredentialIdsInProject(projectId)
	if err != nil {
		return diag.Errorf("could not get cloud credentials associated with project %s: %v", projectId, err)
	}

	d.Set("ids", credentialIds)
	d.SetId(projectId)

	return nil
}
