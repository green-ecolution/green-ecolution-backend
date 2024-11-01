-- +goose Up
-- +goose StatementBegin
INSERT INTO sensors (id, status) VALUES (1, 'online');
ALTER SEQUENCE sensors_id_seq RESTART WITH 2;
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO sensor_data (sensor_id, data)
VALUES 
  (1, '{"end_device_ids": {"device_id": "Device123", "application_ids": {"application_id": "AppID123"}, "dev_eui": "00-14-22-01-23-45", "join_eui": "00-15-33-02-34-56"}, 
       "correlation_ids": ["corrID1", "corrID2"], 
       "received_at": "2023-10-01T12:00:00Z", 
       "uplink_message": {
            "session_key_id": "sessionKey1", 
            "f_port": 1, 
            "f_cnt": 10, 
            "frm_payload": "payloadData", 
            "decoded_payload": {
                "battery": 85.0, 
                "humidity": 55, 
                "raw": 123
            },
            "rx_metadata": [
                {
                    "gateway_ids": {"gateway_id": "Gateway123"}, 
                    "rssi": -45, 
                    "channel_rssi": -42, 
                    "snr": 9.5, 
                    "location": {"latitude": 52.5200, "longitude": 13.4050, "altitude": 34.0}
                }
            ],
            "settings": {
                "data_rate": {"lora": {"bandwidth": 125, "spreading_factor": 7, "coding_rate": "4/5"}}, 
                "frequency": "868100000"
            },
            "confirmed": true, 
            "consumed_airtime": "0.123s"
        }
    }');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sensors;
DELETE FROM sensor_data;
-- +goose StatementEnd