package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitOrganizationResource(t *testing.T) {
	resourceType := "env0_organization"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	organization := client.Organization{
		Id:          "id0",
		Name:        "my-org",
		Description: "organization description",
		PhotoUrl:    "https://example.com/photo.jpg",
	}

	updatedOrganization := client.Organization{
		Id:          organization.Id,
		Name:        "my-updated-org",
		Description: "updated organization description",
		PhotoUrl:    "https://example.com/updated-photo.jpg",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        organization.Name,
						"description": organization.Description,
						"photo_url":   organization.PhotoUrl,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", organization.Id),
						resource.TestCheckResourceAttr(accessor, "name", organization.Name),
						resource.TestCheckResourceAttr(accessor, "description", organization.Description),
						resource.TestCheckResourceAttr(accessor, "photo_url", organization.PhotoUrl),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedOrganization.Name,
						"description": updatedOrganization.Description,
						"photo_url":   updatedOrganization.PhotoUrl,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedOrganization.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedOrganization.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedOrganization.Description),
						resource.TestCheckResourceAttr(accessor, "photo_url", updatedOrganization.PhotoUrl),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationCreate(client.OrganizationCreatePayload{
				Name:        organization.Name,
				Description: organization.Description,
				PhotoUrl:    organization.PhotoUrl,
			}).Times(1).Return(&organization, nil)

			mock.EXPECT().OrganizationUpdate(organization.Id, client.OrganizationUpdatePayload{
				Name:        updatedOrganization.Name,
				Description: updatedOrganization.Description,
				PhotoUrl:    updatedOrganization.PhotoUrl,
			}).Times(1).Return(&updatedOrganization, nil)

			gomock.InOrder(
				mock.EXPECT().OrganizationById(gomock.Any()).Times(2).Return(&organization, nil),
				mock.EXPECT().OrganizationById(gomock.Any()).Times(2).Return(&updatedOrganization, nil),
			)
		})
	})

	t.Run("Failure in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        organization.Name,
						"description": organization.Description,
						"photo_url":   organization.PhotoUrl,
					}),
					ExpectError: regexp.MustCompile("could not create organization: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationCreate(client.OrganizationCreatePayload{
				Name:        organization.Name,
				Description: organization.Description,
				PhotoUrl:    organization.PhotoUrl,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Failure in update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        organization.Name,
						"description": organization.Description,
						"photo_url":   organization.PhotoUrl,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedOrganization.Name,
						"description": updatedOrganization.Description,
						"photo_url":   updatedOrganization.PhotoUrl,
					}),
					ExpectError: regexp.MustCompile("could not update organization: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationCreate(client.OrganizationCreatePayload{
				Name:        organization.Name,
				Description: organization.Description,
				PhotoUrl:    organization.PhotoUrl,
			}).Times(1).Return(&organization, nil)

			mock.EXPECT().OrganizationUpdate(organization.Id, client.OrganizationUpdatePayload{
				Name:        updatedOrganization.Name,
				Description: updatedOrganization.Description,
				PhotoUrl:    updatedOrganization.PhotoUrl,
			}).Times(1).Return(nil, errors.New("error"))

			mock.EXPECT().OrganizationById(gomock.Any()).Times(3).Return(&organization, nil)
		})
	})

	t.Run("Failure in read", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        organization.Name,
						"description": organization.Description,
						"photo_url":   organization.PhotoUrl,
					}),
					ExpectError: regexp.MustCompile("could not get organization: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationCreate(client.OrganizationCreatePayload{
				Name:        organization.Name,
				Description: organization.Description,
				PhotoUrl:    organization.PhotoUrl,
			}).Times(1).Return(&organization, nil)

			mock.EXPECT().OrganizationById(gomock.Any()).Return(nil, errors.New("error"))
		})
	})
}
