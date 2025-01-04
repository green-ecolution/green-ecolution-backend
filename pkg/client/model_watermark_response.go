/*
Green Space Management API

This is the API for the Green Ecolution Management System.

API version: develop
Contact: info@green-ecolution.de
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the WatermarkResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WatermarkResponse{}

// WatermarkResponse struct for WatermarkResponse
type WatermarkResponse struct {
	Centibar int32 `json:"centibar"`
	Depth int32 `json:"depth"`
	Resistance int32 `json:"resistance"`
}

type _WatermarkResponse WatermarkResponse

// NewWatermarkResponse instantiates a new WatermarkResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWatermarkResponse(centibar int32, depth int32, resistance int32) *WatermarkResponse {
	this := WatermarkResponse{}
	this.Centibar = centibar
	this.Depth = depth
	this.Resistance = resistance
	return &this
}

// NewWatermarkResponseWithDefaults instantiates a new WatermarkResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWatermarkResponseWithDefaults() *WatermarkResponse {
	this := WatermarkResponse{}
	return &this
}

// GetCentibar returns the Centibar field value
func (o *WatermarkResponse) GetCentibar() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Centibar
}

// GetCentibarOk returns a tuple with the Centibar field value
// and a boolean to check if the value has been set.
func (o *WatermarkResponse) GetCentibarOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Centibar, true
}

// SetCentibar sets field value
func (o *WatermarkResponse) SetCentibar(v int32) {
	o.Centibar = v
}

// GetDepth returns the Depth field value
func (o *WatermarkResponse) GetDepth() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Depth
}

// GetDepthOk returns a tuple with the Depth field value
// and a boolean to check if the value has been set.
func (o *WatermarkResponse) GetDepthOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Depth, true
}

// SetDepth sets field value
func (o *WatermarkResponse) SetDepth(v int32) {
	o.Depth = v
}

// GetResistance returns the Resistance field value
func (o *WatermarkResponse) GetResistance() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Resistance
}

// GetResistanceOk returns a tuple with the Resistance field value
// and a boolean to check if the value has been set.
func (o *WatermarkResponse) GetResistanceOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Resistance, true
}

// SetResistance sets field value
func (o *WatermarkResponse) SetResistance(v int32) {
	o.Resistance = v
}

func (o WatermarkResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WatermarkResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["centibar"] = o.Centibar
	toSerialize["depth"] = o.Depth
	toSerialize["resistance"] = o.Resistance
	return toSerialize, nil
}

func (o *WatermarkResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"centibar",
		"depth",
		"resistance",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varWatermarkResponse := _WatermarkResponse{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varWatermarkResponse)

	if err != nil {
		return err
	}

	*o = WatermarkResponse(varWatermarkResponse)

	return err
}

type NullableWatermarkResponse struct {
	value *WatermarkResponse
	isSet bool
}

func (v NullableWatermarkResponse) Get() *WatermarkResponse {
	return v.value
}

func (v *NullableWatermarkResponse) Set(val *WatermarkResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWatermarkResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWatermarkResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWatermarkResponse(val *WatermarkResponse) *NullableWatermarkResponse {
	return &NullableWatermarkResponse{value: val, isSet: true}
}

func (v NullableWatermarkResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWatermarkResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


