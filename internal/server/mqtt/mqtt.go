package mqtt

import (
	"context"
	"encoding/json"
	"log/slog"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/green-ecolution/green-ecolution-backend/config"
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
	var sensorData sensor.MqttPayloadResponse
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		slog.Error("Error unmarshalling message", "error", err)
		return
	}

	domainPayload := m.mapper.FromResponse(&sensorData)
	_, err := m.svc.MqttService.HandleMessage(context.Background(), domainPayload)
	if err != nil {
		slog.Error("Error handling message", "error", err)
		return
	}
}
