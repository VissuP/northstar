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

package thingspace

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pmylund/go-cache"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	. "github.com/verizonlabs/northstar/pkg/thingspace/api"
	"strconv"
	"time"
)

const (
	CacheExpiration = 1 * time.Hour
	CachePurgeTime  = 30 * time.Minute
)

// Defines the auth client interface.
type AuthClient interface {
	GetClientToken(clientId string, clientSecret string, scope string) (*Token, *management.Error)
	GetUserToken(clientId, clientSecret, username, password, scope string) (*Token, *management.Error)
	RevokeAccessToken(clientId, clientSecret, token string) *management.Error
	GetTokenInfo(token string) (*TokenInfo, *management.Error)
}

// Defines the ThingSpace Auth client.
type ThingSpaceAuthClient struct {
	BaseClient
	tokenCache *cache.Cache
}

// Returns a new ThingSpace Auth Service client.
func NewThingSpaceAuthClient(thingSpaceHostAndPort string) AuthClient {
	return NewThingSpaceAuthClientWithProtocol("http", thingSpaceHostAndPort)
}

// Returns a new ThingSpace Auth Service client.
func NewThingSpaceAuthClientWithProtocol(protocol, hostAndPort string) AuthClient {
	authClient := &ThingSpaceAuthClient{
		BaseClient: BaseClient{
			httpClient: management.NewHttpClient(),
			baseUrl:    fmt.Sprintf("%s://%s", protocol, hostAndPort),
		},
		tokenCache: cache.New(CacheExpiration, CachePurgeTime),
	}

	return authClient
}

// Returns the token associated with client credentials.
func (client *ThingSpaceAuthClient) GetClientToken(clientId string, clientSecret string, scope string) (*Token, *management.Error) {
	mlog.Debug("GetClientToken")

	// Try to get the token from the cache.
	credentials := clientId + ":" + clientSecret
	cacheItem, found := client.tokenCache.Get(credentials)

	// If the item is found, we check if it has expire. If it did,
	// the request will be made again and item will be store in cache.
	// Otherwise, we return the value of the cache item.
	if found == true {
		mlog.Debug("Found internal token in cache.")
		token := cacheItem.(Token)

		if token.IsExpired() == false {
			mlog.Debug("Cache item still valid, reusing value.")
			return &token, nil
		} else {
			mlog.Debug("Cache item expired ")
		}
	}

	// Create request body.
	reqBody := url.Values{}
	reqBody.Set("grant_type", "client_credentials")
	reqBody.Set("scope", scope)

	// Create the request.
	request, _ := http.NewRequest("POST", client.baseUrl+"/oauth2/token", bytes.NewBufferString(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))
	request.Header.Set("Connection", "close")

	// Send the request.
	response, err := client.httpClient.Do(request)

	if err != nil {
		mlog.Error("Error, failed to get client token from url %s with error: %s.", client.baseUrl, err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	// Validate the response body and status code
	if response == nil || response.Body == nil {
		return nil, management.GetExternalError("Error, the response, or response body, is not valid.")
	}

	// Get the token from the response body
	respBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		mlog.Error("Error, failed to get client token with error: %s.", string(respBody))
		return nil, management.ErrorExternal
	}

	var token Token

	if err = json.Unmarshal(respBody, &token); err != nil {
		mlog.Error("Error, failed to unmarshal token from response with error: %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	// Note that response does not contains created at. We add to facilitate
	// validation when used from cache.
	token.CreatedAt = time.Now()

	// Store in cache using token expiration time.
	client.tokenCache.Set(credentials, token, time.Duration(token.ExpiresIn)*time.Second)

	return &token, nil
}

// Returns the token associated with user credentials.
func (client *ThingSpaceAuthClient) GetUserToken(clientId, clientSecret, username, password, scope string) (*Token, *management.Error) {
	mlog.Debug("GetUserToken")

	// Note that we do not cache user token.

	credentials := clientId + ":" + clientSecret
	reqBody := url.Values{}
	reqBody.Set("grant_type", "password")
	reqBody.Set("scope", scope)
	reqBody.Set("username", username)
	reqBody.Set("password", password)

	request, _ := http.NewRequest("POST", client.baseUrl+"/oauth2/token", bytes.NewBufferString(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))
	request.Header.Set("Connection", "close")

	// Send the request.
	response, err := client.httpClient.Do(request)

	if err != nil {
		mlog.Error("Error, failed to get user token from url %s with error: %s.", client.baseUrl, err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	// Validate the response body and status code
	if response == nil || response.Body == nil {
		return nil, management.GetExternalError("Error, the response, or response body, is not valid.")
	}

	// Get the token from the response body
	respBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		mlog.Error("Error, failed to get user token with error: %s.", string(respBody))
		return nil, management.ErrorExternal
	}

	var token Token

	if err = json.Unmarshal(respBody, &token); err != nil {
		mlog.Error("Error, failed to unmarshal token from response with error: %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	return &token, nil
}

// Revokes, or removes, the specified access token.
func (client *ThingSpaceAuthClient) RevokeAccessToken(clientId, clientSecret, token string) *management.Error {
	mlog.Debug("RevokeAccessToken")

	reqBody := url.Values{}
	reqBody.Set("token", token)
	reqBody.Set("token_type_hint", "access_token")

	credentials := clientId + ":" + clientSecret
	request, _ := http.NewRequest("POST", client.baseUrl+"/oauth2/revoke", bytes.NewBufferString(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))
	request.Header.Set("Connection", "close")

	// Send the request.
	response, err := client.httpClient.Do(request)

	if err != nil {
		mlog.Error("Error, failed to get user token from url %s with error: %s.", client.baseUrl, err.Error())
		return management.GetInternalError(err.Error())
	}

	// Validate the response body and status code
	if response == nil || response.Body == nil {
		return management.GetExternalError("Error, the response, or response body, is not valid.")
	}

	// Validate http status.
	if response.StatusCode != http.StatusOK {
		// Get the error from the response body
		respBody, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		mlog.Error("Error, failed to revoke token with error: %s.", string(respBody))
		return management.ErrorExternal
	}

	return nil
}

// Returns information associated with the specified token.
func (client *ThingSpaceAuthClient) GetTokenInfo(token string) (*TokenInfo, *management.Error) {
	mlog.Debug("GetTokenInfo")

	url := client.baseUrl + "/oauth2/token/info?access_token=" + token

	// Request token information.
	body, err := client.executeHttpMethod("GET", url, "", nil)

	if err != nil {
		mlog.Error("Error, failed to execute http request for url %s with error: %s", url, err.Description)
		return nil, err
	}

	tokenInfo := &TokenInfo{}

	if goErr := json.Unmarshal(body, tokenInfo); goErr != nil {
		mlog.Error("Error, failed to unmarshal token information with error: %s", goErr.Error())
		return nil, management.GetInternalError(goErr.Error())
	}

	return tokenInfo, nil
}
