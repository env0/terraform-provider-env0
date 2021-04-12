package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/env0/terraform-provider-env0/internal/api"
	"github.com/env0/terraform-provider-env0/internal/rest"
)

var apiClient *api.ApiClient

var _ = BeforeSuite(func() {
	restClient, _ := rest.NewRestClientFromEnv()
	apiClient = api.NewApiClient(restClient)
})

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
