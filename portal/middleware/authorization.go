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

package middleware

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/portal/model"
	"github.com/verizonlabs/northstar/portal/utils"
)

const (
	// AccessTokenKeyName defines the name of the Access Token key.
	AccessTokenKeyName = "Ns.Auth.AccessToken"
)

// To authenticate another type of cookie, add the name of the cookie here
// in addition to the name of the function that will handle extracting and
// validating the cookie. This function should take a *http.Cookie object
// and a *gin.Context object, and return a boolean where true indicates that
// the token was successfully extracted and the cookie was validated
var cookieMap = map[string]interface{}{
	"TNDP.Auth.HttpCookie" : HandleTNDPCookie,
	"Ns.Auth.HttpCookie" : HandleNSCookie,
}

func HandleTNDPCookie(cookie *http.Cookie, context *gin.Context) bool {
	jwtObj, err := utils.GetJWTFromString(cookie.Value)
	if err != nil || !utils.ValidateTSCoreJWT(jwtObj) {
		return false
	} else {
		var access_token = jwtObj.Claims.(jwt.MapClaims)["token"].(string)

		decrypted, err := utils.B64DecodeAndDecryptAccessToken(access_token)
		if err != nil {
			serviceError := management.GetInternalError("An error occurred trying to decrypt the access token")
			// per docs, headers need to be set before calling context.JSON method
			for k, v := range serviceError.Header {
				for _, v1 := range v {
					context.Writer.Header().Add(k, v1)
				}
			}
			// now serialize rest of the response
			context.JSON(serviceError.HttpStatus, serviceError)
		}

		context.Set(AccessTokenKeyName, &api.Token{AccessToken : decrypted})
		return true
	}
}

func HandleNSCookie(cookie *http.Cookie, context *gin.Context) bool {
	// Get the token from the request http cookie.
	token, err := utils.GetToken(context)
	context.Set(AccessTokenKeyName, token)
	return err == nil
}

// Search through every potential cookie and authenticate the request if any of them are valid
func Authorization(context *gin.Context) {
	authorized := false
	for key, val := range cookieMap {
		cookie, err := context.Request.Cookie(key)
		if err == nil {
			if val.(func(*http.Cookie, *gin.Context)bool)(cookie, context) {
				mlog.Info("Request successfully validated using %s cookie", key)
				context.Next()
				authorized = true
				break
			}
		}
	}

	if !authorized {
		context.Writer.Header().Set("WWW-Authenticate", fmt.Sprintf("Cookie realm=%s", utils.HTTPCookieName))
		context.JSON(model.ErrorUnauthorized.HttpStatus, model.ErrorUnauthorized)
		context.Abort()
	}
}
