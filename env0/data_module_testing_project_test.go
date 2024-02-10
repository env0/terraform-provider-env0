package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestModuleTestingProjectDataSource(t *testing.T) {
	resourceType := "env0_module_testing_project"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	moduleTestingProject := client.ModuleTestingProject{
		Name:            "namex",
		ParentProjectId: "pidx",
		Id:              "idx",
	}

	getTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", moduleTestingProject.Name),
						resource.TestCheckResourceAttr(accessor, "id", moduleTestingProject.Id),
						resource.TestCheckResourceAttr(accessor, "parent_project_id", moduleTestingProject.ParentProjectId),
					),
				},
			},
		}
	}

	mockModuleTestingProject := func() func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleTestingProject().AnyTimes().Return(&moduleTestingProject, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(),
			mockModuleTestingProject(),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
						ExpectError: regexp.MustCompile("could not get module testing project: error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ModuleTestingProject().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})

}
