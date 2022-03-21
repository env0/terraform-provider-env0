package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Module Client", func() {
	mockModule := Module{
		Id:             "module-id",
		ModuleName:     "module-name",
		ModuleProvider: "module-provider",
		Repository:     "repository-name",
	}

	Describe("Get Single Module", func() {
		var returnedModule *Module

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/modules/"+mockModule.Id, gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *Module) {
					*response = mockModule
				})
			returnedModule, _ = apiClient.Module(mockModule.Id)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return module", func() {
			Expect(*returnedModule).To(Equal(mockModule))
		})
	})

	Describe("Get All Modules", func() {
		var returnedModules []Module
		mockModules := []Module{mockModule}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/modules", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]Module) {
					*response = mockModules
				})
			returnedModules, _ = apiClient.Modules()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return modules", func() {
			Expect(returnedModules).To(Equal(mockModules))
		})
	})

	Describe("Create Module", func() {
		var createdModule *Module
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			createModulePayload := ModuleCreatePayload{}
			copier.Copy(&createModulePayload, &mockModule)

			expectedCreateRequest := ModuleCreatePayloadWith{
				ModuleCreatePayload: createModulePayload,
				OrganizationId:      organizationId,
				Type:                "module",
			}

			httpCall = mockHttpClient.EXPECT().
				Post("/modules", expectedCreateRequest, gomock.Any()).
				Do(func(path string, request interface{}, response *Module) {
					*response = mockModule
				})

			createdModule, err = apiClient.ModuleCreate(createModulePayload)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return created module", func() {
			Expect(*createdModule).To(Equal(mockModule))
		})
	})

	Describe("Delete Module", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/modules/" + mockModule.Id)
			apiClient.ModuleDelete(mockModule.Id)
		})

		It("Should send DELETE request with module id", func() {
			httpCall.Times(1)
		})
	})

	Describe("Update Module", func() {
		var updatedModule *Module
		var err error

		updatedMockModule := mockModule
		updatedMockModule.ModuleName = "updated-module-name"

		BeforeEach(func() {
			updateModulePayload := ModuleUpdatePayload{ModuleName: updatedMockModule.ModuleName}

			httpCall = mockHttpClient.EXPECT().
				Patch("/modules/"+mockModule.Id, updateModulePayload, gomock.Any()).
				Do(func(path string, request interface{}, response *Module) {
					*response = updatedMockModule
				})

			updatedModule, err = apiClient.ModuleUpdate(mockModule.Id, updateModulePayload)
		})

		It("Should send Patch request with expected payload", func() {
			httpCall.Times(1)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return module received from API", func() {
			Expect(*updatedModule).To(Equal(updatedMockModule))
		})
	})
})
