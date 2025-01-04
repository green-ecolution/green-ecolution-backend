/*
Openrouteservice

Testing MatrixServiceAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package ors

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	openapiclient "github.com/green-ecolution/green-ecolution-backend/ors"
)

func Test_ors_MatrixServiceAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test MatrixServiceAPIService GetDefault1", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var profile string

		resp, httpRes, err := apiClient.MatrixServiceAPI.GetDefault1(context.Background(), profile).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
