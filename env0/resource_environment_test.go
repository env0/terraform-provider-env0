package env0

import (
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestUnitEnvironmentResource(t *testing.T) {
	resourceType := "env0_environment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	environment := client.Environment{
		Id:         "id0",
		Name:       "my-environment",
		ProjectId:  "project-id",
		TemplateId: "template-id",
	}

	updatedEnvironment := client.Environment{
		Id:         environment.Id,
		Name:       "my-updated-environment-name",
		ProjectId:  "project-id",
		TemplateId: "template-id",
	}

	createEnvironmentResourceConfig := func(environment client.Environment) string {
		return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"name":        environment.Name,
			"project_id":  environment.ProjectId,
			"template_id": environment.TemplateId,
		})
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createEnvironmentResourceConfig(environment),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environment.Id),
						resource.TestCheckResourceAttr(accessor, "name", environment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", environment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", environment.TemplateId),
					),
				},
				{
					Config: createEnvironmentResourceConfig(updatedEnvironment),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironment.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironment.Name),
						resource.TestCheckResourceAttr(accessor, "project_id", updatedEnvironment.ProjectId),
						resource.TestCheckResourceAttr(accessor, "template_id", updatedEnvironment.TemplateId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: environment.TemplateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
				Name: updatedEnvironment.Name,
			}).Times(1).Return(updatedEnvironment, nil)

			gomock.InOrder(
				mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil),        // 1 after create, 1 before update
				mock.EXPECT().Environment(gomock.Any()).Times(1).Return(updatedEnvironment, nil), // 1 after update
			)

			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})
	})

	t.Run("Failure in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createEnvironmentResourceConfig(environment),
					ExpectError: regexp.MustCompile("could not create environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: environment.TemplateId,
				},
			}).Times(1).Return(client.Environment{}, errors.New("error"))
		})

	})

	t.Run("Failure in update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: createEnvironmentResourceConfig(environment),
				},
				{
					Config:      createEnvironmentResourceConfig(updatedEnvironment),
					ExpectError: regexp.MustCompile("could not update environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: environment.TemplateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().EnvironmentUpdate(updatedEnvironment.Id, client.EnvironmentUpdate{
				Name: updatedEnvironment.Name,
			}).Times(1).Return(client.Environment{}, errors.New("error"))
			mock.EXPECT().Environment(gomock.Any()).Times(2).Return(environment, nil) // 1 after create, 1 before update
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})

	})

	t.Run("Failure in read", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      createEnvironmentResourceConfig(environment),
					ExpectError: regexp.MustCompile("could not get environment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentCreate(client.EnvironmentCreate{
				Name:      environment.Name,
				ProjectId: environment.ProjectId,
				DeployRequest: &client.DeployRequest{
					BlueprintId: environment.TemplateId,
				},
			}).Times(1).Return(environment, nil)
			mock.EXPECT().Environment(gomock.Any()).Return(client.Environment{}, errors.New("error"))
			mock.EXPECT().EnvironmentDestroy(environment.Id).Times(1)
		})

	})
}
