/*
Green Space Management API

Testing TreeSensorAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package client

import (
	"context"
	"testing"

	openapiclient "github.com/green-ecolution/green-ecolution-backend/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_client_TreeSensorAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test TreeSensorAPIService AddSensorToTree", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var treeId string

		resp, httpRes, err := apiClient.TreeSensorAPI.AddSensorToTree(context.Background(), treeId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeSensorAPIService GetTreeBySensorId", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var sensorId string

		resp, httpRes, err := apiClient.TreeSensorAPI.GetTreeBySensorId(context.Background(), sensorId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeSensorAPIService GetTreeSensor", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var treeId string

		resp, httpRes, err := apiClient.TreeSensorAPI.GetTreeSensor(context.Background(), treeId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test TreeSensorAPIService RemoveSensorFromTree", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var treeId string
		var sensorId string

		resp, httpRes, err := apiClient.TreeSensorAPI.RemoveSensorFromTree(context.Background(), treeId, sensorId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
