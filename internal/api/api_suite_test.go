package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/env0/terraform-provider-env0/internal/api"
)

var apiClient *api.ApiClient

var _ = BeforeSuite(func() {
	var err error
	apiClient, err = api.NewClientFromEnv()
	Expect(err).To(BeNil())
	Expect(apiClient).ToNot(BeNil())
})

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
