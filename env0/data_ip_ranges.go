package env0

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataIpRanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIpRangesRead,

		Schema: map[string]*schema.Schema{
			"ipv4": {
				Type:        schema.TypeList,
				Description: "list of env0 ipv4 CIDR addresses. This list can be used to whitelist inconming env0 traffic (E.g.: https://docs.env0.com/docs/templates#on-premises-git-servers-support)",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "ipv4 CIDR address",
				},
			},
		},
	}
}

func dataIpRangesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ipv4s := []string{
		"18.214.210.123/32",
		"44.195.170.230/32",
	}

	d.Set("ipv4", ipv4s)

	d.SetId("ip_ranges")

	return nil
}
