package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAwsCredDataSource(t *testing.T) {
	awsCred := client.ApiKey{
		Id:             "11111",
		Name:           "testdata",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	}

	AwsCredFieldsByName := map[string]interface{}{"name": awsCred.Name}
	AwsCredFieldsById := map[string]interface{}{"id": awsCred.Id}

	resourceType := "env0_aws_credentials"
	resourceName := "testdata"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", awsCred.Id),
						resource.TestCheckResourceAttr(accessor, "name", awsCred.Name),
					),
				},
			},
		}
	}

	  mockGetAwsCredCall := func(returnValue client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
	  	return func(mock *client.MockApiClientInterface) {
	  		mock.EXPECT().AwsCredentials(awsCred.Id).AnyTimes().Return(returnValue, nil)
	  	}
	  }

	mockListAwsCredCall := func(returnValue []client.ApiKey) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AwsCredentialsList().AnyTimes().Return(returnValue, nil)
		}
	}

	  t.Run("By ID", func(t *testing.T) {
	  	runUnitTest(t,
	 		getValidTestCase(AwsCredFieldsById),
	  		mockGetAwsCredCall(awsCred),
	  	)
	  })

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(AwsCredFieldsByName),
			mockListAwsCredCall([]client.ApiKey{awsCred}),
		)
	})

}
