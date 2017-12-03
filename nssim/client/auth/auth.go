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

package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"net/http/httputil"

	"github.com/pmylund/go-cache"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/config"
)

var (
	DEFAULT_CACHE_EXPIRATION = 1 * time.Hour
	DEFAULT_CACHE_PURGE_TIME = 30 * time.Minute
	tokenCache               = cache.New(DEFAULT_CACHE_EXPIRATION, DEFAULT_CACHE_PURGE_TIME)
	authClient               = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5000,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			//HTTP methods do not have a timeout by default. Add one for initiating the connection.
			Dial: func(netw, addr string) (net.Conn, error) {
				mlog.Event("Connecting to %s", addr)
				c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(config.Configuration.HttpTimeout))
				if err != nil {
					if strings.Compare(reflect.TypeOf(err).String(), "*net.OpError") == 0 {
						mlog.Event("Unable to connect to %s", addr)
					}
					mlog.Event("Disconnected from %s", addr)
					return nil, err
				}
				mlog.Event("Connected to %s", addr)
				return c, nil
			},
			//HTTP methods do not have a timeout by default. Add one for waiting for a response.
			ResponseHeaderTimeout: time.Second * time.Duration(config.Configuration.HttpTimeout),
		},
	}
)

// Defines the auth token type.
type Token struct {
	AccessToken  string
	RefreshToken string
	Type         string
	ExpiresIn    float64
	Scope        string
	CreatedAt    time.Time
}

// Returns true if token is expired, false otherwise.
func (token *Token) IsExpired() bool {
	return token.CreatedAt.Add(time.Duration(token.ExpiresIn) * time.Second).Before(time.Now())
}

// Returns the token associated with client credentials.
func GetClientToken() (*Token, *management.Error) {
	mlog.Debug("GetClientToken")

	credentials := config.Configuration.Credentials.Client.Id + ":" + config.Configuration.Credentials.Client.Secret
	/// Try to get the token from the cache.
	cacheItem, found := tokenCache.Get(credentials)

	// If the item is found, we check if it has expire. If it did,
	// the request will be made again and item will be store in cache.
	// Otherwise, we return the value of the cache item.
	if found {
		mlog.Debug("Found internal token in cache.")
		token := cacheItem.(*Token)

		if token.IsExpired() == false {
			mlog.Debug("Cache item still valid, reusing value.")
			return token, nil
		} else {
			mlog.Debug("Cache item expired ")
		}
	}

	reqBody := url.Values{}
	reqBody.Set("grant_type", "client_credentials")
	reqBody.Set("scope", config.Configuration.Credentials.Client.Scope)

	baseUrl := fmt.Sprintf("%s://%s", config.Configuration.ThingspaceProtocol, config.Configuration.ThingspaceAuthHostPort)
	request, _ := http.NewRequest("POST", baseUrl+"/oauth2/token", bytes.NewBufferString(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))
	request.Header.Set("Connection", "close")
	mlog.Debug("client-token request headers %+v", request.Header)

	// Note that we intentionally create a new client here.
	response, err := authClient.Do(request)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Error, failed to get client token with error: %s.", err.Error()))
	} else if response == nil || response.Body == nil {
		return nil, management.GetExternalError(fmt.Sprintf("Error, the response, or response body, is not valid."))
	} else if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		respBody, _ := ioutil.ReadAll(response.Body)
		return nil, management.GetExternalError(fmt.Sprintf("Error, failed to get client token with error: %s.", string(respBody)))
	}

	respBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	mlog.Debug("Response body: %s", string(respBody))

	var tokenInfo map[string]interface{}
	json.Unmarshal(respBody, &tokenInfo)

	token, mErr := getToken(tokenInfo)
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Error, failed to generate token from response with error: %s", mErr.Description))
	}

	// Store in cache using token expiration time.
	tokenCache.Set(credentials, token, time.Duration(token.ExpiresIn)*time.Second)

	return token, nil
}

// Returns the token associated with user credentials.
func GetUserToken(username, password string) (*Token, *management.Error) {
	mlog.Debug("GetUserTokenFromCredentials")

	/// Try to get the token from the cache.
	cacheItem, found := tokenCache.Get(username)

	// If the item is found, we check if it has expire. If it did,
	// the request will be made again and item will be store in cache.
	// Otherwise, we return the value of the cache item.
	if found {
		mlog.Debug("Found internal token in cache.")
		token := cacheItem.(*Token)

		if token.IsExpired() == false {
			mlog.Debug("Cache item still valid, reusing value.")
			return token, nil
		} else {
			mlog.Debug(" Cache item expired ")
		}
	}

	credentials := config.Configuration.Credentials.Client.Id + ":" + config.Configuration.Credentials.Client.Secret
	reqBody := url.Values{}
	reqBody.Set("grant_type", "password")
	reqBody.Set("scope", config.Configuration.Credentials.User.Scope)
	reqBody.Set("username", username)
	reqBody.Set("password", password)

	baseUrl := fmt.Sprintf("%s://%s", config.Configuration.ThingspaceProtocol, config.Configuration.ThingspaceAuthHostPort)
	request, _ := http.NewRequest("POST", baseUrl+"/oauth2/token", bytes.NewBufferString(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))
	// TODO : do we really need this authorization header below?
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))
	request.Header.Set("Connection", "close")

	reqDumpBody, _ := httputil.DumpRequest(request, true)
	mlog.Info("Get  user token request: %s", reqDumpBody)

	// Note that we intentionally create a new client here.
	response, err := authClient.Do(request)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Error, failed to get user token with error: %s.", err.Error()))
	} else if response == nil || response.Body == nil {
		return nil, management.GetExternalError(fmt.Sprintf("Error, the response, or response body, is not valid."))
	} else if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, management.GetExternalError(fmt.Sprintf("Error, failed to get user token for user %s with error: %v.", username, err))
		}

		if respBody != nil {
			return nil, management.GetExternalError(fmt.Sprintf("Error, failed to get user token for user %s with response body string: %s.", username, string(respBody)))
		}

		return nil, management.GetExternalError(fmt.Sprintf(" Error response body for user %s was nil", username))
	}

	respBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	mlog.Debug("Response body: %s", string(respBody))

	var tokenInfo map[string]interface{}
	json.Unmarshal(respBody, &tokenInfo)

	token, mErr := getToken(tokenInfo)
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Error, failed to generate token from response with error: %s", mErr.Description))
	}

	// Store in cache using token expiration time.
	tokenCache.Set(username, token, time.Duration(token.ExpiresIn)*time.Second)

	return token, nil
}

// Revokes the given token (token must belong to expected client associated with user credentials).
func RevokeUserToken(token string) *management.Error {
	reqBody := url.Values{}
	reqBody.Set("token", token)
	reqBody.Set("token_type_hint", "access_token")

	credentials := config.Configuration.Credentials.Client.Id + ":" + config.Configuration.Credentials.Client.Secret
	baseUrl := fmt.Sprintf("%s://%s", config.Configuration.ThingspaceProtocol, config.Configuration.ThingspaceAuthHostPort)
	request, _ := http.NewRequest("POST", baseUrl+"/oauth2/revoke", bytes.NewBufferString(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(reqBody.Encode())))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(credentials))))
	request.Header.Set("Connection", "close")

	response, err := authClient.Do(request)

	if err != nil {
		return management.GetExternalError(fmt.Sprintf("Error, failed to revoke token with error: %s.", err.Error()))
	} else if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		respBody, _ := ioutil.ReadAll(response.Body)
		return management.GetExternalError(fmt.Sprintf("Error, failed to revoke user token with error: %s.", string(respBody)))
	}
	ioutil.ReadAll(response.Body)
	response.Body.Close()
	return nil
}

// Revokes the given client token
func RevokeClientToken(token string) *management.Error {
	// Revoke the token
	mErr := RevokeUserToken(token)
	if mErr != nil {
		return mErr
	}

	// Remove from cache
	credentials := config.Configuration.Credentials.Client.Id + ":" + config.Configuration.Credentials.Client.Secret
	tokenCache.Delete(credentials)
	return nil
}

// Helper method used to get token.
func getToken(tokenInfo map[string]interface{}) (*Token, *management.Error) {
	// Get Access Token
	value, valid := tokenInfo["access_token"]

	if valid == false {
		return nil, management.GetExternalError(fmt.Sprintf("Error response contains an invalid access token."))
	}

	accessToken := string(value.(string))
	refreshToken := ""
	tokenType := ""
	expiresIn := float64(0)

	// Get Refresh Token
	if value, valid := tokenInfo["refresh_token"]; valid {
		mlog.Debug("Got valid refresh token.")
		refreshToken = string(value.(string))
	}

	// Get Token Type
	if value, valid := tokenInfo["token_type"]; valid {
		mlog.Debug("Got valid token type.")
		tokenType = string(value.(string))
	}

	// Get expires in.
	if value, valid := tokenInfo["expires_in"]; valid {
		mlog.Debug("Got valid token type.")
		expiresIn = float64(value.(float64))
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         tokenType,
		ExpiresIn:    expiresIn,
	}

	return token, nil
}
