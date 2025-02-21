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

// checks if the WateringPlanUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WateringPlanUpdate{}

// WateringPlanUpdate struct for WateringPlanUpdate
type WateringPlanUpdate struct {
	AdditionalInformation map[string]interface{} `json:"additional_information,omitempty"`
	CancellationNote      string                 `json:"cancellation_note"`
	Date                  string                 `json:"date"`
	Description           string                 `json:"description"`
	Evaluation            []EvaluationValue      `json:"evaluation,omitempty"`
	Provider              *string                `json:"provider,omitempty"`
	Status                WateringPlanStatus     `json:"status"`
	TrailerId             *int32                 `json:"trailer_id,omitempty"`
	TransporterId         int32                  `json:"transporter_id"`
	TreeClusterIds        []int32                `json:"tree_cluster_ids"`
	UserIds               []string               `json:"user_ids"`
}

type _WateringPlanUpdate WateringPlanUpdate

// NewWateringPlanUpdate instantiates a new WateringPlanUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWateringPlanUpdate(cancellationNote string, date string, description string, status WateringPlanStatus, transporterId int32, treeClusterIds []int32, userIds []string) *WateringPlanUpdate {
	this := WateringPlanUpdate{}
	this.CancellationNote = cancellationNote
	this.Date = date
	this.Description = description
	this.Status = status
	this.TransporterId = transporterId
	this.TreeClusterIds = treeClusterIds
	this.UserIds = userIds
	return &this
}

// NewWateringPlanUpdateWithDefaults instantiates a new WateringPlanUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWateringPlanUpdateWithDefaults() *WateringPlanUpdate {
	this := WateringPlanUpdate{}
	return &this
}

// GetAdditionalInformation returns the AdditionalInformation field value if set, zero value otherwise.
func (o *WateringPlanUpdate) GetAdditionalInformation() map[string]interface{} {
	if o == nil || IsNil(o.AdditionalInformation) {
		var ret map[string]interface{}
		return ret
	}
	return o.AdditionalInformation
}

// GetAdditionalInformationOk returns a tuple with the AdditionalInformation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetAdditionalInformationOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.AdditionalInformation) {
		return map[string]interface{}{}, false
	}
	return o.AdditionalInformation, true
}

// HasAdditionalInformation returns a boolean if a field has been set.
func (o *WateringPlanUpdate) HasAdditionalInformation() bool {
	if o != nil && !IsNil(o.AdditionalInformation) {
		return true
	}

	return false
}

// SetAdditionalInformation gets a reference to the given map[string]interface{} and assigns it to the AdditionalInformation field.
func (o *WateringPlanUpdate) SetAdditionalInformation(v map[string]interface{}) {
	o.AdditionalInformation = v
}

// GetCancellationNote returns the CancellationNote field value
func (o *WateringPlanUpdate) GetCancellationNote() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CancellationNote
}

// GetCancellationNoteOk returns a tuple with the CancellationNote field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetCancellationNoteOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CancellationNote, true
}

// SetCancellationNote sets field value
func (o *WateringPlanUpdate) SetCancellationNote(v string) {
	o.CancellationNote = v
}

// GetDate returns the Date field value
func (o *WateringPlanUpdate) GetDate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Date
}

// GetDateOk returns a tuple with the Date field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetDateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Date, true
}

// SetDate sets field value
func (o *WateringPlanUpdate) SetDate(v string) {
	o.Date = v
}

// GetDescription returns the Description field value
func (o *WateringPlanUpdate) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *WateringPlanUpdate) SetDescription(v string) {
	o.Description = v
}

// GetEvaluation returns the Evaluation field value if set, zero value otherwise.
func (o *WateringPlanUpdate) GetEvaluation() []EvaluationValue {
	if o == nil || IsNil(o.Evaluation) {
		var ret []EvaluationValue
		return ret
	}
	return o.Evaluation
}

// GetEvaluationOk returns a tuple with the Evaluation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetEvaluationOk() ([]EvaluationValue, bool) {
	if o == nil || IsNil(o.Evaluation) {
		return nil, false
	}
	return o.Evaluation, true
}

// HasEvaluation returns a boolean if a field has been set.
func (o *WateringPlanUpdate) HasEvaluation() bool {
	if o != nil && !IsNil(o.Evaluation) {
		return true
	}

	return false
}

// SetEvaluation gets a reference to the given []EvaluationValue and assigns it to the Evaluation field.
func (o *WateringPlanUpdate) SetEvaluation(v []EvaluationValue) {
	o.Evaluation = v
}

// GetProvider returns the Provider field value if set, zero value otherwise.
func (o *WateringPlanUpdate) GetProvider() string {
	if o == nil || IsNil(o.Provider) {
		var ret string
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetProviderOk() (*string, bool) {
	if o == nil || IsNil(o.Provider) {
		return nil, false
	}
	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *WateringPlanUpdate) HasProvider() bool {
	if o != nil && !IsNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given string and assigns it to the Provider field.
func (o *WateringPlanUpdate) SetProvider(v string) {
	o.Provider = &v
}

// GetStatus returns the Status field value
func (o *WateringPlanUpdate) GetStatus() WateringPlanStatus {
	if o == nil {
		var ret WateringPlanStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetStatusOk() (*WateringPlanStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *WateringPlanUpdate) SetStatus(v WateringPlanStatus) {
	o.Status = v
}

// GetTrailerId returns the TrailerId field value if set, zero value otherwise.
func (o *WateringPlanUpdate) GetTrailerId() int32 {
	if o == nil || IsNil(o.TrailerId) {
		var ret int32
		return ret
	}
	return *o.TrailerId
}

// GetTrailerIdOk returns a tuple with the TrailerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetTrailerIdOk() (*int32, bool) {
	if o == nil || IsNil(o.TrailerId) {
		return nil, false
	}
	return o.TrailerId, true
}

// HasTrailerId returns a boolean if a field has been set.
func (o *WateringPlanUpdate) HasTrailerId() bool {
	if o != nil && !IsNil(o.TrailerId) {
		return true
	}

	return false
}

// SetTrailerId gets a reference to the given int32 and assigns it to the TrailerId field.
func (o *WateringPlanUpdate) SetTrailerId(v int32) {
	o.TrailerId = &v
}

// GetTransporterId returns the TransporterId field value
func (o *WateringPlanUpdate) GetTransporterId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.TransporterId
}

// GetTransporterIdOk returns a tuple with the TransporterId field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetTransporterIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TransporterId, true
}

// SetTransporterId sets field value
func (o *WateringPlanUpdate) SetTransporterId(v int32) {
	o.TransporterId = v
}

// GetTreeClusterIds returns the TreeClusterIds field value
func (o *WateringPlanUpdate) GetTreeClusterIds() []int32 {
	if o == nil {
		var ret []int32
		return ret
	}

	return o.TreeClusterIds
}

// GetTreeClusterIdsOk returns a tuple with the TreeClusterIds field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetTreeClusterIdsOk() ([]int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.TreeClusterIds, true
}

// SetTreeClusterIds sets field value
func (o *WateringPlanUpdate) SetTreeClusterIds(v []int32) {
	o.TreeClusterIds = v
}

// GetUserIds returns the UserIds field value
func (o *WateringPlanUpdate) GetUserIds() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.UserIds
}

// GetUserIdsOk returns a tuple with the UserIds field value
// and a boolean to check if the value has been set.
func (o *WateringPlanUpdate) GetUserIdsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.UserIds, true
}

// SetUserIds sets field value
func (o *WateringPlanUpdate) SetUserIds(v []string) {
	o.UserIds = v
}

func (o WateringPlanUpdate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WateringPlanUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AdditionalInformation) {
		toSerialize["additional_information"] = o.AdditionalInformation
	}
	toSerialize["cancellation_note"] = o.CancellationNote
	toSerialize["date"] = o.Date
	toSerialize["description"] = o.Description
	if !IsNil(o.Evaluation) {
		toSerialize["evaluation"] = o.Evaluation
	}
	if !IsNil(o.Provider) {
		toSerialize["provider"] = o.Provider
	}
	toSerialize["status"] = o.Status
	if !IsNil(o.TrailerId) {
		toSerialize["trailer_id"] = o.TrailerId
	}
	toSerialize["transporter_id"] = o.TransporterId
	toSerialize["tree_cluster_ids"] = o.TreeClusterIds
	toSerialize["user_ids"] = o.UserIds
	return toSerialize, nil
}

func (o *WateringPlanUpdate) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"cancellation_note",
		"date",
		"description",
		"status",
		"transporter_id",
		"tree_cluster_ids",
		"user_ids",
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

	varWateringPlanUpdate := _WateringPlanUpdate{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varWateringPlanUpdate)

	if err != nil {
		return err
	}

	*o = WateringPlanUpdate(varWateringPlanUpdate)

	return err
}

type NullableWateringPlanUpdate struct {
	value *WateringPlanUpdate
	isSet bool
}

func (v NullableWateringPlanUpdate) Get() *WateringPlanUpdate {
	return v.value
}

func (v *NullableWateringPlanUpdate) Set(val *WateringPlanUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableWateringPlanUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableWateringPlanUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWateringPlanUpdate(val *WateringPlanUpdate) *NullableWateringPlanUpdate {
	return &NullableWateringPlanUpdate{value: val, isSet: true}
}

func (v NullableWateringPlanUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWateringPlanUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
