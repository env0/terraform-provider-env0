package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudCredentialsProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudCredentialsProjectAssignmentCreate,
		ReadContext:   resourceCloudCredentialsProjectAssignmentRead,
		DeleteContext: resourceCloudCredentialsProjectAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"credential_id": {
				Type:        schema.TypeString,
				Description: "id of cloud credentials",
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudCredentialsProjectAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	result, err := apiClient.AssignCloudCredentialsToProject(projectId, credentialId)
	if err != nil {
		return diag.Errorf("could not assign cloud credentials to project: %v", err)
	}
	d.SetId(getResourceId(result.CredentialId, result.ProjectId))
	return nil
}

func resourceCloudCredentialsProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	isNewResource := d.IsNewResource()
	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	credentialsList, err := apiClient.CloudCredentialIdsInProject(projectId)
	if err != nil {
		if !isNewResource {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not get cloud_credentials: %v", err)
	}
	found := false
	for _, candidate := range credentialsList {
		if candidate == credentialId {
			found = true
		}
	}
	if !found {
		if !isNewResource {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find cloud credential project assignment.\n project id = %v, cloud credentials id = %v", projectId, credentialId)
	}

	d.SetId(getResourceId(credentialId, projectId))

	return nil
}

func getResourceId(credentialId string, projectId string) string {
	return credentialId + "|" + projectId
}

func resourceCloudCredentialsProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	err := apiClient.RemoveCloudCredentialsFromProject(projectId, credentialId)
	if err != nil {
		return diag.Errorf("could not delete cloud credentials from project: %v", err)
	}
	return nil
}
