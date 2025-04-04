/*
Green Space Management API

This is the API for the Green Ecolution Management System.

API version: develop
Contact: info@green-ecolution.de
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the Sensor type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Sensor{}

// Sensor struct for Sensor
type Sensor struct {
	AdditionalInformation map[string]interface{} `json:"additional_information,omitempty"`
	CreatedAt             string                 `json:"created_at"`
	Id                    string                 `json:"id"`
	LatestData            SensorData             `json:"latest_data"`
	Latitude              float32                `json:"latitude"`
	Longitude             float32                `json:"longitude"`
	Provider              string                 `json:"provider"`
	Status                SensorStatus           `json:"status"`
	UpdatedAt             string                 `json:"updated_at"`
}

type _Sensor Sensor

// NewSensor instantiates a new Sensor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSensor(createdAt string, id string, latestData SensorData, latitude float32, longitude float32, provider string, status SensorStatus, updatedAt string) *Sensor {
	this := Sensor{}
	this.CreatedAt = createdAt
	this.Id = id
	this.LatestData = latestData
	this.Latitude = latitude
	this.Longitude = longitude
	this.Provider = provider
	this.Status = status
	this.UpdatedAt = updatedAt
	return &this
}

// NewSensorWithDefaults instantiates a new Sensor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSensorWithDefaults() *Sensor {
	this := Sensor{}
	return &this
}

// GetAdditionalInformation returns the AdditionalInformation field value if set, zero value otherwise.
func (o *Sensor) GetAdditionalInformation() map[string]interface{} {
	if o == nil || IsNil(o.AdditionalInformation) {
		var ret map[string]interface{}
		return ret
	}
	return o.AdditionalInformation
}

// GetAdditionalInformationOk returns a tuple with the AdditionalInformation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Sensor) GetAdditionalInformationOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.AdditionalInformation) {
		return map[string]interface{}{}, false
	}
	return o.AdditionalInformation, true
}

// HasAdditionalInformation returns a boolean if a field has been set.
func (o *Sensor) HasAdditionalInformation() bool {
	if o != nil && !IsNil(o.AdditionalInformation) {
		return true
	}

	return false
}

// SetAdditionalInformation gets a reference to the given map[string]interface{} and assigns it to the AdditionalInformation field.
func (o *Sensor) SetAdditionalInformation(v map[string]interface{}) {
	o.AdditionalInformation = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *Sensor) GetCreatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetCreatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *Sensor) SetCreatedAt(v string) {
	o.CreatedAt = v
}

// GetId returns the Id field value
func (o *Sensor) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Sensor) SetId(v string) {
	o.Id = v
}

// GetLatestData returns the LatestData field value
func (o *Sensor) GetLatestData() SensorData {
	if o == nil {
		var ret SensorData
		return ret
	}

	return o.LatestData
}

// GetLatestDataOk returns a tuple with the LatestData field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetLatestDataOk() (*SensorData, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LatestData, true
}

// SetLatestData sets field value
func (o *Sensor) SetLatestData(v SensorData) {
	o.LatestData = v
}

// GetLatitude returns the Latitude field value
func (o *Sensor) GetLatitude() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Latitude
}

// GetLatitudeOk returns a tuple with the Latitude field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetLatitudeOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Latitude, true
}

// SetLatitude sets field value
func (o *Sensor) SetLatitude(v float32) {
	o.Latitude = v
}

// GetLongitude returns the Longitude field value
func (o *Sensor) GetLongitude() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Longitude
}

// GetLongitudeOk returns a tuple with the Longitude field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetLongitudeOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Longitude, true
}

// SetLongitude sets field value
func (o *Sensor) SetLongitude(v float32) {
	o.Longitude = v
}

// GetProvider returns the Provider field value
func (o *Sensor) GetProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Provider
}

// GetProviderOk returns a tuple with the Provider field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Provider, true
}

// SetProvider sets field value
func (o *Sensor) SetProvider(v string) {
	o.Provider = v
}

// GetStatus returns the Status field value
func (o *Sensor) GetStatus() SensorStatus {
	if o == nil {
		var ret SensorStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetStatusOk() (*SensorStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *Sensor) SetStatus(v SensorStatus) {
	o.Status = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *Sensor) GetUpdatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *Sensor) GetUpdatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *Sensor) SetUpdatedAt(v string) {
	o.UpdatedAt = v
}

func (o Sensor) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Sensor) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AdditionalInformation) {
		toSerialize["additional_information"] = o.AdditionalInformation
	}
	toSerialize["created_at"] = o.CreatedAt
	toSerialize["id"] = o.Id
	toSerialize["latest_data"] = o.LatestData
	toSerialize["latitude"] = o.Latitude
	toSerialize["longitude"] = o.Longitude
	toSerialize["provider"] = o.Provider
	toSerialize["status"] = o.Status
	toSerialize["updated_at"] = o.UpdatedAt
	return toSerialize, nil
}

func (o *Sensor) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"created_at",
		"id",
		"latest_data",
		"latitude",
		"longitude",
		"provider",
		"status",
		"updated_at",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varSensor := _Sensor{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varSensor)

	if err != nil {
		return err
	}

	*o = Sensor(varSensor)

	return err
}

type NullableSensor struct {
	value *Sensor
	isSet bool
}

func (v NullableSensor) Get() *Sensor {
	return v.value
}

func (v *NullableSensor) Set(val *Sensor) {
	v.value = val
	v.isSet = true
}

func (v NullableSensor) IsSet() bool {
	return v.isSet
}

func (v *NullableSensor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSensor(val *Sensor) *NullableSensor {
	return &NullableSensor{value: val, isSet: true}
}

func (v NullableSensor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSensor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
