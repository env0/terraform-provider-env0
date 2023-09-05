package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project Budget Client", func() {
	mockProjectBudget := ProjectBudget{
		Id:         "id",
		ProjectId:  "pid",
		Amount:     50,
		Timeframe:  "WEEKLY",
		Thresholds: []int{1, 2},
	}

	Describe("Get Project Budget", func() {
		var returnedProjectBudget *ProjectBudget

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/costs/project/"+mockProjectBudget.ProjectId+"/budget", gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *ProjectBudget) {
					*response = mockProjectBudget
				})
			httpCall.Times(1)
			returnedProjectBudget, _ = apiClient.ProjectBudget(mockProjectBudget.ProjectId)
		})

		It("Should return project budget", func() {
			Expect(*returnedProjectBudget).To(Equal(mockProjectBudget))
		})
	})

	Describe("Delete Project Budget", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/costs/project/"+mockProjectBudget.ProjectId+"/budget", gomock.Nil())
			httpCall.Times(1)
			apiClient.ProjectBudgetDelete(mockProjectBudget.ProjectId)
		})

		It("Should delete project budget", func() {})
	})

	Describe("Update Project Budget", func() {
		var returnedProjectBudget *ProjectBudget

		BeforeEach(func() {
			var updateProjectBudgetPayload ProjectBudgetUpdatePayload
			copier.Copy(&updateProjectBudgetPayload, &mockProjectBudget)

			httpCall = mockHttpClient.EXPECT().
				Put("/costs/project/"+mockProjectBudget.ProjectId+"/budget", &updateProjectBudgetPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *ProjectBudget) {
					*response = mockProjectBudget
				})
			httpCall.Times(1)
			returnedProjectBudget, _ = apiClient.ProjectBudgetUpdate(mockProjectBudget.ProjectId, &updateProjectBudgetPayload)
		})

		It("Should return project budget", func() {
			Expect(*returnedProjectBudget).To(Equal(mockProjectBudget))
		})
	})
})
