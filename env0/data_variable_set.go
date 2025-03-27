package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataVariableSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataVariableSetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the name of the variable set",
				Required:    true,
			},
			"scope": {
				Type:             schema.TypeString,
				Description:      "the scope of the variable set. Valid values: 'ORGANIZATION', or 'PROJECT'",
				Required:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"ORGANIZATION", "PROJECT"}),
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the id of the 'PROJECT' scope. Is not required for 'ORGANIZATION' scope",
				Optional:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "the id variable set",
				Computed:    true,
			},
		},
	}
}

func dataVariableSetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	resource := struct {
		Name      string
		Scope     string
		ProjectId string
	}{}

	if err := readResourceData(&resource, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	var scopeId string

	switch resource.Scope {
	case "ORGANIZATION":
		var err error

		scopeId, err = apiClient.OrganizationId()
		if err != nil {
			return diag.Errorf("could not get organization id: %v", err)
		}
	case "PROJECT":
		if resource.ProjectId == "" {
			return diag.Errorf("'project_id' is required")
		}

		scopeId = resource.ProjectId
	}

	variableSets, err := apiClient.ConfigurationSets(resource.Scope, scopeId)
	if err != nil {
		return diag.Errorf("could not get variable sets: %v", err)
	}

	for _, variableSet := range variableSets {
		if variableSet.Name == resource.Name {
			d.SetId(variableSet.Id)

			return nil
		}
	}

	return diag.Errorf("variable set not found")
}
