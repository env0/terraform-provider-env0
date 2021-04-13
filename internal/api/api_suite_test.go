package api_test

import (
	"github.com/golang/mock/gomock"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/env0/terraform-provider-env0/internal/api"
	"github.com/env0/terraform-provider-env0/internal/http/mocks"
)

const organizationId = "organization0"

var (
	ctrl           *gomock.Controller
	mockHttpClient *mocks.MockHttpClientInterface
	apiClient      *api.ApiClient
)

var _ = BeforeSuite(func() {
	ctrl = gomock.NewController(GinkgoT())
})

var _ = BeforeEach(func() {
	mockHttpClient = mocks.NewMockHttpClientInterface(ctrl)
	apiClient = api.NewApiClient(mockHttpClient, organizationId)
})

var _ = AfterEach(func() {
	ctrl.Finish()
})

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
