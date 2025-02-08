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

// checks if the GeoJsonGeometry type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GeoJsonGeometry{}

// GeoJsonGeometry struct for GeoJsonGeometry
type GeoJsonGeometry struct {
	Coordinates [][]float32 `json:"coordinates"`
	Type        GeoJsonType `json:"type"`
}

type _GeoJsonGeometry GeoJsonGeometry

// NewGeoJsonGeometry instantiates a new GeoJsonGeometry object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGeoJsonGeometry(coordinates [][]float32, type_ GeoJsonType) *GeoJsonGeometry {
	this := GeoJsonGeometry{}
	this.Coordinates = coordinates
	this.Type = type_
	return &this
}

// NewGeoJsonGeometryWithDefaults instantiates a new GeoJsonGeometry object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGeoJsonGeometryWithDefaults() *GeoJsonGeometry {
	this := GeoJsonGeometry{}
	return &this
}

// GetCoordinates returns the Coordinates field value
func (o *GeoJsonGeometry) GetCoordinates() [][]float32 {
	if o == nil {
		var ret [][]float32
		return ret
	}

	return o.Coordinates
}

// GetCoordinatesOk returns a tuple with the Coordinates field value
// and a boolean to check if the value has been set.
func (o *GeoJsonGeometry) GetCoordinatesOk() ([][]float32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Coordinates, true
}

// SetCoordinates sets field value
func (o *GeoJsonGeometry) SetCoordinates(v [][]float32) {
	o.Coordinates = v
}

// GetType returns the Type field value
func (o *GeoJsonGeometry) GetType() GeoJsonType {
	if o == nil {
		var ret GeoJsonType
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *GeoJsonGeometry) GetTypeOk() (*GeoJsonType, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *GeoJsonGeometry) SetType(v GeoJsonType) {
	o.Type = v
}

func (o GeoJsonGeometry) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GeoJsonGeometry) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["coordinates"] = o.Coordinates
	toSerialize["type"] = o.Type
	return toSerialize, nil
}

func (o *GeoJsonGeometry) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"coordinates",
		"type",
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

	varGeoJsonGeometry := _GeoJsonGeometry{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varGeoJsonGeometry)

	if err != nil {
		return err
	}

	*o = GeoJsonGeometry(varGeoJsonGeometry)

	return err
}

type NullableGeoJsonGeometry struct {
	value *GeoJsonGeometry
	isSet bool
}

func (v NullableGeoJsonGeometry) Get() *GeoJsonGeometry {
	return v.value
}

func (v *NullableGeoJsonGeometry) Set(val *GeoJsonGeometry) {
	v.value = val
	v.isSet = true
}

func (v NullableGeoJsonGeometry) IsSet() bool {
	return v.isSet
}

func (v *NullableGeoJsonGeometry) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGeoJsonGeometry(val *GeoJsonGeometry) *NullableGeoJsonGeometry {
	return &NullableGeoJsonGeometry{value: val, isSet: true}
}

func (v NullableGeoJsonGeometry) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGeoJsonGeometry) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
