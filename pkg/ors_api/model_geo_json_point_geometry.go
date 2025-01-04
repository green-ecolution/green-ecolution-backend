/*
Openrouteservice

This is the openrouteservice API documentation for ORS Core-Version 9.0.0. Documentations for [older Core-Versions](https://github.com/GIScience/openrouteservice-docs/releases) can be rendered with the [Swagger-Editor](https://editor-next.swagger.io/).

API version: v2
Contact: support@smartmobility.heigit.org
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ors

import (
	"encoding/json"
)

// checks if the GeoJSONPointGeometry type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GeoJSONPointGeometry{}

// GeoJSONPointGeometry Feature geometry
type GeoJSONPointGeometry struct {
	// GeoJSON type
	Type *string `json:"type,omitempty"`
	// Lon/Lat coordinates of the snapped location
	Coordinates []float64 `json:"coordinates,omitempty"`
}

// NewGeoJSONPointGeometry instantiates a new GeoJSONPointGeometry object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGeoJSONPointGeometry() *GeoJSONPointGeometry {
	this := GeoJSONPointGeometry{}
	var type_ string = "Point"
	this.Type = &type_
	return &this
}

// NewGeoJSONPointGeometryWithDefaults instantiates a new GeoJSONPointGeometry object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGeoJSONPointGeometryWithDefaults() *GeoJSONPointGeometry {
	this := GeoJSONPointGeometry{}
	var type_ string = "Point"
	this.Type = &type_
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *GeoJSONPointGeometry) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeoJSONPointGeometry) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *GeoJSONPointGeometry) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *GeoJSONPointGeometry) SetType(v string) {
	o.Type = &v
}

// GetCoordinates returns the Coordinates field value if set, zero value otherwise.
func (o *GeoJSONPointGeometry) GetCoordinates() []float64 {
	if o == nil || IsNil(o.Coordinates) {
		var ret []float64
		return ret
	}
	return o.Coordinates
}

// GetCoordinatesOk returns a tuple with the Coordinates field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeoJSONPointGeometry) GetCoordinatesOk() ([]float64, bool) {
	if o == nil || IsNil(o.Coordinates) {
		return nil, false
	}
	return o.Coordinates, true
}

// HasCoordinates returns a boolean if a field has been set.
func (o *GeoJSONPointGeometry) HasCoordinates() bool {
	if o != nil && !IsNil(o.Coordinates) {
		return true
	}

	return false
}

// SetCoordinates gets a reference to the given []float64 and assigns it to the Coordinates field.
func (o *GeoJSONPointGeometry) SetCoordinates(v []float64) {
	o.Coordinates = v
}

func (o GeoJSONPointGeometry) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GeoJSONPointGeometry) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Coordinates) {
		toSerialize["coordinates"] = o.Coordinates
	}
	return toSerialize, nil
}

type NullableGeoJSONPointGeometry struct {
	value *GeoJSONPointGeometry
	isSet bool
}

func (v NullableGeoJSONPointGeometry) Get() *GeoJSONPointGeometry {
	return v.value
}

func (v *NullableGeoJSONPointGeometry) Set(val *GeoJSONPointGeometry) {
	v.value = val
	v.isSet = true
}

func (v NullableGeoJSONPointGeometry) IsSet() bool {
	return v.isSet
}

func (v *NullableGeoJSONPointGeometry) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGeoJSONPointGeometry(val *GeoJSONPointGeometry) *NullableGeoJSONPointGeometry {
	return &NullableGeoJSONPointGeometry{value: val, isSet: true}
}

func (v NullableGeoJSONPointGeometry) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGeoJSONPointGeometry) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


