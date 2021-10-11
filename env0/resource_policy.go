package env0

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: nil,
		ReadContext:   nil,
		UpdateContext: nil,
		DeleteContext: nil,

		Importer: nil,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "id  of the project",
				Required:    true,
				ForceNew:    true,
			},
			"number_of_environments": {
				Type:        schema.TypeInt,
				Description: "number of environments",
				Optional:    true,
			},
			"requires_approval_default": {
				Type:        schema.TypeBool,
				Description: "requires approval",
				Optional:    true,
			},
			"include_cost_estimation": {
				Type:        schema.TypeBool,
				Description: "include cost estimation",
				Optional:    true,
			},
			"skip_apply_when_plan_is_empty": {
				Type:        schema.TypeBool,
				Description: "skip apply when plan is empty",
				Optional:    true,
			},
			"disable_destroy_environments": {
				Type:        schema.TypeBool,
				Description: "disable destroy environments",
				Optional:    true,
			},
		},
	}
}
