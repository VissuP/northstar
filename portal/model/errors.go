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

import (
	"net/http"

	"fmt"
	"github.com/verizonlabs/northstar/pkg/management"
)

// Define user specific errors returned by the API.
var (
	// Define service specific errors returned by the API.
	ErrorParseRequestBody = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The request body is invalid. Failed to parse content."}
	ErrorMissingEmail     = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The email is missing."}
	ErrorMissingPassword  = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The password is missing."}
	ErrorMissingCookie    = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The request does not contain auth information."}
	ErrorInvalidCookie    = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The request contains invalid auth information."}
	ErrorUnauthorized     = &management.Error{HttpStatus: http.StatusUnauthorized, Id: "unauthorized_client", Description: "The request requires user authentication."}
)

var (
	ErrorEventIdMissing   = NewErrorEvent("The event ID is missing.")
	ErrorEventTypeMissing = NewErrorEvent("The event type is missing.")
)

func GetErrorMissingResourceID(resource string) *management.Error {
	return &management.Error{
		HttpStatus:  http.StatusBadRequest,
		Id:          management.ERR_BAD_REQUEST,
		Description: fmt.Sprintf("The %s is missing or invalid.", resource),
	}
}
