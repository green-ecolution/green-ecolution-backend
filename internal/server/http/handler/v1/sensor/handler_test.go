package sensor

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllSensors(t *testing.T) {
	t.Run("should return all sensors successfully with full MqttPayload", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := GetAllSensors(mockSensorService)
	
		currentTime := time.Now()

		expectedSensors := []*entities.Sensor{
			{
				ID:        1,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
				Status:    entities.SensorStatusOnline,
				Data: []*entities.SensorData{
					{
						ID:        101,
						CreatedAt: currentTime,
						UpdatedAt: currentTime,
						Data: &entities.MqttPayload{
							EndDeviceIDs: entities.MqttIdentifierDeviceID{
								DeviceID: "Device123",
								ApplicationIDs: entities.MqttIdentifierApplicationID{
									ApplicationID: "AppID123",
								},
								DevEUI:  "00-14-22-01-23-45",
								JoinEUI: "00-15-33-02-34-56",
							},
							CorrelationIDs: []string{"corrID1", "corrID2"},
							ReceivedAt:     &currentTime,
							UplinkMessage: entities.MqttUplinkMessage{
								SessionKeyID:   "sessionKey1",
								FPort:          1,
								Fcnt:           10,
								FRMPayload:     "payloadData",
								DecodedPayload: entities.MqttDecodedPayload{Battery: 85.0, Humidity: 55, Raw: 123},
								RxMetadata: []entities.MqttRxMetadata{
									{
										GatewayIDs: entities.MqttRxMetadataGatewayIDs{
											GatewayID: "Gateway123",
										},
										Rssi:         -45,
										ChannelRssi:  -42,
										Snr:          9.5,
										Location: entities.MqttLocation{
											Latitude:  52.5200,
											Longitude: 13.4050,
											Altitude:  34.0,
										},
									},
								},
								Settings: entities.MqttUplinkSettings{
									DataRate: entities.MqttUplinkSettingsDataRate{
										Lora: entities.MqttUplinkSettingsLora{
											Bandwidth:       125,
											SpreadingFactor: 7,
											CodingRate:      "4/5",
										},
									},
									Frequency: "868100000",
								},
								Confirmed:       true,
								ConsumedAirtime: "0.123s",
							},
						},
					},
				},
			},
		}

		mockSensorService.EXPECT().GetAll(mock.Anything).Return(expectedSensors, nil)
		app.Get("/v1/sensor", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/sensor", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.SensorListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)

		assert.Len(t, response.Data, 1)
		sensorData := response.Data[0]
		assert.Equal(t, expectedSensors[0].ID, sensorData.ID)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return empty sensor list when no sensors found", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(mock.Anything).Return([]*entities.Sensor{}, nil)
		app.Get("/v1/sensor", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/sensor", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.SensorListResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)
		assert.Len(t, response.Data, 0)

		mockSensorService.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns an error", func(t *testing.T) {
		mockSensorService := serviceMock.NewMockSensorService(t)
		app := fiber.New()
		handler := GetAllSensors(mockSensorService)

		mockSensorService.EXPECT().GetAll(mock.Anything).Return(nil, errors.New("service error"))
		app.Get("/v1/sensor", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/sensor", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSensorService.AssertExpectations(t)
	})
}
