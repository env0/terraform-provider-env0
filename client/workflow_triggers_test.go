package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workflow Triggers", func() {
	const environmentId = "environmentId"
	mockTrigger := WorkflowTrigger{
		Id:            "id1",
		Name:          "name",
		WorkspaceName: "workspaceName",
		ProjectId:     "projectId",
	}

	Describe("WorkflowTriggerCreate", func() {
		var trigger WorkflowTrigger

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Put("/environments/"+environmentId+"/downstream", WorkflowTriggerCreatePayload{
					DownstreamEnvironmentIds: []string{mockTrigger.Id},
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *WorkflowTrigger) {
					*response = mockTrigger
				})

			trigger, _ = apiClient.WorkflowTriggerCreate(environmentId, WorkflowTriggerCreatePayload{
				DownstreamEnvironmentIds: []string{mockTrigger.Id},
			})
		})

		It("Should send PUT request with params", func() {
			httpCall.Times(1)
		})

		It("Should return created trigger", func() {
			Expect(trigger).To(Equal(mockTrigger))
		})
	})

	Describe("WorkflowTrigger", func() {
		var queriedTriggers []WorkflowTrigger
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("environments/"+environmentId+"/downstream", nil, gomock.Any()).
				Do(func(path string, request interface{}, response *[]WorkflowTrigger) {
					*response = []WorkflowTrigger{mockTrigger}
				})
			queriedTriggers, _ = apiClient.WorkflowTrigger(environmentId)
		})

		It("Should send GET request with project id", func() {
			httpCall.Times(1)
		})

		It("Should return correct triggers", func() {
			Expect(queriedTriggers).To(Equal([]WorkflowTrigger{mockTrigger}))
		})
	})
})
