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

// checks if the Pagination type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Pagination{}

// Pagination struct for Pagination
type Pagination struct {
	CurrentPage  int32 `json:"current_page"`
	NextPage     int32 `json:"next_page"`
	PrevPage     int32 `json:"prev_page"`
	TotalPages   int32 `json:"total_pages"`
	TotalRecords int32 `json:"total_records"`
}

type _Pagination Pagination

// NewPagination instantiates a new Pagination object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPagination(currentPage int32, nextPage int32, prevPage int32, totalPages int32, totalRecords int32) *Pagination {
	this := Pagination{}
	this.CurrentPage = currentPage
	this.NextPage = nextPage
	this.PrevPage = prevPage
	this.TotalPages = totalPages
	this.TotalRecords = totalRecords
	return &this
}

// NewPaginationWithDefaults instantiates a new Pagination object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginationWithDefaults() *Pagination {
	this := Pagination{}
	return &this
}

// GetCurrentPage returns the CurrentPage field value
func (o *Pagination) GetCurrentPage() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.CurrentPage
}

// GetCurrentPageOk returns a tuple with the CurrentPage field value
// and a boolean to check if the value has been set.
func (o *Pagination) GetCurrentPageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CurrentPage, true
}

// SetCurrentPage sets field value
func (o *Pagination) SetCurrentPage(v int32) {
	o.CurrentPage = v
}

// GetNextPage returns the NextPage field value
func (o *Pagination) GetNextPage() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NextPage
}

// GetNextPageOk returns a tuple with the NextPage field value
// and a boolean to check if the value has been set.
func (o *Pagination) GetNextPageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NextPage, true
}

// SetNextPage sets field value
func (o *Pagination) SetNextPage(v int32) {
	o.NextPage = v
}

// GetPrevPage returns the PrevPage field value
func (o *Pagination) GetPrevPage() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.PrevPage
}

// GetPrevPageOk returns a tuple with the PrevPage field value
// and a boolean to check if the value has been set.
func (o *Pagination) GetPrevPageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PrevPage, true
}

// SetPrevPage sets field value
func (o *Pagination) SetPrevPage(v int32) {
	o.PrevPage = v
}

// GetTotalPages returns the TotalPages field value
func (o *Pagination) GetTotalPages() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TotalPages
}

// GetTotalPagesOk returns a tuple with the TotalPages field value
// and a boolean to check if the value has been set.
func (o *Pagination) GetTotalPagesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalPages, true
}

// SetTotalPages sets field value
func (o *Pagination) SetTotalPages(v int32) {
	o.TotalPages = v
}

// GetTotalRecords returns the TotalRecords field value
func (o *Pagination) GetTotalRecords() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TotalRecords
}

// GetTotalRecordsOk returns a tuple with the TotalRecords field value
// and a boolean to check if the value has been set.
func (o *Pagination) GetTotalRecordsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TotalRecords, true
}

// SetTotalRecords sets field value
func (o *Pagination) SetTotalRecords(v int32) {
	o.TotalRecords = v
}

func (o Pagination) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Pagination) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["current_page"] = o.CurrentPage
	toSerialize["next_page"] = o.NextPage
	toSerialize["prev_page"] = o.PrevPage
	toSerialize["total_pages"] = o.TotalPages
	toSerialize["total_records"] = o.TotalRecords
	return toSerialize, nil
}

func (o *Pagination) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"current_page",
		"next_page",
		"prev_page",
		"total_pages",
		"total_records",
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

	varPagination := _Pagination{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varPagination)

	if err != nil {
		return err
	}

	*o = Pagination(varPagination)

	return err
}

type NullablePagination struct {
	value *Pagination
	isSet bool
}

func (v NullablePagination) Get() *Pagination {
	return v.value
}

func (v *NullablePagination) Set(val *Pagination) {
	v.value = val
	v.isSet = true
}

func (v NullablePagination) IsSet() bool {
	return v.isSet
}

func (v *NullablePagination) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePagination(val *Pagination) *NullablePagination {
	return &NullablePagination{value: val, isSet: true}
}

func (v NullablePagination) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePagination) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
