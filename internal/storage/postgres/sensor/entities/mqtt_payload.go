package entities

import (
	"time"
)

type MqttIdentifierApplicationID struct {
	ApplicationID string `json:"application_id"`
}

type MqttIdentifierDeviceID struct {
	DeviceID       string                      `json:"device_id"`
	ApplicationIDs MqttIdentifierApplicationID `json:"application_ids"`
	DevEUI         string                      `json:"dev_eui"`
	JoinEUI        string                      `json:"join_eui"`
	DevAddr        string                      `json:"dev_addr"`
}

type MqttIdentifier struct {
	DeviceIDs MqttIdentifierDeviceID `json:"device_ids"`
}

type MqttDecodedPayload struct {
	Battery  float64 `json:"battery"`
	Humidity int     `json:"humidity"`
	Raw      int     `json:"raw"`
}

type MqttRxMetadataGatewayIDs struct {
	GatewayID string `json:"gateway_id"`
}

type MqttRxMetadataPacketBroker struct {
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

type MqttLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type MqttRxMetadata struct {
	GatewayIDs      MqttRxMetadataGatewayIDs   `json:"gateway_ids"`
	PacketBroker    MqttRxMetadataPacketBroker `json:"packet_broker"`
	Time            *time.Time                 `json:"time"`
	Rssi            int                        `json:"rssi"`
	ChannelRssi     int                        `json:"channel_rssi"`
	Snr             float64                    `json:"snr"`
	FrequencyOffset string                     `json:"frequency_offset"`
	Location        MqttLocation               `json:"location"`
	UplinkToken     string                     `json:"uplink_token"`
	ReceivedAt      *time.Time                 `json:"received_at"`
}

type MqttUplinkSettingsLora struct {
	Bandwidth       int    `json:"bandwidth"`
	SpreadingFactor int    `json:"spreading_factor"`
	CodingRate      string `json:"coding_rate"`
}

type MqttUplinkSettingsDataRate struct {
	Lora MqttUplinkSettingsLora `json:"lora"`
}

type MqttUplinkSettings struct {
	DataRate  MqttUplinkSettingsDataRate `json:"data_rate"`
	Frequency string                     `json:"frequency"`
}

type MqttVersionIDs struct {
	BrandID         string `json:"brand_id"`
	ModelID         string `json:"model_id"`
	HardwareVersion string `json:"hardware_version"`
	FirmwareVersion string `json:"firmware_version"`
	BandID          string `json:"band_id"`
}

type MqttNetworkIDs struct {
	NetID          string `json:"net_id"`
	NSID           string `json:"nsid"`
	TenantID       string `json:"tenant_id"`
	ClusterID      string `json:"cluster_id"`
	ClusterAddress string `json:"cluster_address"`
	TenantAddress  string `json:"tenant_address"`
}

type MqttUplinkMessage struct {
	SessionKeyID    string             `json:"session_key_id"`
	FPort           int                `json:"f_port"`
	Fcnt            int                `json:"fcnt"`
	FRMPayload      string             `json:"frm_payload"`
	DecodedPayload  MqttDecodedPayload `json:"decoded_payload"`
	RxMetadata      []MqttRxMetadata   `json:"rx_metadata"`
	Settings        MqttUplinkSettings `json:"settings"`
	ReceivedAt      *time.Time         `json:"received_at"`
	Confirmed       bool               `json:"confirmed"`
	ConsumedAirtime string             `json:"consumed_airtime"`
	VersionIDs      MqttVersionIDs     `json:"version_ids"`
	NetworkIDs      MqttNetworkIDs     `json:"network_ids"`
}

type MqttDataPayload struct {
	Type           string                 `json:"type"`
	EndDeviceIDs   MqttIdentifierDeviceID `json:"end_device_ids"`
	CorrelationIDs []string               `json:"correlation_ids"`
	ReceivedAt     *time.Time             `json:"received_at"`
	UplinkMessage  MqttUplinkMessage      `json:"uplink_message"`
}

type MqttVisibility struct {
	Rights []string `json:"rights"`
}

type MqttData struct {
	Name           string            `json:"name"`
	Time           *time.Time        `json:"time"`
	Identifiers    []MqttIdentifier  `json:"identifiers"`
	Data           MqttDataPayload   `json:"data"`
	CorrelationIDs []string          `json:"correlation_ids"`
	Origin         string            `json:"origin"`
	Context        map[string]string `json:"context"`
	Visibility     MqttVisibility    `json:"visibility"`
	UniqueID       string            `json:"unique_id"`
}

type MqttPayload struct {
	EndDeviceIDs   MqttIdentifierDeviceID `json:"end_device_ids"`
	CorrelationIDs []string               `json:"correlation_ids"`
	ReceivedAt     *time.Time             `json:"received_at"`
	UplinkMessage  MqttUplinkMessage      `json:"uplink_message"`
}
