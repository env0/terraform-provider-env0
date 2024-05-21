package client_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Configuration Set", func() {
	scope := "environment"
	scopeId := "12345"
	setIds := []string{"1", "2", "3"}

	Describe("assign configuration sets", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Post("/configuration-sets/assignments/environment/12345?setIds=1,2,3", nil, nil).Times(1)
			apiClient.AssignConfigurationSets(scope, scopeId, setIds)
		})

		It("Should send post request", func() {})
	})

	Describe("unassign configuration sets", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/configuration-sets/assignments/environment/12345", map[string]string{
				"setIds": "1,2,3",
			}).Times(1)
			apiClient.UnassignConfigurationSets(scope, scopeId, setIds)
		})

		It("Should send delete request", func() {})
	})
})
