package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var project = client.Project{
	Id:          "id0",
	Name:        "my-project-1",
	CreatedBy:   "env0",
	Role:        "role0",
	Description: "A project's description",
}

func testProjectDataSource(t *testing.T, input map[string]string, mockFunc func(mockFunc *client.MockApiClientInterface)) {
	resourceType := "env0_project"
	resourceName := "test_project"
	resourceFullName := dataSourceAccessor(resourceType, resourceName)

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, input),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", project.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", project.Name),
					resource.TestCheckResourceAttr(resourceFullName, "created_by", project.CreatedBy),
					resource.TestCheckResourceAttr(resourceFullName, "role", project.Role),
					resource.TestCheckResourceAttr(resourceFullName, "description", project.Description),
				),
			},
		},
	}

	runUnitTest(t, testCase, mockFunc)
}

func TestUnitProjectDataByName(t *testing.T) {
	testProjectDataSource(
		t,
		map[string]string{"name": project.Name},
		func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Projects().AnyTimes().Return([]client.Project{project}, nil)
		},
	)

}

func TestUnitProjectDataById(t *testing.T) {
	testProjectDataSource(
		t,
		map[string]string{"id": project.Id},
		func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Project(project.Id).AnyTimes().Return(project, nil)
		},
	)
}
