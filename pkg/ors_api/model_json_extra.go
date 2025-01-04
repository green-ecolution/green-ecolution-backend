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

// checks if the JSONExtra type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JSONExtra{}

// JSONExtra An object representing one of the extra info items requested
type JSONExtra struct {
	// A list of values representing a section of the route. The individual values are:  Value 1: Indice of the staring point of the geometry for this section, Value 2: Indice of the end point of the geoemetry for this sections, Value 3: [Value](https://GIScience.github.io/openrouteservice/api-reference/endpoints/directions/extra-info/) assigned to this section.
	Values [][]int64 `json:"values,omitempty"`
	// List representing the summary of the extra info items.
	Summary []JSONExtraSummary `json:"summary,omitempty"`
}

// NewJSONExtra instantiates a new JSONExtra object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJSONExtra() *JSONExtra {
	this := JSONExtra{}
	return &this
}

// NewJSONExtraWithDefaults instantiates a new JSONExtra object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJSONExtraWithDefaults() *JSONExtra {
	this := JSONExtra{}
	return &this
}

// GetValues returns the Values field value if set, zero value otherwise.
func (o *JSONExtra) GetValues() [][]int64 {
	if o == nil || IsNil(o.Values) {
		var ret [][]int64
		return ret
	}
	return o.Values
}

// GetValuesOk returns a tuple with the Values field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONExtra) GetValuesOk() ([][]int64, bool) {
	if o == nil || IsNil(o.Values) {
		return nil, false
	}
	return o.Values, true
}

// HasValues returns a boolean if a field has been set.
func (o *JSONExtra) HasValues() bool {
	if o != nil && !IsNil(o.Values) {
		return true
	}

	return false
}

// SetValues gets a reference to the given [][]int64 and assigns it to the Values field.
func (o *JSONExtra) SetValues(v [][]int64) {
	o.Values = v
}

// GetSummary returns the Summary field value if set, zero value otherwise.
func (o *JSONExtra) GetSummary() []JSONExtraSummary {
	if o == nil || IsNil(o.Summary) {
		var ret []JSONExtraSummary
		return ret
	}
	return o.Summary
}

// GetSummaryOk returns a tuple with the Summary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONExtra) GetSummaryOk() ([]JSONExtraSummary, bool) {
	if o == nil || IsNil(o.Summary) {
		return nil, false
	}
	return o.Summary, true
}

// HasSummary returns a boolean if a field has been set.
func (o *JSONExtra) HasSummary() bool {
	if o != nil && !IsNil(o.Summary) {
		return true
	}

	return false
}

// SetSummary gets a reference to the given []JSONExtraSummary and assigns it to the Summary field.
func (o *JSONExtra) SetSummary(v []JSONExtraSummary) {
	o.Summary = v
}

func (o JSONExtra) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JSONExtra) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Values) {
		toSerialize["values"] = o.Values
	}
	if !IsNil(o.Summary) {
		toSerialize["summary"] = o.Summary
	}
	return toSerialize, nil
}

type NullableJSONExtra struct {
	value *JSONExtra
	isSet bool
}

func (v NullableJSONExtra) Get() *JSONExtra {
	return v.value
}

func (v *NullableJSONExtra) Set(val *JSONExtra) {
	v.value = val
	v.isSet = true
}

func (v NullableJSONExtra) IsSet() bool {
	return v.isSet
}

func (v *NullableJSONExtra) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJSONExtra(val *JSONExtra) *NullableJSONExtra {
	return &NullableJSONExtra{value: val, isSet: true}
}

func (v NullableJSONExtra) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJSONExtra) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


