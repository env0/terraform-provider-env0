package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProjectDataSource(t *testing.T) {
	project := client.Project{
		Id:          "id0",
		Name:        "my-project-1",
		CreatedBy:   "env0",
		Role:        "role0",
		Description: "A project's description",
	}

	projectDataByName := map[string]interface{}{"name": project.Name}
	projectDataById := map[string]interface{}{"id": project.Id}

	resourceType := "env0_project"
	resourceName := "test_project"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", project.Id),
						resource.TestCheckResourceAttr(accessor, "name", project.Name),
						resource.TestCheckResourceAttr(accessor, "created_by", project.CreatedBy),
						resource.TestCheckResourceAttr(accessor, "role", project.Role),
						resource.TestCheckResourceAttr(accessor, "description", project.Description),
					),
				},
			},
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(projectDataById),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Project(project.Id).AnyTimes().Return(project, nil)
			})
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(projectDataByName),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Projects().AnyTimes().Return([]client.Project{project}, nil)
			})
	})

	t.Run("Throw error when by name and more than one project exists", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, projectDataByName),
					ExpectError: regexp.MustCompile(`Found multiple Projects for name`),
				},
			},
		}

		runUnitTest(t,
			testCase,
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Projects().AnyTimes().Return([]client.Project{project, project}, nil)
			})
	})
}
