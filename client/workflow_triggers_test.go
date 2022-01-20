package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workflow Triggers", func() {
	const environmentId = "environmentId"
	mockTrigger := []WorkflowTrigger{
		{
			Id:            "id1",
			Name:          "name",
			WorkspaceName: "workspaceName",
			ProjectId:     "projectId",
		},
	}

	var triggers []WorkflowTrigger
	Describe("WorkflowTriggerUpsert", func() {

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Put("/environments/"+environmentId+"/downstream", WorkflowTriggerUpsertPayload{
					DownstreamEnvironmentIds: []string{mockTrigger[0].Id},
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *[]WorkflowTrigger) {
					*response = mockTrigger
				})

			triggers, _ = apiClient.WorkflowTriggerUpsert(environmentId, WorkflowTriggerUpsertPayload{
				DownstreamEnvironmentIds: []string{mockTrigger[0].Id},
			})
		})

		It("Should return created triggers", func() {
			Expect(triggers).To(Equal(mockTrigger))
		})
	})

	Describe("WorkflowTrigger", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("environments/"+environmentId+"/downstream", nil, gomock.Any()).
				Do(func(path string, request interface{}, response *[]WorkflowTrigger) {
					*response = mockTrigger
				})
			triggers, _ = apiClient.WorkflowTrigger(environmentId)
		})

		It("Should return correct triggers", func() {
			Expect(triggers).To(Equal(mockTrigger))
		})
	})
})
