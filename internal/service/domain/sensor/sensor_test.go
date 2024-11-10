package sensor

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)


func TestSensorService_GetAll(t *testing.T) {
	t.Run("should return all sensor", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		svc := NewSensorService(repo)

		// when
		repo.EXPECT().GetAll(context.Background()).Return(getTestSensors(), nil)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, getTestSensors(), sensors)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		svc := NewSensorService(repo)

		repo.EXPECT().GetAll(context.Background()).Return(nil, storage.ErrSensorNotFound)
		sensors, err := svc.GetAll(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, sensors)
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		repo := storageMock.NewMockSensorRepository(t)
		svc := NewSensorService(repo)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		// given
		svc := NewSensorService(nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

func getTestSensors() []*entities.Sensor {
	testDate := time.Date(2023, time.November, 10, 10, 0, 0, 0, time.UTC)

	return []*entities.Sensor{
		{
			ID:        1,
			CreatedAt: testDate,
			UpdatedAt: testDate,
			Status:    entities.SensorStatusOnline,
			Data: []*entities.SensorData{
				{
					ID:        101,
					CreatedAt: testDate,
					UpdatedAt: testDate,
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
						ReceivedAt:     &testDate,
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
									Rssi:        -45,
									ChannelRssi: -42,
									Snr:         9.5,
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
}
