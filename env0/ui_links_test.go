package env0

import (
	"errors"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUnitUiBaseUrl(t *testing.T) {
	testCases := []struct {
		apiEndpoint string
		expected    string
	}{
		{"https://api.env0.com/", "https://app.env0.com"},
		{"https://api.env0.com", "https://app.env0.com"},
		{"https://api-dev.dev.env0.com/", "https://dev.dev.env0.com"},
		{"https://api-bors.dev.env0.com/", "https://bors.dev.env0.com"},
		{"https://api-pr-1234.dev.env0.com/", "http://pr-1234.dev.env0.com"},
		{"https://self-hosted.example.com/", ""},
		{"https://api-.dev.env0.com/", ""},
		{"", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.apiEndpoint, func(t *testing.T) {
			assert.Equal(t, testCase.expected, uiBaseUrl(testCase.apiEndpoint))
		})
	}
}

func TestUnitDeploymentUrl(t *testing.T) {
	t.Run("builds the deployment link with the organization id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := client.NewMockApiClientInterface(ctrl)
		mock.EXPECT().ApiEndpoint().Times(1).Return("https://api-dev.dev.env0.com/")
		mock.EXPECT().OrganizationId().Times(1).Return("organization0", nil)

		assert.Equal(t,
			"https://dev.dev.env0.com/p/project0/environments/environment0/deployments/deployment0?organizationId=organization0",
			deploymentUrl(mock, "project0", "environment0", "deployment0"))
	})

	t.Run("omits the organization id when it cannot be fetched", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := client.NewMockApiClientInterface(ctrl)
		mock.EXPECT().ApiEndpoint().Times(1).Return("https://api.env0.com/")
		mock.EXPECT().OrganizationId().Times(1).Return("", errors.New("error"))

		assert.Equal(t,
			"https://app.env0.com/p/project0/environments/environment0/deployments/deployment0",
			deploymentUrl(mock, "project0", "environment0", "deployment0"))
	})

	t.Run("returns empty for an unknown api endpoint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := client.NewMockApiClientInterface(ctrl)
		mock.EXPECT().ApiEndpoint().Times(1).Return("https://self-hosted.example.com/")

		assert.Empty(t, deploymentUrl(mock, "project0", "environment0", "deployment0"))
	})
}
