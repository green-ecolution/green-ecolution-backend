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

// checks if the RegionList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RegionList{}

// RegionList struct for RegionList
type RegionList struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Regions    []Region    `json:"regions"`
}

type _RegionList RegionList

// NewRegionList instantiates a new RegionList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegionList(regions []Region) *RegionList {
	this := RegionList{}
	this.Regions = regions
	return &this
}

// NewRegionListWithDefaults instantiates a new RegionList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegionListWithDefaults() *RegionList {
	this := RegionList{}
	return &this
}

// GetPagination returns the Pagination field value if set, zero value otherwise.
func (o *RegionList) GetPagination() Pagination {
	if o == nil || IsNil(o.Pagination) {
		var ret Pagination
		return ret
	}
	return *o.Pagination
}

// GetPaginationOk returns a tuple with the Pagination field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RegionList) GetPaginationOk() (*Pagination, bool) {
	if o == nil || IsNil(o.Pagination) {
		return nil, false
	}
	return o.Pagination, true
}

// HasPagination returns a boolean if a field has been set.
func (o *RegionList) HasPagination() bool {
	if o != nil && !IsNil(o.Pagination) {
		return true
	}

	return false
}

// SetPagination gets a reference to the given Pagination and assigns it to the Pagination field.
func (o *RegionList) SetPagination(v Pagination) {
	o.Pagination = &v
}

// GetRegions returns the Regions field value
func (o *RegionList) GetRegions() []Region {
	if o == nil {
		var ret []Region
		return ret
	}

	return o.Regions
}

// GetRegionsOk returns a tuple with the Regions field value
// and a boolean to check if the value has been set.
func (o *RegionList) GetRegionsOk() ([]Region, bool) {
	if o == nil {
		return nil, false
	}
	return o.Regions, true
}

// SetRegions sets field value
func (o *RegionList) SetRegions(v []Region) {
	o.Regions = v
}

func (o RegionList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RegionList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Pagination) {
		toSerialize["pagination"] = o.Pagination
	}
	toSerialize["regions"] = o.Regions
	return toSerialize, nil
}

func (o *RegionList) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"regions",
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

	varRegionList := _RegionList{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varRegionList)

	if err != nil {
		return err
	}

	*o = RegionList(varRegionList)

	return err
}

type NullableRegionList struct {
	value *RegionList
	isSet bool
}

func (v NullableRegionList) Get() *RegionList {
	return v.value
}

func (v *NullableRegionList) Set(val *RegionList) {
	v.value = val
	v.isSet = true
}

func (v NullableRegionList) IsSet() bool {
	return v.isSet
}

func (v *NullableRegionList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegionList(val *RegionList) *NullableRegionList {
	return &NullableRegionList{value: val, isSet: true}
}

func (v NullableRegionList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegionList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
