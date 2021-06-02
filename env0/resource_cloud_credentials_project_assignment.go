package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudCredentialsProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudCredentialsProjectAssignmenetCreate,
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

func resourceCloudCredentialsProjectAssignmenetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	result, err := apiClient.AssignCloudCredentialsToProject(projectId, credentialId)
	if err != nil {
		return diag.Errorf("could not assign cloud credentials to project: %v", err)
	}
	d.SetId(result.Id)
	return nil
}

func resourceCloudCredentialsProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	credentialsList, err := apiClient.CloudCredentialProjectAssignments(projectId)
	if err != nil {
		return diag.Errorf("could not get cloud_credentials: %v", err)
	}
	var credentials *client.CloudCredentialsProjectAssignment
	for _, candidate := range credentialsList {
		if candidate.CredentialId == credentialId {
			credentials = &candidate
		}
	}
	if credentials == nil {
		return diag.Errorf("could not find cloud credential project assignment.\n project id = %v, cloud credentials id = %v", projectId, credentialId)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceCloudCredentialsProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	credential_id := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	err := apiClient.RemoveCloudCredentialsFromProject(credential_id, projectId)
	if err != nil {
		return diag.Errorf("could not delete cloud credentials from project: %v", err)
	}
	return nil
}
