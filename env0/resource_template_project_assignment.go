package env0

import (
	"context"

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

func resourceTemplateProjectAssignmenetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*apiClient := meta.(*client.ApiClient)

	request, problem := templateCreatePayloadFromParameters(d)
	if problem != nil {
		return problem
	}
	template, err := apiClient.TemplateCreate(request)
	if err != nil {
		return diag.Errorf("could not create template: %v", err)
	}

	d.SetId(template.Id)

	return nil*/
}

func resourceTemplateProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*apiClient := meta.(*client.ApiClient)

	template, err := apiClient.Template(d.Id())
	if err != nil {
		return diag.Errorf("could not get template: %v", err)
	}

	d.Set("name", template.Name)
	d.Set("description", template.Description)
	d.Set("github_installation_id", template.GithubInstallationId)
	d.Set("repository", template.Repository)
	d.Set("path", template.Path)
	d.Set("revision", template.Revision)
	d.Set("type", template.Type)
	d.Set("project_ids", template.ProjectIds)
	d.Set("terraform_version", template.TerraformVersion)
	if template.Retry.OnDeploy != nil {
		d.Set("retries_on_deploy", template.Retry.OnDeploy.Times)
		d.Set("retry_on_deploy_only_when_matches_regex", template.Retry.OnDeploy.ErrorRegex)
	} else {
		d.Set("retries_on_deploy", 0)
		d.Set("retry_on_deploy_only_when_matches_regex", "")
	}
	if template.Retry.OnDestroy != nil {
		d.Set("retries_on_destroy", template.Retry.OnDestroy.Times)
		d.Set("retry_on_destroy_only_when_matches_regex", template.Retry.OnDestroy.ErrorRegex)
	} else {
		d.Set("retries_on_destroy", 0)
		d.Set("retry_on_destroy_only_when_matches_regex", "")
	}

	return nil*/
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


