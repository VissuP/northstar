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

package model

// JSON structure to respond to the /oauth2/token request coming from the portal
type UserToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

// JSON structure to respond to the /oauth2/token/info request coming from the NS API
type UserTokenInfo struct {
	GrantType       string   `json:"grant_type,omitempty"`
	ExpiresIn       int      `json:"expires_in"`
	CreatedAt       string   `json:"created_at,omitempty"`
	UserName        string   `json:"username,omitempty"`
	Scope           []string `json:"scope,omitempty"`
	ClientId        string   `json:"clientid,omitempty"`
	CredentialsType string   `json:"credentialstype,omitempty"`
}

// JSON structure to respond to the /api/v2/users/me request coming from the portal
type GetUserInfo struct {
	Id              string `json:"id,omitempty"`
	Kind            string `json:"kind,omitempty"`
	Version         string `json:"version,omitempty"`
	VersionId       string `json:"versionid,omitempty"`
	CreatedOn       string `json:"createdon,omitempty"`
	LastUpdated     string `json:"lastupdated,omitempty"`
	Number          string `json:"number,omitempty"`
	RequestedOn     string `json:"requestedon,omitempty"`
	Expiration      int    `json:"expiration"`
	CredentialsId   string `json:"credentialsid,omitempty"`
	CredentialsType string `json:"credentialstype,omitempty"`
	State           string `json:"state,omitempty"`
	DisplayName     string `json:"displayname,omitempty"`
	AckTermsOn      string `json:"acktermson,omitempty"`
	FirstName       string `json:"firstname,omitempty"`
	LastName        string `json:"lastname,omitempty"`
	Email           string `json:"email,omitempty"`
}

// JSON structure to respond to the /south/v2/users request coming from the NS API
type SouthUserInfo struct {
	Id              string `json:"id,omitempty"`
	Kind            string `json:"kind,omitempty"`
	Version         string `json:"version,omitempty"`
	VersionId       string `json:"versionid,omitempty"`
	CreatedOn       string `json:"createdon,omitempty"`
	LastUpdated     string `json:"lastupdated,omitempty"`
	ForeignId       string `json:"foreignid,omitempty"`
	Number          string `json:"number,omitempty"`
	RequestedOn     string `json:"requestedon,omitempty"`
	Expiration      int    `json:"expiration"`
	CredentialsId   string `json:"credentialsid,omitempty"`
	CredentialsType string `json:"credentialstype,omitempty"`
	State           string `json:"state,omitempty"`
	DisplayName     string `json:"displayname,omitempty"`
	AckTermsOn      string `json:"acktermson,omitempty"`
	FirstName       string `json:"firstname,omitempty"`
	LastName        string `json:"lastname,omitempty"`
	Email           string `json:"email,omitempty"`
}
