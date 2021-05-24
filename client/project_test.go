package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const projectName = "project_test"

var _ = Describe("Project", func() {
	var project Project
	mockProject := &Project{
		Id:             "id",
		Name:           "config-key",
		OrganizationId: organizationId,
	}

	Describe("Create", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().Post(
				"/projects",
				map[string]interface{}{
					"name":           projectName,
					"organizationId": organizationId,
				},
				&project,
			)

			project, _ = apiClient.ProjectCreate(projectName)
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return project variable", func() {
			Expect(project).To(Equal(mockProject))
		})
	})
})
