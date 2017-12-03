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

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"time"
)

const (
	//HTTPCookieName defines the constant used as the HTTP Cookie Name.
	HTTPCookieName = "Ns.Auth.HttpCookie"
)

// SetCookie is a helper method used to create HTTP Cookie for the specified Auth Token.
func SetCookie(context *gin.Context, token *api.Token) error {
	mlog.Debug("SetCookie")

	// TODO - We need to protect the token, e.g., encode with something
	// like https://github.com/gorilla/securecookie or encryp, use JWT, etc.
	// For now, we just create JSON and base64 encoding.
	value, err := json.Marshal(&token)

	if err != nil {
		return fmt.Errorf("Failed to marshal auth token with error: %v", err)
	}

	encodedToken := base64.StdEncoding.EncodeToString(value)

	// TODO - We need to look at the values for: Path, Domain, and Secure.

	// Create http cookie.
	cookie := &http.Cookie{
		Name:     HTTPCookieName,
		Value:    encodedToken,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}

	mlog.Debug("Set HTTP cookie: %+v", cookie)
	http.SetCookie(context.Writer, cookie)

	return nil
}

// GetToken is a helper method used to get Auth Token from HTTP Cookie.
func GetToken(context *gin.Context) (*api.Token, error) {
	mlog.Debug("GetToken")

	// Get the cookie from the request
	cookie, err := context.Request.Cookie(HTTPCookieName)

	if err != nil {
		return nil, fmt.Errorf("Failed to get HTTP cookie %s with error: %v", HTTPCookieName, err)
	}

	// Decode the cookie value.
	decodedToken, err := base64.StdEncoding.DecodeString(cookie.Value)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode HTTP cookie %s with error: %v", HTTPCookieName, err)
	}

	// Umarshal the auth token.
	var token api.Token

	if err := json.Unmarshal(decodedToken, &token); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal HTTP cookie %s with error: %v", HTTPCookieName, err)
	}

	return &token, nil
}

// DeleteCookie is a helper method used to remove HTTP Cookie.
func DeleteCookie(context *gin.Context) {
	mlog.Debug("DeleteCookie")

	// Invalid the http cookie.
	cookie := &http.Cookie{
		Name:    HTTPCookieName,
		Value:   "",                           // Set empty value.
		Path:    "/",                          // Route back to login page.
		MaxAge:  -1,                           // Note: -1 means delete right now.
		Expires: time.Now().AddDate(0, 0, -1), // Some browsers don't support max age. Include expires just in case.
	}

	http.SetCookie(context.Writer, cookie)
}
