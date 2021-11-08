package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitPolicyResource(t *testing.T) {
	resourceType := "env0_policy"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	policy := client.Policy{
		Id: "id0",
	}

	updatedPolicy := client.Policy{
		Id: policy.Id,
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"id": policy.Id,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", policy.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"id": updatedPolicy.Id,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", updatedPolicy.Id),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
			NumberOfEnvironments: policy.NumberOfEnvironments,
		}).Times(1).Return(policy, nil)
		mock.EXPECT().PolicyUpdate(client.PolicyUpdatePayload{
			NumberOfEnvironments: updatedPolicy.NumberOfEnvironments,
		}).Times(1).Return(updatedPolicy, nil)

		gomock.InOrder(
			mock.EXPECT().Policy(gomock.Any()).Times(2).Return(policy, nil),        // 1 after create, 1 before update
			mock.EXPECT().Policy(gomock.Any()).Times(1).Return(updatedPolicy, nil), // 1 after update
		)
	})
}

func TestUnitPolicyInvalidParams(t *testing.T) {
	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate("env0_policy", "test", map[string]interface{}{"project_id": ""}),
				ExpectError: regexp.MustCompile("Project id cannot be empty"),
			},
		},
	}

	runUnitTest(t, testCase, func(mockFunc *client.MockApiClientInterface) {})
}
