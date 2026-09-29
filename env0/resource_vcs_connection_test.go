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
	t.Parallel()

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

	t.Run("Azure DevOps", func(t *testing.T) {
		tokenId := uuid.NewString()

		adoConnection := client.VcsConnection{
			Id:            uuid.NewString(),
			Name:          "ado-connection",
			Type:          "AzureDevOps",
			OauthProvider: "entra",
			TokenId:       tokenId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":     adoConnection.Name,
						"type":     adoConnection.Type,
						"token_id": tokenId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", adoConnection.Id),
						resource.TestCheckResourceAttr(accessor, "token_id", tokenId),
						resource.TestCheckResourceAttr(accessor, "oauth_provider", "entra"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":     "renamed-ado-connection",
						"type":     adoConnection.Type,
						"token_id": tokenId,
					}),
					ExpectError: regexp.MustCompile("an AzureDevOps VCS connection cannot be renamed"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().VcsConnectionCreate(client.VcsConnectionCreatePayload{
				Name:    adoConnection.Name,
				Type:    adoConnection.Type,
				TokenId: tokenId,
			}).Times(1).Return(&adoConnection, nil)
			mock.EXPECT().VcsConnection(adoConnection.Id).AnyTimes().Return(&adoConnection, nil)
			mock.EXPECT().VcsConnectionDelete(adoConnection.Id).Times(1)
		})
	})

	t.Run("Azure DevOps Import", func(t *testing.T) {
		adoConnection := client.VcsConnection{
			Id:   uuid.NewString(),
			Type: "AzureDevOps",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"type":     adoConnection.Type,
						"token_id": uuid.NewString(),
					}),
					ResourceName:       resourceNameImport,
					ImportState:        true,
					ImportStateId:      adoConnection.Id,
					ImportStatePersist: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"type":     adoConnection.Type,
						"token_id": uuid.NewString(),
					}),
					PlanOnly: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().VcsConnection(adoConnection.Id).AnyTimes().Return(&adoConnection, nil)
			mock.EXPECT().VcsConnectionDelete(adoConnection.Id).Times(1)
		})
	})

	for vcsType, expectedError := range map[string]string{
		"AzureDevOps":      "token_id is required for AzureDevOps",
		"GitHubEnterprise": "url is required for GitHubEnterprise",
	} {
		t.Run("Missing Required Field "+vcsType, func(t *testing.T) {
			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
							"name": "test",
							"type": vcsType,
						}),
						ExpectError: regexp.MustCompile(expectedError),
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
		})
	}
}
