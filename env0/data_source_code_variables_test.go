package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSourceCodeVariablesDataSource(t *testing.T) {
	template := client.Template{
		Id:               "id0",
		Name:             "template0",
		Description:      "description0",
		Repository:       "env0/repo",
		Path:             "path/zero",
		Revision:         "branch-zero",
		Type:             "terraform",
		TokenId:          "1",
		TerraformVersion: "0.12.24",
	}

	payload := &client.VariablesFromRepositoryPayload{
		BitbucketClientKey:   template.BitbucketClientKey,
		GithubInstallationId: template.GithubInstallationId,
		Path:                 template.Path,
		Revision:             template.Revision,
		TokenId:              template.TokenId,
		Repository:           template.Repository,
	}

	var1 := client.ConfigurationVariable{
		Id:             "id0",
		Name:           "name0",
		Description:    "desc0",
		ScopeId:        "scope0",
		Value:          "value0",
		OrganizationId: "organization0",
		UserId:         "user0",
		Schema:         &client.ConfigurationVariableSchema{Type: "string", Format: client.HCL},
		Regex:          "regex",
	}

	var2 := client.ConfigurationVariable{
		Id:             "id1",
		Name:           "name1",
		Description:    "desc1",
		ScopeId:        "scope1",
		Value:          "value1",
		OrganizationId: "organization0",
		UserId:         "user1",
		Schema:         &client.ConfigurationVariableSchema{Type: "string", Format: client.JSON},
		Regex:          "regex",
	}

	vars := []client.ConfigurationVariable{var1, var2}

	resourceType := "env0_source_code_variables"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"template_id": template.Id}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "variables.0.name", var1.Name),
					resource.TestCheckResourceAttr(accessor, "variables.1.name", var2.Name),
					resource.TestCheckResourceAttr(accessor, "variables.0.value", var1.Value),
					resource.TestCheckResourceAttr(accessor, "variables.1.value", var2.Value),
					resource.TestCheckResourceAttr(accessor, "variables.0.format", string(var1.Schema.Format)),
					resource.TestCheckResourceAttr(accessor, "variables.1.format", string(var2.Schema.Format)),
					resource.TestCheckResourceAttr(accessor, "variables.0.description", var1.Description),
					resource.TestCheckResourceAttr(accessor, "variables.1.description", var2.Description),
				),
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Template(template.Id).AnyTimes().Return(template, nil)
			mock.EXPECT().VariablesFromRepository(payload).AnyTimes().Return(vars, nil)
		})
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{"template_id": template.Id}),
						ExpectError: regexp.MustCompile("error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Template(template.Id).AnyTimes().Return(template, nil)
				mock.EXPECT().VariablesFromRepository(payload).AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})

}
