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

// checks if the RouteRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RouteRequest{}

// RouteRequest struct for RouteRequest
type RouteRequest struct {
	ClusterIds    []int32 `json:"cluster_ids"`
	TrailerId     *int32  `json:"trailer_id,omitempty"`
	TransporterId int32   `json:"transporter_id"`
}

type _RouteRequest RouteRequest

// NewRouteRequest instantiates a new RouteRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRouteRequest(clusterIds []int32, transporterId int32) *RouteRequest {
	this := RouteRequest{}
	this.ClusterIds = clusterIds
	this.TransporterId = transporterId
	return &this
}

// NewRouteRequestWithDefaults instantiates a new RouteRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRouteRequestWithDefaults() *RouteRequest {
	this := RouteRequest{}
	return &this
}

// GetClusterIds returns the ClusterIds field value
func (o *RouteRequest) GetClusterIds() []int32 {
	if o == nil {
		var ret []int32
		return ret
	}

	return o.ClusterIds
}

// GetClusterIdsOk returns a tuple with the ClusterIds field value
// and a boolean to check if the value has been set.
func (o *RouteRequest) GetClusterIdsOk() ([]int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.ClusterIds, true
}

// SetClusterIds sets field value
func (o *RouteRequest) SetClusterIds(v []int32) {
	o.ClusterIds = v
}

// GetTrailerId returns the TrailerId field value if set, zero value otherwise.
func (o *RouteRequest) GetTrailerId() int32 {
	if o == nil || IsNil(o.TrailerId) {
		var ret int32
		return ret
	}
	return *o.TrailerId
}

// GetTrailerIdOk returns a tuple with the TrailerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RouteRequest) GetTrailerIdOk() (*int32, bool) {
	if o == nil || IsNil(o.TrailerId) {
		return nil, false
	}
	return o.TrailerId, true
}

// HasTrailerId returns a boolean if a field has been set.
func (o *RouteRequest) HasTrailerId() bool {
	if o != nil && !IsNil(o.TrailerId) {
		return true
	}

	return false
}

// SetTrailerId gets a reference to the given int32 and assigns it to the TrailerId field.
func (o *RouteRequest) SetTrailerId(v int32) {
	o.TrailerId = &v
}

// GetTransporterId returns the TransporterId field value
func (o *RouteRequest) GetTransporterId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TransporterId
}

// GetTransporterIdOk returns a tuple with the TransporterId field value
// and a boolean to check if the value has been set.
func (o *RouteRequest) GetTransporterIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TransporterId, true
}

// SetTransporterId sets field value
func (o *RouteRequest) SetTransporterId(v int32) {
	o.TransporterId = v
}

func (o RouteRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RouteRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["cluster_ids"] = o.ClusterIds
	if !IsNil(o.TrailerId) {
		toSerialize["trailer_id"] = o.TrailerId
	}
	toSerialize["transporter_id"] = o.TransporterId
	return toSerialize, nil
}

func (o *RouteRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"cluster_ids",
		"transporter_id",
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

	varRouteRequest := _RouteRequest{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varRouteRequest)

	if err != nil {
		return err
	}

	*o = RouteRequest(varRouteRequest)

	return err
}

type NullableRouteRequest struct {
	value *RouteRequest
	isSet bool
}

func (v NullableRouteRequest) Get() *RouteRequest {
	return v.value
}

func (v *NullableRouteRequest) Set(val *RouteRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableRouteRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableRouteRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRouteRequest(val *RouteRequest) *NullableRouteRequest {
	return &NullableRouteRequest{value: val, isSet: true}
}

func (v NullableRouteRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRouteRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
