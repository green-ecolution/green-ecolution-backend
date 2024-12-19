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
	"fmt"
)

// VehicleStatus the model 'VehicleStatus'
type VehicleStatus string

// List of VehicleStatus
const (
	VehicleStatusActive       VehicleStatus = "active"
	VehicleStatusAvailable    VehicleStatus = "available"
	VehicleStatusNotAvailable VehicleStatus = "not available"
	VehicleStatusUnknown      VehicleStatus = "unknown"
)

// All allowed values of VehicleStatus enum
var AllowedVehicleStatusEnumValues = []VehicleStatus{
	"active",
	"available",
	"not available",
	"unknown",
}

func (v *VehicleStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := VehicleStatus(value)
	for _, existing := range AllowedVehicleStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid VehicleStatus", value)
}

// NewVehicleStatusFromValue returns a pointer to a valid VehicleStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewVehicleStatusFromValue(v string) (*VehicleStatus, error) {
	ev := VehicleStatus(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for VehicleStatus: valid values are %v", v, AllowedVehicleStatusEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v VehicleStatus) IsValid() bool {
	for _, existing := range AllowedVehicleStatusEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to VehicleStatus value
func (v VehicleStatus) Ptr() *VehicleStatus {
	return &v
}

type NullableVehicleStatus struct {
	value *VehicleStatus
	isSet bool
}

func (v NullableVehicleStatus) Get() *VehicleStatus {
	return v.value
}

func (v *NullableVehicleStatus) Set(val *VehicleStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableVehicleStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableVehicleStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVehicleStatus(val *VehicleStatus) *NullableVehicleStatus {
	return &NullableVehicleStatus{value: val, isSet: true}
}

func (v NullableVehicleStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVehicleStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
