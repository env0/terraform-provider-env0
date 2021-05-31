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
})
