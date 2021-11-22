package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "the environment's id",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "the environments name",
				Optional:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the environment's project id",
				Required:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the environment's template id",
				Computed:    true,
			},
			//"workspace": {
			//	Type:        schema.TypeString,
			//	Description: "The environments workspace",
			//	Computed:    true,
			//},
			//"revision": {
			//	Type:        schema.TypeString,
			//	Description: "The environments revision",
			//	Computed:    true,
			//},
			//"redeploy_on_git_push": {
			//	Type:        schema.TypeBool,
			//	Description: "Indicate if git push events should trigger deployments",
			//	Computed:    true,
			//},
			//"run_plan_on_pull_request": {
			//	Type:        schema.TypeBool,
			//	Description: "Indicate if a pull request creation should trigger a plan",
			//	Computed:    true,
			//},
			//"approve_plan_automatically": {
			//	Type:        schema.TypeBool,
			//	Description: "Indicate if the environment requires approval after plan",
			//	Computed:    true,
			//},
			//"force_destroy": {
			//	Type:        schema.TypeBool,
			//	// TODO: update this description
			//	Description: "idk",
			//	Computed:    true,
			//},
			//"configuration": {
			//	Type: schema.TypeList,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeSet,
			//		Elem: &schema.Resource{
			//			Schema: map[string]*schema.Schema{
			//				"name": {
			//					Type:     schema.TypeString,
			//					Required: true,
			//				},
			//				"type": {
			//					Type:     schema.TypeString,
			//					Description: "variable type, either environment or terraform (defaults to environment)",
			//					Optional: true,
			//					Default:  "environment",
			//				},
			//				"value": {
			//					Type:     schema.TypeString,
			//					Required: true,
			//				},
			//			},
			//		},
			//	},
			//},
		},
	}
}

func dataEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var environment client.Environment

	environmentId, ok := d.GetOk("id")
	if ok {
		environment, err = getEnvironment(environmentId.(string), meta)
		if err != nil {
			return err
		}
	}
	//d.SetId(policy.Id)
	setEnvironmentSchema(d, environment)
	return nil
}

func getEnvironment(environmentId string, meta interface{}) (client.Environment, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	environment, err := apiClient.Environment(environmentId)
	if err != nil {
		return client.Environment{}, diag.Errorf("Could not find environment: %v", err)
	}
	return environment, nil
}
