package client_test

import (
	"encoding/json"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("TemplateCreatePayload", func() {
		It("Should omit Github Installation Id when it's not there", func() {
			payload := TemplateCreatePayload{}
			jsonPayload, _ := json.Marshal(payload)
			var parsedPayload map[string]interface{}
			json.Unmarshal(jsonPayload, &parsedPayload)
			Expect(parsedPayload["githubInstallationId"]).To(BeNil())
		})
	})
})
