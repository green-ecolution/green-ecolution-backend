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
	"time"
)

// checks if the JSONPtStop type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &JSONPtStop{}

// JSONPtStop Stop of a public transport leg
type JSONPtStop struct {
	// The ID of the stop.
	StopId *string `json:"stop_id,omitempty"`
	// The name of the stop.
	Name *string `json:"name,omitempty"`
	// The location of the stop.
	Location []float64 `json:"location,omitempty"`
	// Arrival time of the stop.
	ArrivalTime *time.Time `json:"arrival_time,omitempty"`
	// Planned arrival time of the stop.
	PlannedArrivalTime *time.Time `json:"planned_arrival_time,omitempty"`
	// Predicted arrival time of the stop.
	PredictedArrivalTime *time.Time `json:"predicted_arrival_time,omitempty"`
	// Whether arrival at the stop was cancelled.
	ArrivalCancelled *bool `json:"arrival_cancelled,omitempty"`
	// Departure time of the stop.
	DepartureTime *time.Time `json:"departure_time,omitempty"`
	// Planned departure time of the stop.
	PlannedDepartureTime *time.Time `json:"planned_departure_time,omitempty"`
	// Predicted departure time of the stop.
	PredictedDepartureTime *time.Time `json:"predicted_departure_time,omitempty"`
	// Whether departure at the stop was cancelled.
	DepartureCancelled *bool `json:"departure_cancelled,omitempty"`
}

// NewJSONPtStop instantiates a new JSONPtStop object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewJSONPtStop() *JSONPtStop {
	this := JSONPtStop{}
	return &this
}

// NewJSONPtStopWithDefaults instantiates a new JSONPtStop object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewJSONPtStopWithDefaults() *JSONPtStop {
	this := JSONPtStop{}
	return &this
}

// GetStopId returns the StopId field value if set, zero value otherwise.
func (o *JSONPtStop) GetStopId() string {
	if o == nil || IsNil(o.StopId) {
		var ret string
		return ret
	}
	return *o.StopId
}

// GetStopIdOk returns a tuple with the StopId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetStopIdOk() (*string, bool) {
	if o == nil || IsNil(o.StopId) {
		return nil, false
	}
	return o.StopId, true
}

// HasStopId returns a boolean if a field has been set.
func (o *JSONPtStop) HasStopId() bool {
	if o != nil && !IsNil(o.StopId) {
		return true
	}

	return false
}

// SetStopId gets a reference to the given string and assigns it to the StopId field.
func (o *JSONPtStop) SetStopId(v string) {
	o.StopId = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *JSONPtStop) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *JSONPtStop) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *JSONPtStop) SetName(v string) {
	o.Name = &v
}

// GetLocation returns the Location field value if set, zero value otherwise.
func (o *JSONPtStop) GetLocation() []float64 {
	if o == nil || IsNil(o.Location) {
		var ret []float64
		return ret
	}
	return o.Location
}

// GetLocationOk returns a tuple with the Location field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetLocationOk() ([]float64, bool) {
	if o == nil || IsNil(o.Location) {
		return nil, false
	}
	return o.Location, true
}

// HasLocation returns a boolean if a field has been set.
func (o *JSONPtStop) HasLocation() bool {
	if o != nil && !IsNil(o.Location) {
		return true
	}

	return false
}

// SetLocation gets a reference to the given []float64 and assigns it to the Location field.
func (o *JSONPtStop) SetLocation(v []float64) {
	o.Location = v
}

// GetArrivalTime returns the ArrivalTime field value if set, zero value otherwise.
func (o *JSONPtStop) GetArrivalTime() time.Time {
	if o == nil || IsNil(o.ArrivalTime) {
		var ret time.Time
		return ret
	}
	return *o.ArrivalTime
}

// GetArrivalTimeOk returns a tuple with the ArrivalTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetArrivalTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ArrivalTime) {
		return nil, false
	}
	return o.ArrivalTime, true
}

// HasArrivalTime returns a boolean if a field has been set.
func (o *JSONPtStop) HasArrivalTime() bool {
	if o != nil && !IsNil(o.ArrivalTime) {
		return true
	}

	return false
}

// SetArrivalTime gets a reference to the given time.Time and assigns it to the ArrivalTime field.
func (o *JSONPtStop) SetArrivalTime(v time.Time) {
	o.ArrivalTime = &v
}

// GetPlannedArrivalTime returns the PlannedArrivalTime field value if set, zero value otherwise.
func (o *JSONPtStop) GetPlannedArrivalTime() time.Time {
	if o == nil || IsNil(o.PlannedArrivalTime) {
		var ret time.Time
		return ret
	}
	return *o.PlannedArrivalTime
}

// GetPlannedArrivalTimeOk returns a tuple with the PlannedArrivalTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetPlannedArrivalTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PlannedArrivalTime) {
		return nil, false
	}
	return o.PlannedArrivalTime, true
}

// HasPlannedArrivalTime returns a boolean if a field has been set.
func (o *JSONPtStop) HasPlannedArrivalTime() bool {
	if o != nil && !IsNil(o.PlannedArrivalTime) {
		return true
	}

	return false
}

// SetPlannedArrivalTime gets a reference to the given time.Time and assigns it to the PlannedArrivalTime field.
func (o *JSONPtStop) SetPlannedArrivalTime(v time.Time) {
	o.PlannedArrivalTime = &v
}

// GetPredictedArrivalTime returns the PredictedArrivalTime field value if set, zero value otherwise.
func (o *JSONPtStop) GetPredictedArrivalTime() time.Time {
	if o == nil || IsNil(o.PredictedArrivalTime) {
		var ret time.Time
		return ret
	}
	return *o.PredictedArrivalTime
}

// GetPredictedArrivalTimeOk returns a tuple with the PredictedArrivalTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetPredictedArrivalTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PredictedArrivalTime) {
		return nil, false
	}
	return o.PredictedArrivalTime, true
}

// HasPredictedArrivalTime returns a boolean if a field has been set.
func (o *JSONPtStop) HasPredictedArrivalTime() bool {
	if o != nil && !IsNil(o.PredictedArrivalTime) {
		return true
	}

	return false
}

// SetPredictedArrivalTime gets a reference to the given time.Time and assigns it to the PredictedArrivalTime field.
func (o *JSONPtStop) SetPredictedArrivalTime(v time.Time) {
	o.PredictedArrivalTime = &v
}

// GetArrivalCancelled returns the ArrivalCancelled field value if set, zero value otherwise.
func (o *JSONPtStop) GetArrivalCancelled() bool {
	if o == nil || IsNil(o.ArrivalCancelled) {
		var ret bool
		return ret
	}
	return *o.ArrivalCancelled
}

// GetArrivalCancelledOk returns a tuple with the ArrivalCancelled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetArrivalCancelledOk() (*bool, bool) {
	if o == nil || IsNil(o.ArrivalCancelled) {
		return nil, false
	}
	return o.ArrivalCancelled, true
}

// HasArrivalCancelled returns a boolean if a field has been set.
func (o *JSONPtStop) HasArrivalCancelled() bool {
	if o != nil && !IsNil(o.ArrivalCancelled) {
		return true
	}

	return false
}

// SetArrivalCancelled gets a reference to the given bool and assigns it to the ArrivalCancelled field.
func (o *JSONPtStop) SetArrivalCancelled(v bool) {
	o.ArrivalCancelled = &v
}

// GetDepartureTime returns the DepartureTime field value if set, zero value otherwise.
func (o *JSONPtStop) GetDepartureTime() time.Time {
	if o == nil || IsNil(o.DepartureTime) {
		var ret time.Time
		return ret
	}
	return *o.DepartureTime
}

// GetDepartureTimeOk returns a tuple with the DepartureTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetDepartureTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DepartureTime) {
		return nil, false
	}
	return o.DepartureTime, true
}

// HasDepartureTime returns a boolean if a field has been set.
func (o *JSONPtStop) HasDepartureTime() bool {
	if o != nil && !IsNil(o.DepartureTime) {
		return true
	}

	return false
}

// SetDepartureTime gets a reference to the given time.Time and assigns it to the DepartureTime field.
func (o *JSONPtStop) SetDepartureTime(v time.Time) {
	o.DepartureTime = &v
}

// GetPlannedDepartureTime returns the PlannedDepartureTime field value if set, zero value otherwise.
func (o *JSONPtStop) GetPlannedDepartureTime() time.Time {
	if o == nil || IsNil(o.PlannedDepartureTime) {
		var ret time.Time
		return ret
	}
	return *o.PlannedDepartureTime
}

// GetPlannedDepartureTimeOk returns a tuple with the PlannedDepartureTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetPlannedDepartureTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PlannedDepartureTime) {
		return nil, false
	}
	return o.PlannedDepartureTime, true
}

// HasPlannedDepartureTime returns a boolean if a field has been set.
func (o *JSONPtStop) HasPlannedDepartureTime() bool {
	if o != nil && !IsNil(o.PlannedDepartureTime) {
		return true
	}

	return false
}

// SetPlannedDepartureTime gets a reference to the given time.Time and assigns it to the PlannedDepartureTime field.
func (o *JSONPtStop) SetPlannedDepartureTime(v time.Time) {
	o.PlannedDepartureTime = &v
}

// GetPredictedDepartureTime returns the PredictedDepartureTime field value if set, zero value otherwise.
func (o *JSONPtStop) GetPredictedDepartureTime() time.Time {
	if o == nil || IsNil(o.PredictedDepartureTime) {
		var ret time.Time
		return ret
	}
	return *o.PredictedDepartureTime
}

// GetPredictedDepartureTimeOk returns a tuple with the PredictedDepartureTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetPredictedDepartureTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PredictedDepartureTime) {
		return nil, false
	}
	return o.PredictedDepartureTime, true
}

// HasPredictedDepartureTime returns a boolean if a field has been set.
func (o *JSONPtStop) HasPredictedDepartureTime() bool {
	if o != nil && !IsNil(o.PredictedDepartureTime) {
		return true
	}

	return false
}

// SetPredictedDepartureTime gets a reference to the given time.Time and assigns it to the PredictedDepartureTime field.
func (o *JSONPtStop) SetPredictedDepartureTime(v time.Time) {
	o.PredictedDepartureTime = &v
}

// GetDepartureCancelled returns the DepartureCancelled field value if set, zero value otherwise.
func (o *JSONPtStop) GetDepartureCancelled() bool {
	if o == nil || IsNil(o.DepartureCancelled) {
		var ret bool
		return ret
	}
	return *o.DepartureCancelled
}

// GetDepartureCancelledOk returns a tuple with the DepartureCancelled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *JSONPtStop) GetDepartureCancelledOk() (*bool, bool) {
	if o == nil || IsNil(o.DepartureCancelled) {
		return nil, false
	}
	return o.DepartureCancelled, true
}

// HasDepartureCancelled returns a boolean if a field has been set.
func (o *JSONPtStop) HasDepartureCancelled() bool {
	if o != nil && !IsNil(o.DepartureCancelled) {
		return true
	}

	return false
}

// SetDepartureCancelled gets a reference to the given bool and assigns it to the DepartureCancelled field.
func (o *JSONPtStop) SetDepartureCancelled(v bool) {
	o.DepartureCancelled = &v
}

func (o JSONPtStop) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o JSONPtStop) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.StopId) {
		toSerialize["stop_id"] = o.StopId
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Location) {
		toSerialize["location"] = o.Location
	}
	if !IsNil(o.ArrivalTime) {
		toSerialize["arrival_time"] = o.ArrivalTime
	}
	if !IsNil(o.PlannedArrivalTime) {
		toSerialize["planned_arrival_time"] = o.PlannedArrivalTime
	}
	if !IsNil(o.PredictedArrivalTime) {
		toSerialize["predicted_arrival_time"] = o.PredictedArrivalTime
	}
	if !IsNil(o.ArrivalCancelled) {
		toSerialize["arrival_cancelled"] = o.ArrivalCancelled
	}
	if !IsNil(o.DepartureTime) {
		toSerialize["departure_time"] = o.DepartureTime
	}
	if !IsNil(o.PlannedDepartureTime) {
		toSerialize["planned_departure_time"] = o.PlannedDepartureTime
	}
	if !IsNil(o.PredictedDepartureTime) {
		toSerialize["predicted_departure_time"] = o.PredictedDepartureTime
	}
	if !IsNil(o.DepartureCancelled) {
		toSerialize["departure_cancelled"] = o.DepartureCancelled
	}
	return toSerialize, nil
}

type NullableJSONPtStop struct {
	value *JSONPtStop
	isSet bool
}

func (v NullableJSONPtStop) Get() *JSONPtStop {
	return v.value
}

func (v *NullableJSONPtStop) Set(val *JSONPtStop) {
	v.value = val
	v.isSet = true
}

func (v NullableJSONPtStop) IsSet() bool {
	return v.isSet
}

func (v *NullableJSONPtStop) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableJSONPtStop(val *JSONPtStop) *NullableJSONPtStop {
	return &NullableJSONPtStop{value: val, isSet: true}
}

func (v NullableJSONPtStop) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableJSONPtStop) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


