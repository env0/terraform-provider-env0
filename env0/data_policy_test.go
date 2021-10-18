package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestPolicyDataSource(t *testing.T) {
	policy := client.Policy{
		Id:        "id0",
		ProjectId: "project0",
	}

	resourceType := "env0_policy"
	resourceName := "test_policy"
	accessor := dataSourceAccessor(resourceType, resourceName)

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, make(map[string]interface{})),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "project_id", policy.ProjectId),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().Policy(policy.Id).AnyTimes().Return(policy, nil)
	})
}
