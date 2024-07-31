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
		"3.209.36.240/32",
		"3.222.51.117/32",
		"3.226.24.146/32",
		"18.214.49.142/32",
		"18.214.210.123/32",
		"35.81.146.242/32",
		"35.85.88.233/32",
		"44.195.170.230/32",
		"44.205.134.220/32",
		"44.212.144.113/32",
		"44.227.16.37/32",
		"44.228.227.2/32",
		"44.240.181.100/32",
		"52.73.227.111/32",
		"54.68.137.240/32",
		"54.88.50.2/32",
		"54.149.16.114/32",
		"54.165.19.49/32",
	}

	if err := d.Set("ipv4", ipv4s); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ip_ranges")

	return nil
}
