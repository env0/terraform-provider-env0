package env0

import (
	"context"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTemplateProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateProjectAssignmenetCreate,
		ReadContext:   resourceTemplateProjectAssignmentRead,
		//UpdateContext: resourceTemplateProjectAssignmentUpdate,
		DeleteContext: resourceTemplateProjectAssignmentDelete,


		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Description: "id of the template",
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

func templateProjectAssignmentPayloadFromParameters(d *schema.ResourceData) client.TemplateAssignmentToProjectPayload {
	result := client.TemplateAssignmentToProjectPayload{
		ProjectId: d.Get("project_id").(string),
	}

	return result
}

func resourceTemplateProjectAssignmenetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Eyal!! start")
	apiClient := meta.(*client.ApiClient)
	templateId := d.Get("template_id").(string)
	projectId := d.Get("project_id").(string)
	request := templateProjectAssignmentPayloadFromParameters(d)
	result, err := apiClient.AssignTemplateToProject(templateId, request)
	if err != nil {
		return diag.Errorf("could not assign template to project: %v", err)
	}
	resourceId := result.Id + "|" + projectId
	d.SetId(resourceId)
	return nil
}

func resourceTemplateProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)
	templateId := d.Get("template_id").(string)
	template, err := apiClient.Template(templateId)
	if err != nil {
		return diag.Errorf("could not get template: %v", err)
	}
	var assignProjectId = d.Get("project_id").(string)
	d.Set("project_id", "")
	
	for _, projectId := range template.ProjectIds {
		if assignProjectId == projectId {
			d.Set("project_id", projectId)
			
		}
	}
	return nil
}

func resourceTemplateProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	templateId := d.Get("template_id").(string)
	projectId := d.Get("project_id").(string)
	err := apiClient.RemoveTemplateFromProject(templateId, projectId)
	if err != nil {
		return diag.Errorf("could not delete template from project: %v", err)
	}
	return nil
}
