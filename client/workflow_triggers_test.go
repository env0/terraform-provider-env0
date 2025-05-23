package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Workflow Triggers", func() {
	const environmentId = "environmentId"
	mockTrigger := []WorkflowTrigger{
		{
			Id: "id1",
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
				Do(func(path string, request any, response *[]WorkflowTrigger) {
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
				Get("/environments/"+environmentId+"/downstream", nil, gomock.Any()).
				Do(func(path string, request any, response *[]WorkflowTrigger) {
					*response = mockTrigger
				})
			triggers, _ = apiClient.WorkflowTrigger(environmentId)
		})

		It("Should return correct triggers", func() {
			Expect(triggers).To(Equal(mockTrigger))
		})
	})

	Describe("SubscribeWorkflowTrigger", func() {
		var err error

		subscribePayload := WorkflowTriggerEnvironments{
			DownstreamEnvironmentIds: []string{"1", "2"},
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/environments/"+environmentId+"/downstream/subscribe", subscribePayload, nil)
			httpCall.Times(1)
			err = apiClient.SubscribeWorkflowTrigger(environmentId, subscribePayload)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("UnsubscribeWorkflowTrigger", func() {
		var err error

		unsubscribePayload := WorkflowTriggerEnvironments{
			DownstreamEnvironmentIds: []string{"1", "2"},
		}

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/environments/"+environmentId+"/downstream/unsubscribe", unsubscribePayload, nil)
			httpCall.Times(1)
			err = apiClient.UnsubscribeWorkflowTrigger(environmentId, unsubscribePayload)
		})

		It("Should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})
})
