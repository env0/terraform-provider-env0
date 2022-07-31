package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestIpRangesDataSource(t *testing.T) {
	resourceType := "env0_ip_ranges"
	resourceName := "ip_ranges"
	accessor := dataSourceAccessor(resourceType, resourceName)

	testCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "ipv4.0", "18.214.210.123/32"),
						resource.TestCheckResourceAttr(accessor, "ipv4.1", "44.195.170.230/32"),
					),
				},
			},
		}
	}

	t.Run("Get IP Ranges", func(t *testing.T) {
		runUnitTest(t,
			testCase(),
			func(mockFunc *client.MockApiClientInterface) {},
		)
	})
}
