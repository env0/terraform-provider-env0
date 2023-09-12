package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Role", func() {
	roleName := "role_test"
	roleId := uuid.New().String()
	rolePermissions := []string{"VIEW_ORGANIZATION", "EDIT_ORGANIZATION_SETTINGS"}
	roleIsDefaultRole := true

	updatedRoleIsDefaultRole := false
	updatedRolePermissions := []string{"VIEW_ORGANIZATION", "CREATE_AND_EDIT_TEMPLATES"}

	var role *Role
	mockRole := Role{
		Id:             roleId,
		Name:           roleName,
		OrganizationId: organizationId,
		Permissions:    rolePermissions,
		IsDefaultRole:  roleIsDefaultRole,
	}

	updatedMockRole := Role{
		Id:             roleId,
		Name:           roleName,
		OrganizationId: organizationId,
		Permissions:    updatedRolePermissions,
		IsDefaultRole:  updatedRoleIsDefaultRole,
	}

	Describe("RoleCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Post("/roles", RoleCreatePayload{
					Name:           roleName,
					OrganizationId: organizationId,
					Permissions:    rolePermissions,
					IsDefaultRole:  roleIsDefaultRole,
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *Role) {
					*response = mockRole
				})

			role, _ = apiClient.RoleCreate(RoleCreatePayload{
				Name:          roleName,
				Permissions:   rolePermissions,
				IsDefaultRole: roleIsDefaultRole,
			})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return role", func() {
			Expect(*role).To(Equal(mockRole))
		})
	})

	Describe("RoleDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/"+mockRole.Id, nil)
			apiClient.RoleDelete(mockRole.Id)
		})

		It("Should send DELETE request with role id", func() {
			httpCall.Times(1)
		})
	})

	Describe("RoleUpdate", func() {
		BeforeEach(func() {
			payload := RoleUpdatePayload{
				Name:          updatedMockRole.Name,
				Permissions:   updatedMockRole.Permissions,
				IsDefaultRole: updatedMockRole.IsDefaultRole,
			}

			httpCall = mockHttpClient.EXPECT().
				Put("/roles/"+updatedMockRole.Id, payload, gomock.Any()).
				Do(func(path string, request interface{}, response *Role) {
					*response = updatedMockRole
				})
			role, _ = apiClient.RoleUpdate(updatedMockRole.Id, payload)
		})

		It("Should send PUT request with role ID and expected payload", func() {
			httpCall.Times(1)
		})

		It("Should return role received from API", func() {
			Expect(*role).To(Equal(updatedMockRole))
		})
	})

	Describe("Role", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/roles/"+mockRole.Id, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *Role) {
					*response = mockRole
				})
			role, _ = apiClient.Role(mockRole.Id)
		})

		It("Should send GET request with role id", func() {
			httpCall.Times(1)
		})

		It("Should return role", func() {
			Expect(*role).To(Equal(mockRole))
		})
	})

	Describe("Roles", func() {
		var roles []Role
		mockRoles := []Role{mockRole}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Get("/roles", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]Role) {
					*response = mockRoles
				})
			roles, _ = apiClient.Roles()
		})

		It("Should send GET request with organization id param", func() {
			httpCall.Times(1)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should return roles", func() {
			Expect(roles).To(Equal(mockRoles))
		})
	})
})
