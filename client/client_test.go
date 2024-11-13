package client_test

import (
	"testing"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

// This file wraps the test suite for the entire client folder

const (
	organizationId        = "organization0"
	defaultOrganizationId = "organization1"
)

var (
	ctrl               *gomock.Controller
	mockHttpClient     *http.MockHttpClientInterface
	apiClient          ApiClientInterface
	httpCall           *gomock.Call
	httpCall2          *gomock.Call
	organizationIdCall *gomock.Call
)

var _ = BeforeSuite(func() {
	ctrl = gomock.NewController(GinkgoT())
})

var _ = BeforeEach(func() {
	mockHttpClient = http.NewMockHttpClientInterface(ctrl)
	apiClient = NewApiClient(mockHttpClient, defaultOrganizationId)
})

var _ = AfterSuite(func() {
	ctrl.Finish()
})

func mockOrganizationIdCall() *gomock.Call {
	organizations := []Organization{{
		Id: organizationId,
	}}

	organizationIdCall = mockHttpClient.EXPECT().Get("/organizations", nil, gomock.Any()).Do(func(path string, params interface{}, response *[]Organization) {
		*response = organizations
	})

	return organizationIdCall
}

func TestApiClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Client Tests")
}
