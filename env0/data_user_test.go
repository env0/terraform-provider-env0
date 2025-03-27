package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUserDataSource(t *testing.T) {
	user := client.OrganizationUser{
		User: client.User{
			Email:  "a@b.com",
			UserId: "1",
		},
	}

	otherUser := client.OrganizationUser{
		User: client.User{
			Email:  "c@d.com",
			UserId: "2",
		},
	}

	resourceType := "env0_user"
	resourceName := "test_user"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", user.User.UserId),
						resource.TestCheckResourceAttr(accessor, "email", user.User.Email),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockUsersCall := func(returnValue []client.OrganizationUser) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Users().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Get user by email", func(t *testing.T) {
		input := map[string]any{"email": user.User.Email}

		runUnitTest(t,
			getValidTestCase(input),
			mockUsersCall([]client.OrganizationUser{user, otherUser}),
		)
	})

	t.Run("Return error when user by email not found", func(t *testing.T) {
		input := map[string]any{"email": user.User.Email}

		runUnitTest(t,
			getErrorTestCase(input, "not find a user"),
			mockUsersCall([]client.OrganizationUser{otherUser}),
		)
	})

	t.Run("Throw error when multiple users by email found", func(t *testing.T) {
		input := map[string]any{"email": user.User.Email}

		runUnitTest(t,
			getErrorTestCase(input, "multiple users"),
			mockUsersCall([]client.OrganizationUser{user, otherUser, user}),
		)
	})
}
