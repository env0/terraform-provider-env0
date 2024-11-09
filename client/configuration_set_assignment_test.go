package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Configuration Set", func() {
	scope := "environment"
	scopeId := "12345"
	setIds := []string{"1", "2", "3"}
	mockConfigurationSets := []ConfigurationSet{
		{
			Id: "1",
		},
		{
			Id: "2",
		},
		{
			Id: "3",
		},
	}

	Describe("assign configuration sets", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/configuration-sets/assignments/environment/12345?setIds=1,2,3", nil, nil).
				Do(func(path string, request interface{}, response *interface{}) {}).
				Times(1)
			err = apiClient.AssignConfigurationSets(scope, scopeId, setIds)
		})

		It("Should send post request", func() {})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("unassign configuration sets", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/configuration-sets/assignments/environment/12345", map[string]string{
				"setIds": "1,2,3",
			}).
				Do(func(path string, request interface{}) {}).
				Times(1)
			err = apiClient.UnassignConfigurationSets(scope, scopeId, setIds)
		})

		It("Should send delete request", func() {})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("get configuration sets by scope and scope id", func() {
		var configurationSets []ConfigurationSet

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/configuration-sets/assignments/environment/12345", nil, gomock.Any()).
				Do(func(path string, request interface{}, response *[]ConfigurationSet) {
					*response = mockConfigurationSets
				}).Times(1)
			configurationSets, _ = apiClient.ConfigurationSetsAssignments(scope, scopeId)
		})

		It("Should return configuration sets", func() {
			Expect(configurationSets).To(Equal(mockConfigurationSets))
		})
	})
})
