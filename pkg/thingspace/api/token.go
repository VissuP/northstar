/*
Copyright (C) 2017 Verizon. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"time"
)

const (
	// Defines the access token grant types.
	AuthorizationCodeGrantType string = "authorization_code"
	PasswordGrantType          string = "password"
	ClientCredentialsGrantType string = "client_credentials"
)

// Defines the auth token type.
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Type         string    `json:"token_type"`
	ExpiresIn    float64   `json:"expires_in"`
	Scope        string    `json:"scope"`
	CreatedAt    time.Time `json:"-"`
}

// Defines the type used to represent token information.
type TokenInfo struct {
	GrantType  string    `json:"grant_type,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ExpiresIn  int32     `json:"expires_in,omitempty"`
	UserName   string    `json:"username"`
	CustomData string    `json:"custom_data,omitempty"`
	Scopes     []string  `json:"scope,omitempty"`
	ClientId   string    `json:"clientid"`
}

// Returns true if token is expired, false otherwise.
func (token Token) IsExpired() bool {
	return token.CreatedAt.Add(time.Duration(token.ExpiresIn) * time.Second).Before(time.Now())
}

// Returns true if token info is expired, false otherwise.
func (tokenInfo TokenInfo) IsExpired() bool {
	return tokenInfo.CreatedAt.Add(time.Duration(tokenInfo.ExpiresIn) * time.Second).Before(time.Now())
}
