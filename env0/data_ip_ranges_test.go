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
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "ipv4.0", "3.209.36.240/32"),
						resource.TestCheckResourceAttr(accessor, "ipv4.1", "3.222.51.117/32"),
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
