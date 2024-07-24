package sensor

import (
	"github.com/google/uuid"
	"time"
)

type MqttIdentifierApplicationIDEntity struct {
	ApplicationID string `json:"application_id"`
}

type MqttIdentifierDeviceIDEntity struct {
	DeviceID       string                            `json:"device_id"`
	ApplicationIDs MqttIdentifierApplicationIDEntity `json:"application_ids"`
	DevEUI         string                            `json:"dev_eui"`
	JoinEUI        string                            `json:"join_eui"`
	DevAddr        string                            `json:"dev_addr"`
}

type MqttIdentifierEntity struct {
	DeviceIDs MqttIdentifierDeviceIDEntity `bson:"device_ids"`
}

type MqttDecodedPayloadEntity struct {
	Battery  float64 `json:"battery"`
	Humidity int     `json:"humidity"`
	Raw      int     `json:"raw"`
}

type MqttRxMetadataGatewayIDsEntity struct {
	GatewayID string `json:"gateway_id"`
}

type MqttRxMetadataPacketBrokerEntity struct {
	MessageID            string `json:"message_id"`
	ForwarderNetID       string `json:"forwarder_net_id"`
	ForwarderTenantID    string `json:"forwarder_tenant_id"`
	ForwarderClusterID   string `json:"forwarder_cluster_id"`
	ForwarderGatewayID   string `json:"forwarder_gateway_id"`
	ForwarderGatewayEUI  string `json:"forwarder_gateway_eui"`
	HomeNetworkNetID     string `json:"home_network_net_id"`
	HomeNetworkTenantID  string `json:"home_network_tenant_id"`
	HomeNetworkClusterID string `json:"home_network_cluster_id"`
}

type MqttLocationEntity struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type MqttRxMetadataEntity struct {
<<<<<<< HEAD:internal/storage/mongodb/entities/sensor/sensor.go
	GatewayIDs      MqttRxMetadataGatewayIDsEntity   `bson:"gateway_i_ds"`
	PacketBroker    MqttRxMetadataPacketBrokerEntity `bson:"packet_broker"`
	Time            *time.Time                       `bson:"time"`
	Rssi            int                              `bson:"rssi"`
	ChannelRssi     int                              `bson:"channel_rssi"`
	Snr             float64                          `bson:"snr"`
	FrequencyOffset string                           `bson:"frequency_offset"`
	Location        MqttLocationEntity               `bson:"location"`
	UplinkToken     string                           `bson:"uplink_token"`
	ReceivedAt      *time.Time                       `bson:"received_at"`
=======
	GatewayIDs      MqttRxMetadataGatewayIDsEntity   `json:"gateway_i_ds"`
	PacketBroker    MqttRxMetadataPacketBrokerEntity `json:"packet_broker"`
	Time            *time.Time                       `json:"time"`
	Rssi            int                              `json:"rssi"`
	ChannelRssi     int                              `json:"channel_rssi"`
	Snr             float64                          `json:"snr"`
	FrequencyOffset string                           `json:"frequency_offset"`
	Location        MqttLocationEntity               `json:"location"`
	UplinkToken     string                           `json:"uplink_token"`
	RecievedAt      *time.Time                       `json:"recieved_at"`
>>>>>>> 3e26622 (swithch from mangodb to postgreSQL):internal/storage/entities/sensor/sensor.go
}

type MqttUplinkSettingsLoraEntity struct {
	Bandwidth       int    `json:"bandwidth"`
	SpreadingFactor int    `json:"spreading_factor"`
	CodingRate      string `json:"coding_rate"`
}

type MqttUplinkSettingsDataRateEntity struct {
	Lora MqttUplinkSettingsLoraEntity `bson:"lora"`
}

type MqttUplinkSettingsEntity struct {
	DataRate  MqttUplinkSettingsDataRateEntity `json:"data_rate"`
	Frequency string                           `json:"frequency" bson:"frequency"`
}

type MqttVersionIDsEntity struct {
	BrandID         string `json:"brand_id"`
	ModelID         string `json:"model_id"`
	HardwareVersion string `json:"hardware_version"`
	FirmwareVersion string `json:"firmware_version"`
	BandID          string `json:"band_id"`
}

type MqttNetworkIDsEntity struct {
	NetID          string `json:"net_id"`
	NSID           string `json:"nsid"`
	TenantID       string `json:"tenant_id"`
	ClusterID      string `json:"cluster_id"`
	ClusterAddress string `json:"cluster_address"`
	TenantAddress  string `json:"tenant_address"`
}

type MqttUplinkMessageEntity struct {
	SessionKeyID    string                   `json:"session_key_id"`
	FPort           int                      `json:"f_port"`
	Fcnt            int                      `json:"fcnt"`
	FRMPayload      string                   `json:"frm_payload"`
	DecodedPayload  MqttDecodedPayloadEntity `json:"decoded_payload"`
	RxMetadata      []MqttRxMetadataEntity   `json:"rx_metadata"`
	Settings        MqttUplinkSettingsEntity `json:"settings"`
	ReceivedAt      *time.Time               `json:"received_at"`
	Confirmed       bool                     `json:"confirmed"`
	ConsumedAirtime string                   `json:"consumed_airtime"`
	VersionIDs      MqttVersionIDsEntity     `json:"version_ids"`
	NetworkIDs      MqttNetworkIDsEntity     `json:"network_ids"`
}

type MqttDataPayloadEntity struct {
	Type           string                       `json:"type"`
	EndDeviceIDs   MqttIdentifierDeviceIDEntity `json:"end_device_ids"`
	CorrelationIDs []string                     `json:"correlation_ids"`
	ReceivedAt     *time.Time                   `json:"received_at"`
	UplinkMessage  MqttUplinkMessageEntity      `json:"uplink_message"`
}

type MqttVisibilityEntity struct {
	Rights []string `json:"rights"`
}

type MqttDataEntity struct {
	Name           string                 `json:"name"`
	Time           *time.Time             `json:"time"`
	Identifiers    []MqttIdentifierEntity `json:"identifiers"`
	Data           MqttDataPayloadEntity  `json:"data"`
	CorrelationIDs []string               `json:"correlation_ids"`
	Origin         string                 `json:"origin"`
	Context        map[string]string      `json:"context"`
	Visibility     MqttVisibilityEntity   `json:"visibility"`
	UniqueID       string                 `json:"unique_id"`
}

type MqttPayloadEntity struct {
	EndDeviceIDs   MqttIdentifierDeviceIDEntity `json:"end_device_ids"`
	CorrelationIDs []string                     `json:"correlation_ids"`
	ReceivedAt     *time.Time                   `json:"received_at"`
	UplinkMessage  MqttUplinkMessageEntity      `json:"uplink_message"`
}

type MqttEntity struct {
	ID     uuid.UUID         `json:"id"`
	TreeID string            `json:"tree_id"`
	Data   MqttPayloadEntity `json:"data"`
}

func (m *MqttEntity) GetID() string {
	return m.ID.String()
}

func (m *MqttEntity) SetID(id string) error {
	objID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	m.ID = objID
	return nil
}
