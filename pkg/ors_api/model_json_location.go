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

// checks if the JSONLocation type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JSONLocation{}

// JSONLocation The snapped locations as coordinates and snapping distance.
type JSONLocation struct {
	// {longitude},{latitude} coordinates of the closest accessible point on the routing graph
	Location []float64 `json:"location,omitempty"`
	// Name of the street the closest accessible point is situated on. Only for `resolve_locations=true` and only if name is available.
	Name *string `json:"name,omitempty"`
	// Distance between the `source/destination` Location and the used point on the routing graph in meters.
	SnappedDistance *float64 `json:"snapped_distance,omitempty"`
}

// NewJSONLocation instantiates a new JSONLocation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJSONLocation() *JSONLocation {
	this := JSONLocation{}
	return &this
}

// NewJSONLocationWithDefaults instantiates a new JSONLocation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJSONLocationWithDefaults() *JSONLocation {
	this := JSONLocation{}
	return &this
}

// GetLocation returns the Location field value if set, zero value otherwise.
func (o *JSONLocation) GetLocation() []float64 {
	if o == nil || IsNil(o.Location) {
		var ret []float64
		return ret
	}
	return o.Location
}

// GetLocationOk returns a tuple with the Location field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONLocation) GetLocationOk() ([]float64, bool) {
	if o == nil || IsNil(o.Location) {
		return nil, false
	}
	return o.Location, true
}

// HasLocation returns a boolean if a field has been set.
func (o *JSONLocation) HasLocation() bool {
	if o != nil && !IsNil(o.Location) {
		return true
	}

	return false
}

// SetLocation gets a reference to the given []float64 and assigns it to the Location field.
func (o *JSONLocation) SetLocation(v []float64) {
	o.Location = v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *JSONLocation) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONLocation) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *JSONLocation) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *JSONLocation) SetName(v string) {
	o.Name = &v
}

// GetSnappedDistance returns the SnappedDistance field value if set, zero value otherwise.
func (o *JSONLocation) GetSnappedDistance() float64 {
	if o == nil || IsNil(o.SnappedDistance) {
		var ret float64
		return ret
	}
	return *o.SnappedDistance
}

// GetSnappedDistanceOk returns a tuple with the SnappedDistance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONLocation) GetSnappedDistanceOk() (*float64, bool) {
	if o == nil || IsNil(o.SnappedDistance) {
		return nil, false
	}
	return o.SnappedDistance, true
}

// HasSnappedDistance returns a boolean if a field has been set.
func (o *JSONLocation) HasSnappedDistance() bool {
	if o != nil && !IsNil(o.SnappedDistance) {
		return true
	}

	return false
}

// SetSnappedDistance gets a reference to the given float64 and assigns it to the SnappedDistance field.
func (o *JSONLocation) SetSnappedDistance(v float64) {
	o.SnappedDistance = &v
}

func (o JSONLocation) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JSONLocation) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Location) {
		toSerialize["location"] = o.Location
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.SnappedDistance) {
		toSerialize["snapped_distance"] = o.SnappedDistance
	}
	return toSerialize, nil
}

type NullableJSONLocation struct {
	value *JSONLocation
	isSet bool
}

func (v NullableJSONLocation) Get() *JSONLocation {
	return v.value
}

func (v *NullableJSONLocation) Set(val *JSONLocation) {
	v.value = val
	v.isSet = true
}

func (v NullableJSONLocation) IsSet() bool {
	return v.isSet
}

func (v *NullableJSONLocation) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJSONLocation(val *JSONLocation) *NullableJSONLocation {
	return &NullableJSONLocation{value: val, isSet: true}
}

func (v NullableJSONLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJSONLocation) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


