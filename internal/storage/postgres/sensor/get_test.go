package sensor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_GetAll(t *testing.T) {
	t.Run("should return all sensors", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetAll(context.Background())

		fmt.Print(got[len(got)-1].Data)

		// then
		assert.NoError(t, err)
		assert.Equal(t, len(getTestSensors()), len(got))
		for i, sensor := range got {
			assert.Equal(t, getTestSensors()[i].ID, sensor.ID)
			assert.Equal(t, getTestSensors()[i].Status, sensor.Status)
			assert.NotZero(t, sensor.CreatedAt)
			assert.NotZero(t, sensor.UpdatedAt)
		}
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetByID(t *testing.T) {
	t.Run("should return sensor by id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, getTestSensors()[0].ID, got.ID)
		assert.Equal(t, getTestSensors()[0].Status, got.Status)
		assert.NotZero(t, got.CreatedAt)
		assert.NotZero(t, got.UpdatedAt)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is negative", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is zero", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetByID(ctx, 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetStatusByID(t *testing.T) {
	t.Run("should return sensor status by id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, getTestSensors()[0].Status, *got)
	})

	t.Run("should return error when sensor not found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is negative", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when sensor id is zero", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetStatusByID(ctx, 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetStatusByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetSensorByStatus(t *testing.T) {
	t.Run("should return sensors by status", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorByStatus(ctx, &getTestSensors()[0].Status)

		// then
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, getTestSensors()[0].ID, got[0].ID)
		assert.Equal(t, getTestSensors()[0].Status, got[0].Status)
	})

	t.Run("should return empty slice when no sensors match status", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		status := entities.SensorStatus("offline")
		got, err := r.GetSensorByStatus(ctx, &status)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when status is nil", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorByStatus(ctx, nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetSensorByStatus(ctx, &getTestSensors()[0].Status)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestSensorRepository_GetSensorDataByID(t *testing.T) {
	t.Run("should return sensor data for valid id", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/sensor")
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorDataByID(ctx, 1)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, got)
		for _, data := range got {
			assert.Equal(t, int32(1), data.ID)
			assert.NotZero(t, data.CreatedAt)
			assert.NotZero(t, data.UpdatedAt)
		}
	})

	t.Run("should return empty slice when no data found", func(t *testing.T) {
		// given
		ctx := context.Background()
		suite.ResetDB(t)
		r := NewSensorRepository(suite.Store, defaultSensorMappers())

		// when
		got, err := r.GetSensorDataByID(ctx, 999)

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewSensorRepository(suite.Store, defaultSensorMappers())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetSensorDataByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
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
