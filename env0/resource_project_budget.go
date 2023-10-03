package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectBudget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectBudgetCreateOrUpdate,
		UpdateContext: resourceProjectBudgetCreateOrUpdate,
		ReadContext:   resourceProjectBudgetRead,
		DeleteContext: resourceProjectBudgetDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
				ForceNew:    true,
			},
			"amount": {
				Type:        schema.TypeInt,
				Description: "amount of the project budget",
				Required:    true,
			},
			"timeframe": {
				Type:             schema.TypeString,
				Description:      "budget timeframe (valid values: WEEKLY, MONTHLY, QUARTERLY, YEARLY)",
				Required:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"WEEKLY", "MONTHLY", "QUARTERLY", "YEARLY"}),
			},
			"thresholds": {
				Type:        schema.TypeList,
				Description: "list of notification thresholds",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeInt,
					Description: "a threshold in %",
				},
			},
		},
	}
}

func resourceProjectBudgetCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var payload client.ProjectBudgetUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	projectId := d.Get("project_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	budget, err := apiClient.ProjectBudgetUpdate(projectId, &payload)
	if err != nil {
		return diag.Errorf("could not create or update budget: %v", err)
	}

	d.SetId(budget.Id)

	return nil
}

func resourceProjectBudgetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)

	budget, err := apiClient.ProjectBudget(projectId)
	if err != nil {
		return ResourceGetFailure(ctx, "project budget", d, err)
	}

	if err := writeResourceData(budget, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceProjectBudgetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ProjectBudgetDelete(projectId); err != nil {
		return diag.Errorf("could not delete project budget: %v", err)
	}

	return nil
}
