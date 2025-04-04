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

// checks if the ServerInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ServerInfo{}

// ServerInfo struct for ServerInfo
type ServerInfo struct {
	Arch      string `json:"arch"`
	Hostname  string `json:"hostname"`
	Interface string `json:"interface"`
	Ip        string `json:"ip"`
	Os        string `json:"os"`
	Port      int32  `json:"port"`
	Uptime    string `json:"uptime"`
	Url       string `json:"url"`
}

type _ServerInfo ServerInfo

// NewServerInfo instantiates a new ServerInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerInfo(arch string, hostname string, interface_ string, ip string, os string, port int32, uptime string, url string) *ServerInfo {
	this := ServerInfo{}
	this.Arch = arch
	this.Hostname = hostname
	this.Interface = interface_
	this.Ip = ip
	this.Os = os
	this.Port = port
	this.Uptime = uptime
	this.Url = url
	return &this
}

// NewServerInfoWithDefaults instantiates a new ServerInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerInfoWithDefaults() *ServerInfo {
	this := ServerInfo{}
	return &this
}

// GetArch returns the Arch field value
func (o *ServerInfo) GetArch() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Arch
}

// GetArchOk returns a tuple with the Arch field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetArchOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Arch, true
}

// SetArch sets field value
func (o *ServerInfo) SetArch(v string) {
	o.Arch = v
}

// GetHostname returns the Hostname field value
func (o *ServerInfo) GetHostname() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetHostnameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Hostname, true
}

// SetHostname sets field value
func (o *ServerInfo) SetHostname(v string) {
	o.Hostname = v
}

// GetInterface returns the Interface field value
func (o *ServerInfo) GetInterface() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Interface
}

// GetInterfaceOk returns a tuple with the Interface field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetInterfaceOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Interface, true
}

// SetInterface sets field value
func (o *ServerInfo) SetInterface(v string) {
	o.Interface = v
}

// GetIp returns the Ip field value
func (o *ServerInfo) GetIp() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Ip
}

// GetIpOk returns a tuple with the Ip field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetIpOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ip, true
}

// SetIp sets field value
func (o *ServerInfo) SetIp(v string) {
	o.Ip = v
}

// GetOs returns the Os field value
func (o *ServerInfo) GetOs() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Os
}

// GetOsOk returns a tuple with the Os field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetOsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Os, true
}

// SetOs sets field value
func (o *ServerInfo) SetOs(v string) {
	o.Os = v
}

// GetPort returns the Port field value
func (o *ServerInfo) GetPort() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Port
}

// GetPortOk returns a tuple with the Port field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetPortOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Port, true
}

// SetPort sets field value
func (o *ServerInfo) SetPort(v int32) {
	o.Port = v
}

// GetUptime returns the Uptime field value
func (o *ServerInfo) GetUptime() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Uptime
}

// GetUptimeOk returns a tuple with the Uptime field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetUptimeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uptime, true
}

// SetUptime sets field value
func (o *ServerInfo) SetUptime(v string) {
	o.Uptime = v
}

// GetUrl returns the Url field value
func (o *ServerInfo) GetUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Url
}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
func (o *ServerInfo) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Url, true
}

// SetUrl sets field value
func (o *ServerInfo) SetUrl(v string) {
	o.Url = v
}

func (o ServerInfo) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ServerInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["arch"] = o.Arch
	toSerialize["hostname"] = o.Hostname
	toSerialize["interface"] = o.Interface
	toSerialize["ip"] = o.Ip
	toSerialize["os"] = o.Os
	toSerialize["port"] = o.Port
	toSerialize["uptime"] = o.Uptime
	toSerialize["url"] = o.Url
	return toSerialize, nil
}

func (o *ServerInfo) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"arch",
		"hostname",
		"interface",
		"ip",
		"os",
		"port",
		"uptime",
		"url",
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

	varServerInfo := _ServerInfo{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varServerInfo)

	if err != nil {
		return err
	}

	*o = ServerInfo(varServerInfo)

	return err
}

type NullableServerInfo struct {
	value *ServerInfo
	isSet bool
}

func (v NullableServerInfo) Get() *ServerInfo {
	return v.value
}

func (v *NullableServerInfo) Set(val *ServerInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableServerInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableServerInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServerInfo(val *ServerInfo) *NullableServerInfo {
	return &NullableServerInfo{value: val, isSet: true}
}

func (v NullableServerInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServerInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
