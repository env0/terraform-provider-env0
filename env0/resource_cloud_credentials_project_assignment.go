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
	result, err := apiClient.AssignCloudCredentialsToProject(credentialId, request)
	if err != nil {
		return diag.Errorf("could not assign cloud credentials to project: %v", err)
	}
	resourceId := result.Id + "|" + projectId
	d.SetId(resourceId)
	return nil
}

func resourceCloudCredentialsProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	credentialId := d.Get("credential_id").(string)
	cloud_credentials, err := apiClient.CloudCredentials(credentialId)
	if err != nil {
		return diag.Errorf("could not get cloud_credentials: %v", err)
	}
	var assignProjectId = d.Get("project_id").(string)
	isProjectIdInCloud_credentials := false
	for _, projectId := range cloud_credentials.ProjectIds {
		if assignProjectId == projectId {
			isProjectIdInCloud_credentials = true
		}
	}
	if !isProjectIdInCloud_credentials {
		return diag.Errorf("could not find projectId in cloud credentials.\n projectId = %v, cloud credentialsId = %v", assignProjectId, credentialId)

	}

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
