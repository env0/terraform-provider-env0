package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Agent Client", func() {
	errorMock := errors.New("error")

	Describe("Agents", func() {
		mockAgent := Agent{
			AgentKey: "key",
		}

		expectedResponse := []Agent{mockAgent}

		Describe("Successful", func() {
			var (
				actualResult []Agent
				err          error
			)

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Get("/agents", gomock.Any(), gomock.Any()).
					Do(func(path string, request any, response *[]Agent) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.Agents()
			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("Should return the GET result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var (
				actualResult []Agent
				err          error
			)

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Get("/agents", gomock.Any(), gomock.Any()).
					Return(errorMock)

				actualResult, err = apiClient.Agents()
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("AgentPools", func() {
		mockPool := AgentPool{
			Id:       "pool-id",
			Name:     "pool-name",
			AgentKey: "pool-key",
		}

		expectedResponse := []AgentPool{mockPool}

		Describe("Successful", func() {
			var (
				actualResult []AgentPool
				err          error
			)

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Get("/agents", gomock.Any(), gomock.Any()).
					Do(func(path string, request any, response *[]AgentPool) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.AgentPools()
			})

			It("Should return the result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Get("/agents", gomock.Any(), gomock.Any()).
					Return(errorMock)

				_, err = apiClient.AgentPools()
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentPoolCreate", func() {
		mockPool := AgentPool{
			Id:       "pool-id",
			Name:     "pool-name",
			AgentKey: "pool-key",
		}

		payload := AgentPoolCreatePayload{
			Name:        "pool-name",
			Description: "desc",
		}

		Describe("Successful", func() {
			var (
				actualResult *AgentPool
				err          error
			)

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Post("/agents", gomock.Any(), gomock.Any()).
					Do(func(path string, request any, response *AgentPool) {
						*response = mockPool
					})
				actualResult, err = apiClient.AgentPoolCreate(payload)
			})

			It("Should return the result", func() {
				Expect(*actualResult).To(Equal(mockPool))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Post("/agents", gomock.Any(), gomock.Any()).
					Return(errorMock)

				_, err = apiClient.AgentPoolCreate(payload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentPool", func() {
		mockPool := AgentPool{
			Id:   "pool-id",
			Name: "pool-name",
		}

		Describe("Successful", func() {
			var (
				actualResult *AgentPool
				err          error
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+mockPool.Id, nil, gomock.Any()).
					Do(func(path string, request any, response *AgentPool) {
						*response = mockPool
					})
				actualResult, err = apiClient.AgentPool(mockPool.Id)
			})

			It("Should return the result", func() {
				Expect(*actualResult).To(Equal(mockPool))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+mockPool.Id, nil, gomock.Any()).
					Return(errorMock)

				_, err = apiClient.AgentPool(mockPool.Id)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentPoolUpdate", func() {
		mockPool := AgentPool{
			Id:   "pool-id",
			Name: "updated-name",
		}

		payload := AgentPoolUpdatePayload{
			Name: "updated-name",
		}

		Describe("Successful", func() {
			var (
				actualResult *AgentPool
				err          error
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Patch("/agents/"+mockPool.Id, payload, gomock.Any()).
					Do(func(path string, request any, response *AgentPool) {
						*response = mockPool
					})
				actualResult, err = apiClient.AgentPoolUpdate(mockPool.Id, payload)
			})

			It("Should return the result", func() {
				Expect(*actualResult).To(Equal(mockPool))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Patch("/agents/"+mockPool.Id, payload, gomock.Any()).
					Return(errorMock)

				_, err = apiClient.AgentPoolUpdate(mockPool.Id, payload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentPoolDelete", func() {
		poolId := "pool-id"

		Describe("Successful", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Delete("/agents/"+poolId, nil).
					Return(nil)
				err = apiClient.AgentPoolDelete(poolId)
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Delete("/agents/"+poolId, nil).
					Return(errorMock)

				err = apiClient.AgentPoolDelete(poolId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentSecretCreate", func() {
		agentId := "agent-id"
		mockSecret := AgentSecret{
			Id:      "secret-id",
			Secret:  "secret-value",
			AgentId: agentId,
		}

		payload := AgentSecretCreatePayload{
			Description: "desc",
		}

		Describe("Successful", func() {
			var (
				actualResult *AgentSecret
				err          error
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Post("/agents/"+agentId+"/secrets", payload, gomock.Any()).
					Do(func(path string, request any, response *AgentSecret) {
						*response = mockSecret
					})
				actualResult, err = apiClient.AgentSecretCreate(agentId, payload)
			})

			It("Should return the result", func() {
				Expect(*actualResult).To(Equal(mockSecret))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Post("/agents/"+agentId+"/secrets", payload, gomock.Any()).
					Return(errorMock)

				_, err = apiClient.AgentSecretCreate(agentId, payload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentSecrets", func() {
		agentId := "agent-id"
		mockSecret := AgentSecret{
			Id:      "secret-id",
			AgentId: agentId,
		}

		expectedResponse := []AgentSecret{mockSecret}

		Describe("Successful", func() {
			var (
				actualResult []AgentSecret
				err          error
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+agentId+"/secrets", nil, gomock.Any()).
					Do(func(path string, request any, response *[]AgentSecret) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.AgentSecrets(agentId)
			})

			It("Should return the result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+agentId+"/secrets", nil, gomock.Any()).
					Return(errorMock)

				_, err = apiClient.AgentSecrets(agentId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentSecretDelete", func() {
		agentId := "agent-id"
		secretId := "secret-id"

		Describe("Successful", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Delete("/agents/"+agentId+"/secrets/"+secretId, nil).
					Return(nil)
				err = apiClient.AgentSecretDelete(agentId, secretId)
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Delete("/agents/"+agentId+"/secrets/"+secretId, nil).
					Return(errorMock)

				err = apiClient.AgentSecretDelete(agentId, secretId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})
		})
	})

	Describe("AgentValues", func() {
		expectedResponse := "response"
		agentId := "id"

		Describe("Successful", func() {
			var (
				actualResult string
				err          error
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+agentId+"/values", nil, gomock.Any()).
					Do(func(path string, request any, response *string) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.AgentValues(agentId)
			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("Should return the GET result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var (
				err          error
				actualResult string
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+agentId+"/values", nil, gomock.Any()).
					Return(errorMock)

				actualResult, err = apiClient.AgentValues(agentId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
				Expect(actualResult).To(Equal(""))
			})
		})
	})
})
