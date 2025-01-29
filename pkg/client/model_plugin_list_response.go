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

// checks if the PluginListResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PluginListResponse{}

// PluginListResponse struct for PluginListResponse
type PluginListResponse struct {
	Plugins []Plugin `json:"plugins"`
}

type _PluginListResponse PluginListResponse

// NewPluginListResponse instantiates a new PluginListResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPluginListResponse(plugins []Plugin) *PluginListResponse {
	this := PluginListResponse{}
	this.Plugins = plugins
	return &this
}

// NewPluginListResponseWithDefaults instantiates a new PluginListResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPluginListResponseWithDefaults() *PluginListResponse {
	this := PluginListResponse{}
	return &this
}

// GetPlugins returns the Plugins field value
func (o *PluginListResponse) GetPlugins() []Plugin {
	if o == nil {
		var ret []Plugin
		return ret
	}

	return o.Plugins
}

// GetPluginsOk returns a tuple with the Plugins field value
// and a boolean to check if the value has been set.
func (o *PluginListResponse) GetPluginsOk() ([]Plugin, bool) {
	if o == nil {
		return nil, false
	}
	return o.Plugins, true
}

// SetPlugins sets field value
func (o *PluginListResponse) SetPlugins(v []Plugin) {
	o.Plugins = v
}

func (o PluginListResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PluginListResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["plugins"] = o.Plugins
	return toSerialize, nil
}

func (o *PluginListResponse) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"plugins",
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

	varPluginListResponse := _PluginListResponse{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varPluginListResponse)

	if err != nil {
		return err
	}

	*o = PluginListResponse(varPluginListResponse)

	return err
}

type NullablePluginListResponse struct {
	value *PluginListResponse
	isSet bool
}

func (v NullablePluginListResponse) Get() *PluginListResponse {
	return v.value
}

func (v *NullablePluginListResponse) Set(val *PluginListResponse) {
	v.value = val
	v.isSet = true
}

func (v NullablePluginListResponse) IsSet() bool {
	return v.isSet
}

func (v *NullablePluginListResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePluginListResponse(val *PluginListResponse) *NullablePluginListResponse {
	return &NullablePluginListResponse{value: val, isSet: true}
}

func (v NullablePluginListResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePluginListResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
