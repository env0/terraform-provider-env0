package env0

import (
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

	testProjectDataSource := func(input map[string]string, mockFunc func(mockFunc *client.MockApiClientInterface)) {
		resourceType := "env0_project"
		resourceName := "test_project"
		accessor := dataSourceAccessor(resourceType, resourceName)

		testCase := resource.TestCase{
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

		runUnitTest(t, testCase, mockFunc)
	}

	t.Run("By ID", func(t *testing.T) {
		testProjectDataSource(
			map[string]string{"id": project.Id},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Project(project.Id).AnyTimes().Return(project, nil)
			},
		)
	})

	t.Run("By Name", func(t *testing.T) {
		testProjectDataSource(
			map[string]string{"name": project.Name},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Projects().AnyTimes().Return([]client.Project{project}, nil)
			},
		)
	})
}
