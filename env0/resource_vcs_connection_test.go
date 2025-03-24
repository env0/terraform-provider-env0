package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitVcsConnectionResource(t *testing.T) {
	resourceType := "env0_vcs_connection"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	vcsConnection := client.VcsConnection{
		Id:          uuid.NewString(),
		Name:        "test-connection",
		Type:        "GitHubEnterprise",
		Url:         "https://github.example.com",
		VcsAgentKey: "ENV0_DEFAULT",
	}

	updatedVcsConnection := client.VcsConnection{
		Id:          vcsConnection.Id,
		Name:        "updated-connection",
		Type:        vcsConnection.Type,
		Url:         vcsConnection.Url,
		VcsAgentKey: "custom-agent",
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          vcsConnection.Name,
						"type":          vcsConnection.Type,
						"url":           vcsConnection.Url,
						"vcs_agent_key": vcsConnection.VcsAgentKey,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsConnection.Id),
						resource.TestCheckResourceAttr(accessor, "name", vcsConnection.Name),
						resource.TestCheckResourceAttr(accessor, "type", vcsConnection.Type),
						resource.TestCheckResourceAttr(accessor, "url", vcsConnection.Url),
						resource.TestCheckResourceAttr(accessor, "vcs_agent_key", vcsConnection.VcsAgentKey),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          updatedVcsConnection.Name,
						"type":          updatedVcsConnection.Type,
						"url":           updatedVcsConnection.Url,
						"vcs_agent_key": updatedVcsConnection.VcsAgentKey,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedVcsConnection.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedVcsConnection.Name),
						resource.TestCheckResourceAttr(accessor, "type", updatedVcsConnection.Type),
						resource.TestCheckResourceAttr(accessor, "url", updatedVcsConnection.Url),
						resource.TestCheckResourceAttr(accessor, "vcs_agent_key", updatedVcsConnection.VcsAgentKey),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().VcsConnectionCreate(gomock.Any()).Times(1).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnection(vcsConnection.Id).Times(2).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnectionUpdate(vcsConnection.Id, gomock.Any()).Times(1).Return(&updatedVcsConnection, nil),
				mock.EXPECT().VcsConnection(updatedVcsConnection.Id).Times(2).Return(&updatedVcsConnection, nil),
				mock.EXPECT().VcsConnectionDelete(updatedVcsConnection.Id).Times(1),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          vcsConnection.Name,
						"type":          vcsConnection.Type,
						"url":           vcsConnection.Url,
						"vcs_agent_key": vcsConnection.VcsAgentKey,
					}),
					ExpectError: regexp.MustCompile("could not create VCS connection: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().Times(1).Return("org-1", nil)
			mock.EXPECT().VcsConnectionCreate(gomock.Any()).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          vcsConnection.Name,
						"type":          vcsConnection.Type,
						"url":           vcsConnection.Url,
						"vcs_agent_key": vcsConnection.VcsAgentKey,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     vcsConnection.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().VcsConnectionCreate(gomock.Any()).Times(1).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnection(vcsConnection.Id).Times(3).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnectionDelete(vcsConnection.Id).Times(1),
			)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          vcsConnection.Name,
						"type":          vcsConnection.Type,
						"url":           vcsConnection.Url,
						"vcs_agent_key": vcsConnection.VcsAgentKey,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     vcsConnection.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().VcsConnectionCreate(gomock.Any()).Times(1).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnection(vcsConnection.Id).Times(1).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnections().Times(1).Return([]client.VcsConnection{vcsConnection}, nil),
				mock.EXPECT().VcsConnection(vcsConnection.Id).Times(1).Return(&vcsConnection, nil),
				mock.EXPECT().VcsConnectionDelete(vcsConnection.Id).Times(1),
			)
		})
	})

	t.Run("Import By Name - Multiple Found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          vcsConnection.Name,
						"type":          vcsConnection.Type,
						"url":           vcsConnection.Url,
						"vcs_agent_key": vcsConnection.VcsAgentKey,
					}),
				},
				{
					ResourceName:  resourceNameImport,
					ImportState:   true,
					ImportStateId: vcsConnection.Name,
					ExpectError:   regexp.MustCompile("found multiple VCS connections with name: .* Use id instead or make sure VCS connection names are unique"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			duplicateConnection := vcsConnection
			duplicateConnection.Id = "different-id"

			mock.EXPECT().VcsConnectionCreate(gomock.Any()).Times(1).Return(&vcsConnection, nil)
			mock.EXPECT().VcsConnection(vcsConnection.Id).Times(1).Return(&vcsConnection, nil)
			mock.EXPECT().VcsConnections().Times(1).Return(
				[]client.VcsConnection{vcsConnection, duplicateConnection},
				nil,
			)
			mock.EXPECT().VcsConnectionDelete(vcsConnection.Id).Times(1)
		})
	})

	t.Run("Import By Name - Not Found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":          vcsConnection.Name,
						"type":          vcsConnection.Type,
						"url":           vcsConnection.Url,
						"vcs_agent_key": vcsConnection.VcsAgentKey,
					}),
				},
				{
					ResourceName:  resourceNameImport,
					ImportState:   true,
					ImportStateId: "non-existent-name",
					ExpectError:   regexp.MustCompile("VCS connection with name .* not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().VcsConnectionCreate(gomock.Any()).Times(1).Return(&vcsConnection, nil)
			mock.EXPECT().VcsConnection(vcsConnection.Id).Times(1).Return(&vcsConnection, nil)
			mock.EXPECT().VcsConnections().Times(1).Return([]client.VcsConnection{}, nil)
			mock.EXPECT().VcsConnectionDelete(vcsConnection.Id).Times(1)
		})
	})
}
