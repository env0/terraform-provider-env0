package env0

import (
	"context"
	"github.com/adhocore/gronx"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironmentScheduling() *schema.Resource {

	validateCronExpression := func(i interface{}, path cty.Path) diag.Diagnostics {
		expr := i.(string)
		parser := gronx.New()
		isValid := parser.IsValid(expr)

		if isValid != true {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid cron expression",
					AttributePath: path,
				}}
		}

		return nil
	}

	return &schema.Resource{
		CreateContext: resourceEnvironmentSchedulingCreateOrUpdate,
		ReadContext:   resourceEnvironmentSchedulingRead,
		UpdateContext: resourceEnvironmentSchedulingCreateOrUpdate,
		DeleteContext: resourceEnvironmentSchedulingDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Description: "the environment's id",
				Required:    true,
				ForceNew:    true,
			},
			"destroy_cron": {
				Type:             schema.TypeString,
				Description:      "Cron expression for scheduled destroy of the environment",
				AtLeastOneOf:     []string{"destroy_cron", "deploy_cron"},
				Optional:         true,
				ValidateDiagFunc: validateCronExpression,
			},
			"deploy_cron": {
				Type:             schema.TypeString,
				Description:      "Cron expression for scheduled deploy of the environment",
				AtLeastOneOf:     []string{"destroy_cron", "deploy_cron"},
				Optional:         true,
				ValidateDiagFunc: validateCronExpression,
			},
		},
	}
}

func resourceEnvironmentSchedulingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Id()

	environmentScheduling, err := apiClient.EnvironmentScheduling(environmentId)

	if err != nil {
		return diag.Errorf("could not get environment scheduling: %v", err)
	}

	d.Set("deploy_cron", environmentScheduling.Deploy.Cron)
	d.Set("destroy_cron", environmentScheduling.Destroy.Cron)
	return nil
}

func resourceEnvironmentSchedulingCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Get("environment_id").(string)
	deployCron := d.Get("deploy_cron").(string)
	destroyCron := d.Get("destroy_cron").(string)

	payload := EnvironmentScheduling{
		Deploy:  EnvironmentSchedulingExpression{Cron: deployCron, Enabled: true},
		Destroy: EnvironmentSchedulingExpression{Cron: destroyCron, Enabled: true},
	}

	_, err := apiClient.EnvironmentSchedulingUpdate(environmentId, payload)

	if err != nil {
		return diag.Errorf("could not create or update environment scheduling: %v", err)
	}

	d.SetId(environmentId)
	return nil
}

func resourceEnvironmentSchedulingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(ApiClientInterface)

	environmentId := d.Id()

	err := apiClient.EnvironmentSchedulingDelete(environmentId)

	if err != nil {
		return diag.Errorf("could not delete environment scheduling: %v", err)
	}

	return nil
}
