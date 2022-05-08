package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Client", func() {
	mockUser := OrganizationUser{
		User: User{
			Email:  "a@b.com",
			UserId: "1",
		},
	}

	Describe("Users", func() {
		var users []OrganizationUser
		mockUsers := []OrganizationUser{mockUser}
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Get("/organizations/"+organizationId+"/users", gomock.Any(), gomock.Any()).
					Do(func(path string, request interface{}, response *[]OrganizationUser) {
						*response = mockUsers
					}).Times(1)

				users, err = apiClient.Users()
			})

			It("Should return the user", func() {
				Expect(users).To(Equal(mockUsers))
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				mockOrganizationIdCall(organizationId)

				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get("/organizations/"+organizationId+"/users", gomock.Any(), gomock.Any()).
					Times(1).
					Return(expectedErr)
				_, err = apiClient.Users()
				Expect(expectedErr).Should(Equal(err))
			})
		})
	})
})
