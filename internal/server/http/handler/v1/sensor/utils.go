package sensor

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var (
	currentTime = time.Now()

	TestSensor = &entities.Sensor{
		ID:        1,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Status:    entities.SensorStatusOnline,
		Data: []*entities.SensorData{
			{
				ID:        1,
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
	}

	TestSensorList = []*entities.Sensor{
		TestSensor,
	}
)
