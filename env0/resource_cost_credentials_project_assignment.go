package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCostCredentialsProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudCredentialsProjectAssignmentCreate,
		ReadContext:   resourceCloudCredentialsProjectAssignmentRead,
		DeleteContext: resourceCloudCredentialsProjectAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"credential_id": {
				Type:        schema.TypeString,
				Description: "id of cost credentials",
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

func resourceCostCredentialsProjectAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	result, err := apiClient.AssignCostCredentialsToProject(projectId, credentialId)
	if err != nil {
		return diag.Errorf("could not assign cost credentials to project: %v", err)
	}
	d.SetId(getResourceId(result.CredentialsId, result.ProjectId))
	return nil
}

func resourceCostdCredentialsProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	credentialsList, err := apiClient.CostCredentialIdsInProject(projectId)
	if err != nil {
		return diag.Errorf("could not get cost_credentials: %v", err)
	}
	found := false
	for _, candidate := range credentialsList {
		if candidate.CredentialsId == credentialId {
			found = true
		}
	}
	if !found && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	d.SetId(getResourceId(credentialId, projectId))

	return nil
}

func resourceCostCredentialsProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentialId := d.Get("credential_id").(string)
	projectId := d.Get("project_id").(string)
	err := apiClient.RemoveCostCredentialsFromProject(projectId, credentialId)
	if err != nil {
		return diag.Errorf("could not delete cost credentials from project: %v", err)
	}
	return nil
}
