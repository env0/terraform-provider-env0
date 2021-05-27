package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTemplateProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateProjectAssignmenetCreate,
		ReadContext:   resourceTemplateProjectAssignmentRead,
		UpdateContext: resourceTemplateProjectAssignmentUpdate,
		DeleteContext: resourceTemplateProjectAssignmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceTemplateImport},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Description: "id of the template",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Computed:    true,
			},
		},
	}
}

func templateProjectAssignmentPayloadFromParameters(d *schema.ResourceData) (client.TemplateCreatePayload) {
	result := client.TemplateCreatePayload{
		projectId:       d.Get("project_id").(string)
	}

	return result
}

func resourceTemplateProjectAssignmenetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	request := templateProjectAssignmentPayloadFromParameters(d)
	
	template, err := apiClient.AssignTemplateToProject(d.templateId, request)
	if err != nil {
		return diag.Errorf("could not assign template to project: %v", err)
	}

	return nil
}

func resourceTemplateProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	template, err := apiClient.Template(d.Id())
	if err != nil {
		return diag.Errorf("could not get template: %v", err)
	}
	var assignProjectId = d.Get("project_id").(string)
	d.set("project_id", "")

	for _, projectId := range templat.projectIds.([]interface{}) {
		if assignProjectId = projectId {
			d.Set("project_id", projectId)
		}
	}
	return nil
}

func resourceTemplateProjectAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*apiClient := meta.(*client.ApiClient)

	request, problem := templateCreatePayloadFromParameters(d)
	if problem != nil {
		return problem
	}
	_, err := apiClient.TemplateUpdate(d.Id(), request)
	if err != nil {
		return diag.Errorf("could not update template: %v", err)
	}

	return nil*/
}

func resourceTemplateProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*apiClient := meta.(*client.ApiClient)

	id := d.Id()
	err := apiClient.TemplateDelete(id)
	if err != nil {
		return diag.Errorf("could not delete template: %v", err)
	}
	return nil*/
}


