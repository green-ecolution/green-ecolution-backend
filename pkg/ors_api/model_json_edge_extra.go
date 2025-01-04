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

// checks if the JsonEdgeExtra type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JsonEdgeExtra{}

// JsonEdgeExtra struct for JsonEdgeExtra
type JsonEdgeExtra struct {
	// Id of the corresponding edge in the graph
	EdgeId *string `json:"edgeId,omitempty"`
	// Extra info stored on the edge
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// NewJsonEdgeExtra instantiates a new JsonEdgeExtra object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJsonEdgeExtra() *JsonEdgeExtra {
	this := JsonEdgeExtra{}
	return &this
}

// NewJsonEdgeExtraWithDefaults instantiates a new JsonEdgeExtra object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJsonEdgeExtraWithDefaults() *JsonEdgeExtra {
	this := JsonEdgeExtra{}
	return &this
}

// GetEdgeId returns the EdgeId field value if set, zero value otherwise.
func (o *JsonEdgeExtra) GetEdgeId() string {
	if o == nil || IsNil(o.EdgeId) {
		var ret string
		return ret
	}
	return *o.EdgeId
}

// GetEdgeIdOk returns a tuple with the EdgeId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JsonEdgeExtra) GetEdgeIdOk() (*string, bool) {
	if o == nil || IsNil(o.EdgeId) {
		return nil, false
	}
	return o.EdgeId, true
}

// HasEdgeId returns a boolean if a field has been set.
func (o *JsonEdgeExtra) HasEdgeId() bool {
	if o != nil && !IsNil(o.EdgeId) {
		return true
	}

	return false
}

// SetEdgeId gets a reference to the given string and assigns it to the EdgeId field.
func (o *JsonEdgeExtra) SetEdgeId(v string) {
	o.EdgeId = &v
}

// GetExtra returns the Extra field value if set, zero value otherwise.
func (o *JsonEdgeExtra) GetExtra() map[string]interface{} {
	if o == nil || IsNil(o.Extra) {
		var ret map[string]interface{}
		return ret
	}
	return o.Extra
}

// GetExtraOk returns a tuple with the Extra field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JsonEdgeExtra) GetExtraOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Extra) {
		return map[string]interface{}{}, false
	}
	return o.Extra, true
}

// HasExtra returns a boolean if a field has been set.
func (o *JsonEdgeExtra) HasExtra() bool {
	if o != nil && !IsNil(o.Extra) {
		return true
	}

	return false
}

// SetExtra gets a reference to the given map[string]interface{} and assigns it to the Extra field.
func (o *JsonEdgeExtra) SetExtra(v map[string]interface{}) {
	o.Extra = v
}

func (o JsonEdgeExtra) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JsonEdgeExtra) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EdgeId) {
		toSerialize["edgeId"] = o.EdgeId
	}
	if !IsNil(o.Extra) {
		toSerialize["extra"] = o.Extra
	}
	return toSerialize, nil
}

type NullableJsonEdgeExtra struct {
	value *JsonEdgeExtra
	isSet bool
}

func (v NullableJsonEdgeExtra) Get() *JsonEdgeExtra {
	return v.value
}

func (v *NullableJsonEdgeExtra) Set(val *JsonEdgeExtra) {
	v.value = val
	v.isSet = true
}

func (v NullableJsonEdgeExtra) IsSet() bool {
	return v.isSet
}

func (v *NullableJsonEdgeExtra) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJsonEdgeExtra(val *JsonEdgeExtra) *NullableJsonEdgeExtra {
	return &NullableJsonEdgeExtra{value: val, isSet: true}
}

func (v NullableJsonEdgeExtra) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJsonEdgeExtra) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


