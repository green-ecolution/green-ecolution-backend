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
	"bytes"
	"fmt"
)

// checks if the ClientToken type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClientToken{}

// ClientToken struct for ClientToken
type ClientToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int32 `json:"expires_in"`
	IdToken string `json:"id_token"`
	NotBeforePolicy int32 `json:"not_before_policy"`
	RefreshExpiresIn int32 `json:"refresh_expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	SessionState string `json:"session_state"`
	TokenType string `json:"token_type"`
}

type _ClientToken ClientToken

// NewClientToken instantiates a new ClientToken object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClientToken(accessToken string, expiresIn int32, idToken string, notBeforePolicy int32, refreshExpiresIn int32, refreshToken string, scope string, sessionState string, tokenType string) *ClientToken {
	this := ClientToken{}
	this.AccessToken = accessToken
	this.ExpiresIn = expiresIn
	this.IdToken = idToken
	this.NotBeforePolicy = notBeforePolicy
	this.RefreshExpiresIn = refreshExpiresIn
	this.RefreshToken = refreshToken
	this.Scope = scope
	this.SessionState = sessionState
	this.TokenType = tokenType
	return &this
}

// NewClientTokenWithDefaults instantiates a new ClientToken object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClientTokenWithDefaults() *ClientToken {
	this := ClientToken{}
	return &this
}

// GetAccessToken returns the AccessToken field value
func (o *ClientToken) GetAccessToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AccessToken
}

// GetAccessTokenOk returns a tuple with the AccessToken field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetAccessTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AccessToken, true
}

// SetAccessToken sets field value
func (o *ClientToken) SetAccessToken(v string) {
	o.AccessToken = v
}

// GetExpiresIn returns the ExpiresIn field value
func (o *ClientToken) GetExpiresIn() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.ExpiresIn
}

// GetExpiresInOk returns a tuple with the ExpiresIn field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetExpiresInOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpiresIn, true
}

// SetExpiresIn sets field value
func (o *ClientToken) SetExpiresIn(v int32) {
	o.ExpiresIn = v
}

// GetIdToken returns the IdToken field value
func (o *ClientToken) GetIdToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IdToken
}

// GetIdTokenOk returns a tuple with the IdToken field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetIdTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IdToken, true
}

// SetIdToken sets field value
func (o *ClientToken) SetIdToken(v string) {
	o.IdToken = v
}

// GetNotBeforePolicy returns the NotBeforePolicy field value
func (o *ClientToken) GetNotBeforePolicy() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NotBeforePolicy
}

// GetNotBeforePolicyOk returns a tuple with the NotBeforePolicy field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetNotBeforePolicyOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NotBeforePolicy, true
}

// SetNotBeforePolicy sets field value
func (o *ClientToken) SetNotBeforePolicy(v int32) {
	o.NotBeforePolicy = v
}

// GetRefreshExpiresIn returns the RefreshExpiresIn field value
func (o *ClientToken) GetRefreshExpiresIn() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.RefreshExpiresIn
}

// GetRefreshExpiresInOk returns a tuple with the RefreshExpiresIn field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetRefreshExpiresInOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RefreshExpiresIn, true
}

// SetRefreshExpiresIn sets field value
func (o *ClientToken) SetRefreshExpiresIn(v int32) {
	o.RefreshExpiresIn = v
}

// GetRefreshToken returns the RefreshToken field value
func (o *ClientToken) GetRefreshToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RefreshToken
}

// GetRefreshTokenOk returns a tuple with the RefreshToken field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetRefreshTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RefreshToken, true
}

// SetRefreshToken sets field value
func (o *ClientToken) SetRefreshToken(v string) {
	o.RefreshToken = v
}

// GetScope returns the Scope field value
func (o *ClientToken) GetScope() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Scope
}

// GetScopeOk returns a tuple with the Scope field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetScopeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Scope, true
}

// SetScope sets field value
func (o *ClientToken) SetScope(v string) {
	o.Scope = v
}

// GetSessionState returns the SessionState field value
func (o *ClientToken) GetSessionState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SessionState
}

// GetSessionStateOk returns a tuple with the SessionState field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetSessionStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SessionState, true
}

// SetSessionState sets field value
func (o *ClientToken) SetSessionState(v string) {
	o.SessionState = v
}

// GetTokenType returns the TokenType field value
func (o *ClientToken) GetTokenType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TokenType
}

// GetTokenTypeOk returns a tuple with the TokenType field value
// and a boolean to check if the value has been set.
func (o *ClientToken) GetTokenTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TokenType, true
}

// SetTokenType sets field value
func (o *ClientToken) SetTokenType(v string) {
	o.TokenType = v
}

func (o ClientToken) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClientToken) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["access_token"] = o.AccessToken
	toSerialize["expires_in"] = o.ExpiresIn
	toSerialize["id_token"] = o.IdToken
	toSerialize["not_before_policy"] = o.NotBeforePolicy
	toSerialize["refresh_expires_in"] = o.RefreshExpiresIn
	toSerialize["refresh_token"] = o.RefreshToken
	toSerialize["scope"] = o.Scope
	toSerialize["session_state"] = o.SessionState
	toSerialize["token_type"] = o.TokenType
	return toSerialize, nil
}

func (o *ClientToken) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"access_token",
		"expires_in",
		"id_token",
		"not_before_policy",
		"refresh_expires_in",
		"refresh_token",
		"scope",
		"session_state",
		"token_type",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varClientToken := _ClientToken{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varClientToken)

	if err != nil {
		return err
	}

	*o = ClientToken(varClientToken)

	return err
}

type NullableClientToken struct {
	value *ClientToken
	isSet bool
}

func (v NullableClientToken) Get() *ClientToken {
	return v.value
}

func (v *NullableClientToken) Set(val *ClientToken) {
	v.value = val
	v.isSet = true
}

func (v NullableClientToken) IsSet() bool {
	return v.isSet
}

func (v *NullableClientToken) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClientToken(val *ClientToken) *NullableClientToken {
	return &NullableClientToken{value: val, isSet: true}
}

func (v NullableClientToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClientToken) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


