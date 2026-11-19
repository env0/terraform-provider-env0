package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestApprovalPolicyDataSource(t *testing.T) {
	t.Parallel()

	approvalPolicy := client.ApprovalPolicy{
		Id:         "id0",
		Name:       "name0",
		Repository: "rep1",
		Path:       "path",
	}

	otherApprovalPolicy := client.ApprovalPolicy{
		Id:         "id1",
		Name:       "name0",
		Repository: "rep2",
		Path:       "path",
	}

	approvalPolicyTemplate := client.Template{
		Id:         approvalPolicy.Id,
		Name:       approvalPolicy.Name,
		Repository: approvalPolicy.Repository,
		Path:       approvalPolicy.Path,
		Type:       string(ApprovalPolicy),
	}

	approvalPolicyFieldsByName := map[string]any{"name": approvalPolicy.Name}
	approvalPolicyFieldsById := map[string]any{"id": approvalPolicy.Id}

	resourceType := "env0_approval_policy"
	resourceName := "test_approval_policy"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", approvalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", approvalPolicy.Name),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockApprovalPoliciesCall := func(returnValue []client.ApprovalPolicy) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApprovalPolicies(approvalPolicy.Name).AnyTimes().Return(returnValue, nil)
		}
	}

	mockTemplateCall := func(returnValue client.Template) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(approvalPolicy.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(approvalPolicyFieldsById),
			mockTemplateCall(approvalPolicyTemplate),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(approvalPolicyFieldsByName),
			mockApprovalPoliciesCall([]client.ApprovalPolicy{approvalPolicy}),
		)
	})

	t.Run("Throw error when by name and more than one approval policy exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(approvalPolicyFieldsByName, "found multiple approval policies with name"),
			mockApprovalPoliciesCall([]client.ApprovalPolicy{approvalPolicy, otherApprovalPolicy}),
		)
	})

	t.Run("Throw error when by name and no approval policy found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(approvalPolicyFieldsByName, "approval policy with name name0 not found"),
			mockApprovalPoliciesCall([]client.ApprovalPolicy{}),
		)
	})

	t.Run("Throw error when by id and no approval policy found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(approvalPolicyFieldsById, "failed to get approval policy by id"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(client.Template{}, http.NewMockFailedResponseError(404))
			},
		)
	})

	t.Run("Throw error when by id and the template is not an approval policy", func(t *testing.T) {
		customFlowTemplate := approvalPolicyTemplate
		customFlowTemplate.Type = string(CustomFlow)

		runUnitTest(t,
			getErrorTestCase(approvalPolicyFieldsById, "template type requires type approval-policy but received type custom-flow"),
			mockTemplateCall(customFlowTemplate),
		)
	})
}
