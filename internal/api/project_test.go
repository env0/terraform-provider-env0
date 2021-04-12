package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/env0/terraform-provider-env0/internal/api"
)

var _ = Describe("Fetching projects list", func() {
	var projects []Project
	var projectsErr error

	JustBeforeEach(func() {
		projects, projectsErr = apiClient.Projects()
	})

	When("Looking for the Default Organization Project", func() {
		var defaultProject *Project

		JustBeforeEach(func() {
			defaultProject = nil
			for _, project := range projects {
				if project.Name == "Default Organization Project" {
					defaultProject = &project
				}
			}
		})

		Specify("fetching the projects list should not fail", func() {
			Expect(projectsErr).To(BeNil())
		})
		Specify("there should be at least one project", func() {
			Expect(len(projects)).ToNot(BeZero())
		})
		Specify("The default project was found", func() {
			Expect(defaultProject).ToNot(BeNil())
		})

		When("Refetching the Default Organization Project using get by id", func() {
			var defaultProjectRefetched Project
			var defaultProjectRefetchedErr error

			JustBeforeEach(func() {
				defaultProjectRefetched, defaultProjectRefetchedErr = apiClient.Project(defaultProject.Id)
			})

			It("Should not fail", func() {
				Expect(defaultProjectRefetchedErr).To(BeNil())
			})
			It("Should return the correct project name", func() {
				Expect(defaultProjectRefetched.Name, "Default Organization Project")
			})

		})
	})
})
