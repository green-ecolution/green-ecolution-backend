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

// checks if the VehicleUpdate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &VehicleUpdate{}

// VehicleUpdate struct for VehicleUpdate
type VehicleUpdate struct {
	AdditionalInformation map[string]map[string]interface{} `json:"additional_information,omitempty"`
	Description           string                            `json:"description"`
	DrivingLicense        DrivingLicense                    `json:"driving_license"`
	Height                float32                           `json:"height"`
	Length                float32                           `json:"length"`
	Model                 string                            `json:"model"`
	NumberPlate           string                            `json:"number_plate"`
	Provider              *string                           `json:"provider,omitempty"`
	Status                VehicleStatus                     `json:"status"`
	Type                  VehicleType                       `json:"type"`
	WaterCapacity         float32                           `json:"water_capacity"`
	Weight                float32                           `json:"weight"`
	Width                 float32                           `json:"width"`
}

type _VehicleUpdate VehicleUpdate

// NewVehicleUpdate instantiates a new VehicleUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVehicleUpdate(description string, drivingLicense DrivingLicense, height float32, length float32, model string, numberPlate string, status VehicleStatus, type_ VehicleType, waterCapacity float32, weight float32, width float32) *VehicleUpdate {
	this := VehicleUpdate{}
	this.Description = description
	this.DrivingLicense = drivingLicense
	this.Height = height
	this.Length = length
	this.Model = model
	this.NumberPlate = numberPlate
	this.Status = status
	this.Type = type_
	this.WaterCapacity = waterCapacity
	this.Weight = weight
	this.Width = width
	return &this
}

// NewVehicleUpdateWithDefaults instantiates a new VehicleUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVehicleUpdateWithDefaults() *VehicleUpdate {
	this := VehicleUpdate{}
	return &this
}

// GetAdditionalInformation returns the AdditionalInformation field value if set, zero value otherwise.
func (o *VehicleUpdate) GetAdditionalInformation() map[string]map[string]interface{} {
	if o == nil || IsNil(o.AdditionalInformation) {
		var ret map[string]map[string]interface{}
		return ret
	}
	return o.AdditionalInformation
}

// GetAdditionalInformationOk returns a tuple with the AdditionalInformation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetAdditionalInformationOk() (map[string]map[string]interface{}, bool) {
	if o == nil || IsNil(o.AdditionalInformation) {
		return map[string]map[string]interface{}{}, false
	}
	return o.AdditionalInformation, true
}

// HasAdditionalInformation returns a boolean if a field has been set.
func (o *VehicleUpdate) HasAdditionalInformation() bool {
	if o != nil && !IsNil(o.AdditionalInformation) {
		return true
	}

	return false
}

// SetAdditionalInformation gets a reference to the given map[string]map[string]interface{} and assigns it to the AdditionalInformation field.
func (o *VehicleUpdate) SetAdditionalInformation(v map[string]map[string]interface{}) {
	o.AdditionalInformation = v
}

// GetDescription returns the Description field value
func (o *VehicleUpdate) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *VehicleUpdate) SetDescription(v string) {
	o.Description = v
}

// GetDrivingLicense returns the DrivingLicense field value
func (o *VehicleUpdate) GetDrivingLicense() DrivingLicense {
	if o == nil {
		var ret DrivingLicense
		return ret
	}

	return o.DrivingLicense
}

// GetDrivingLicenseOk returns a tuple with the DrivingLicense field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetDrivingLicenseOk() (*DrivingLicense, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DrivingLicense, true
}

// SetDrivingLicense sets field value
func (o *VehicleUpdate) SetDrivingLicense(v DrivingLicense) {
	o.DrivingLicense = v
}

// GetHeight returns the Height field value
func (o *VehicleUpdate) GetHeight() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Height
}

// GetHeightOk returns a tuple with the Height field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetHeightOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Height, true
}

// SetHeight sets field value
func (o *VehicleUpdate) SetHeight(v float32) {
	o.Height = v
}

// GetLength returns the Length field value
func (o *VehicleUpdate) GetLength() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Length
}

// GetLengthOk returns a tuple with the Length field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetLengthOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Length, true
}

// SetLength sets field value
func (o *VehicleUpdate) SetLength(v float32) {
	o.Length = v
}

// GetModel returns the Model field value
func (o *VehicleUpdate) GetModel() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Model
}

// GetModelOk returns a tuple with the Model field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetModelOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Model, true
}

// SetModel sets field value
func (o *VehicleUpdate) SetModel(v string) {
	o.Model = v
}

// GetNumberPlate returns the NumberPlate field value
func (o *VehicleUpdate) GetNumberPlate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NumberPlate
}

// GetNumberPlateOk returns a tuple with the NumberPlate field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetNumberPlateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NumberPlate, true
}

// SetNumberPlate sets field value
func (o *VehicleUpdate) SetNumberPlate(v string) {
	o.NumberPlate = v
}

// GetProvider returns the Provider field value if set, zero value otherwise.
func (o *VehicleUpdate) GetProvider() string {
	if o == nil || IsNil(o.Provider) {
		var ret string
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetProviderOk() (*string, bool) {
	if o == nil || IsNil(o.Provider) {
		return nil, false
	}
	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *VehicleUpdate) HasProvider() bool {
	if o != nil && !IsNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given string and assigns it to the Provider field.
func (o *VehicleUpdate) SetProvider(v string) {
	o.Provider = &v
}

// GetStatus returns the Status field value
func (o *VehicleUpdate) GetStatus() VehicleStatus {
	if o == nil {
		var ret VehicleStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetStatusOk() (*VehicleStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *VehicleUpdate) SetStatus(v VehicleStatus) {
	o.Status = v
}

// GetType returns the Type field value
func (o *VehicleUpdate) GetType() VehicleType {
	if o == nil {
		var ret VehicleType
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetTypeOk() (*VehicleType, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *VehicleUpdate) SetType(v VehicleType) {
	o.Type = v
}

// GetWaterCapacity returns the WaterCapacity field value
func (o *VehicleUpdate) GetWaterCapacity() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.WaterCapacity
}

// GetWaterCapacityOk returns a tuple with the WaterCapacity field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetWaterCapacityOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WaterCapacity, true
}

// SetWaterCapacity sets field value
func (o *VehicleUpdate) SetWaterCapacity(v float32) {
	o.WaterCapacity = v
}

// GetWeight returns the Weight field value
func (o *VehicleUpdate) GetWeight() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Weight
}

// GetWeightOk returns a tuple with the Weight field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetWeightOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Weight, true
}

// SetWeight sets field value
func (o *VehicleUpdate) SetWeight(v float32) {
	o.Weight = v
}

// GetWidth returns the Width field value
func (o *VehicleUpdate) GetWidth() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Width
}

// GetWidthOk returns a tuple with the Width field value
// and a boolean to check if the value has been set.
func (o *VehicleUpdate) GetWidthOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Width, true
}

// SetWidth sets field value
func (o *VehicleUpdate) SetWidth(v float32) {
	o.Width = v
}

func (o VehicleUpdate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o VehicleUpdate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AdditionalInformation) {
		toSerialize["additional_information"] = o.AdditionalInformation
	}
	toSerialize["description"] = o.Description
	toSerialize["driving_license"] = o.DrivingLicense
	toSerialize["height"] = o.Height
	toSerialize["length"] = o.Length
	toSerialize["model"] = o.Model
	toSerialize["number_plate"] = o.NumberPlate
	if !IsNil(o.Provider) {
		toSerialize["provider"] = o.Provider
	}
	toSerialize["status"] = o.Status
	toSerialize["type"] = o.Type
	toSerialize["water_capacity"] = o.WaterCapacity
	toSerialize["weight"] = o.Weight
	toSerialize["width"] = o.Width
	return toSerialize, nil
}

func (o *VehicleUpdate) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"description",
		"driving_license",
		"height",
		"length",
		"model",
		"number_plate",
		"status",
		"type",
		"water_capacity",
		"weight",
		"width",
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

	varVehicleUpdate := _VehicleUpdate{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varVehicleUpdate)

	if err != nil {
		return err
	}

	*o = VehicleUpdate(varVehicleUpdate)

	return err
}

type NullableVehicleUpdate struct {
	value *VehicleUpdate
	isSet bool
}

func (v NullableVehicleUpdate) Get() *VehicleUpdate {
	return v.value
}

func (v *NullableVehicleUpdate) Set(val *VehicleUpdate) {
	v.value = val
	v.isSet = true
}

func (v NullableVehicleUpdate) IsSet() bool {
	return v.isSet
}

func (v *NullableVehicleUpdate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVehicleUpdate(val *VehicleUpdate) *NullableVehicleUpdate {
	return &NullableVehicleUpdate{value: val, isSet: true}
}

func (v NullableVehicleUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVehicleUpdate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
