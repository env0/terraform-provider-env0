package client_test

import (
	"encoding/json"
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Models", func() {
	Describe("TemplateCreatePayload", func() {
		DescribeTable("Github Installation Id",
			func(value int, expected types.GomegaMatcher) {
				payload := TemplateCreatePayload{
					GithubInstallationId: value,
				}
				jsonPayload, _ := json.Marshal(payload)
				var parsedPayload map[string]interface{}
				json.Unmarshal(jsonPayload, &parsedPayload)
				Expect(parsedPayload["githubInstallationId"]).To(expected)
			},
			Entry("Has value", 123, BeEquivalentTo(123)),
			Entry("No value", nil, BeNil()),
		)
	})

	Describe("ConfigurationVariable", func() {
		Describe("Deserialize", func() {
			It("On schema type is free text, enum should be nil", func() {
				var parsedPayload ConfigurationVariable
				json.Unmarshal([]byte(`{"schema": {"type": "string"}}`), &parsedPayload)
				Expect(parsedPayload.Schema.Type).Should(Equal("string"))
				Expect(parsedPayload.Schema.Enum).Should(BeNil())
			})

			It("On schema type is dropdown, enum should be present", func() {
				var parsedPayload ConfigurationVariable
				json.Unmarshal([]byte(`{"schema": {"type": "string", "enum": ["hello"]}}`), &parsedPayload)
				Expect(parsedPayload.Schema.Type).Should(Equal("string"))
				Expect(parsedPayload.Schema.Enum).Should(BeEquivalentTo([]string{"hello"}))
			})
		})

	})
})
