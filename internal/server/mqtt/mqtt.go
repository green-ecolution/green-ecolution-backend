package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt/entities/sensor/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

type Mqtt struct {
	cfg    *config.Config
	svc    *service.Services
	mapper sensor.MqttMqttMapper
}

func NewMqtt(cfg *config.Config, services *service.Services) *Mqtt {
	return &Mqtt{
		cfg:    cfg,
		svc:    services,
		mapper: &generated.MqttMqttMapperImpl{},
	}
}

func (m *Mqtt) RunSubscriber(ctx context.Context) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(m.cfg.MQTT.Broker)
	opts.SetClientID(m.cfg.MQTT.ClientID)
	opts.SetUsername(m.cfg.MQTT.Username)
	opts.SetPassword(m.cfg.MQTT.Password)

	opts.OnConnect = func(_ MQTT.Client) {
		slog.Info("Connected to MQTT Broker")
		m.svc.MqttService.SetConnected(true)
	}
	opts.OnConnectionLost = func(_ MQTT.Client, err error) {
		slog.Error("Connection to MQTT Broker lost", "error", err)
		m.svc.MqttService.SetConnected(false)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("Error connecting to MQTT Broker", "error", token.Error())
		return
	}

	token := client.Subscribe(m.cfg.MQTT.Topic, 1, m.handleMqttMessage)
	go func(token MQTT.Token) {
		_ = token.Wait()
		if token.Error() != nil {
			slog.Error("Error while subscribing to MQTT Broker", "error", token.Error())
		}
	}(token)

	<-ctx.Done()
	slog.Info("Shutting down MQTT Subscriber")
}

func (m *Mqtt) handleMqttMessage(_ MQTT.Client, msg MQTT.Message) {
	sensorData, err := m.convertToMqttPayloadResponse(msg)
	if err != nil {
		slog.Error("Error while converting MQTT payload to sensor data", "error", err)
	}

	slog.Info("Logging sensor data", "sensorData", sensorData)

	domainPayload := m.mapper.FromResponse(sensorData)
	_, err = m.svc.MqttService.HandleMessage(context.Background(), domainPayload)
	if err != nil {
		slog.Error("Error handling message", "error", err)
		return
	}
}

func (m *Mqtt) convertToMqttPayloadResponse(msg MQTT.Message) (*sensor.MqttPayloadResponse, error) {
	var raw map[string]any
	if err := json.Unmarshal(msg.Payload(), &raw); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	endDeviceIDs := raw["end_device_ids"].(map[string]any)
	uplinkMessage := raw["uplink_message"].(map[string]any)
	decodedPayload := uplinkMessage["decoded_payload"].(map[string]any)

	// Parse temperature from string to float64
	temperature, err := strconv.ParseFloat(decodedPayload["temperature"].(string), 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing temperature: %w", err)
	}

	payload := &sensor.MqttPayloadResponse{
		DeviceID:    endDeviceIDs["device_id"].(string),
		Battery:     decodedPayload["battery"].(float64),
		Humidity:    decodedPayload["humidity"].(float64),
		Temperature: temperature,
		Latitude:    decodedPayload["latitude"].(float64),
		Longitude:   decodedPayload["longitude"].(float64),
		Watermarks: []sensor.WatermarkResponse{
			{
				Resistance: int(decodedPayload["watermarkOneResistanceValue"].(float64)),
				Centibar:   int(decodedPayload["watermarkOneCentibarValue"].(float64)),
				Depth:      30,
			},
			{
				Resistance: int(decodedPayload["watermarkTwoResistanceValue"].(float64)),
				Centibar:   int(decodedPayload["watermarkTwoCentibarValue"].(float64)),
				Depth:      60,
			},
			{
				Resistance: int(decodedPayload["watermarkThreeResistanceValue"].(float64)),
				Centibar:   int(decodedPayload["watermarkThreeCentibarValue"].(float64)),
				Depth:      90,
			},
		},
	}

	return payload, nil
}
