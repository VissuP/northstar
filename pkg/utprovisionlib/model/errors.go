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
	"github.com/verizonlabs/northstar/pkg/management"
)

//Error IDs
const (
	ERR_CONFLICT            = "conflict_error"
	ERR_PRECONDITION_FAILED = "precondition_error"
	ERR_REGISTERED          = "device_registered"
	ERR_INVALID_DEVICE      = "invalid_device"
	ERR_INVALID_IMSI        = "invalid_imsi"
	ERR_INVALID_IMEI        = "invalid_imei"
	ERR_INVALID_ENCODING    = "invalid_encoding"
	ERR_INVALID_CERT        = "invalid_certificate"
	ERR_INVALID_SUBJECT     = "invalid_subject"
)

//Defined management Errors
var (
	ErrorMissingRequestBody = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The request body is missing."}
	ErrorInvalidResourceId  = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The resource id is missing or invalid."}
	ErrorInvalidDevice      = &management.Error{HttpStatus: http.StatusBadRequest, Id: ERR_INVALID_DEVICE, Description: "The specified device contains missing or invalid fields."}
	ErrorMissingPathParam   = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The path parameter is missing."}
	ErrorMissingImsi        = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The IMSI is missing."}
	ErrorInvalidImsi        = &management.Error{HttpStatus: http.StatusBadRequest, Id: ERR_INVALID_IMSI, Description: "The IMSI is invalid."}
	ErrorInvalidImei        = &management.Error{HttpStatus: http.StatusBadRequest, Id: ERR_INVALID_IMEI, Description: "The IMEI is invalid."}
	ErrorInvalidDeviceModel = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The Model Id was not found or the service was unable to validate."}
	ErrorInvalidEncoding    = &management.Error{HttpStatus: http.StatusBadRequest, Id: ERR_INVALID_ENCODING, Description: "The specified certificate encoding is invalid."}
	ErrorInvalidCert        = &management.Error{HttpStatus: http.StatusBadRequest, Id: ERR_INVALID_CERT, Description: "The specified certificate is invalid."}
	ErrorInvalidSubject     = &management.Error{HttpStatus: http.StatusBadRequest, Id: ERR_INVALID_SUBJECT, Description: "The specified certificate contains an invalid subject."}
	ErrorLoginNameNotFound  = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service was unable to find information about the authenticated user."}
	ErrorUserNameNotFound   = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service was unable to find the user name of the authenticated user."}
	ErrorDeviceRegistered   = &management.Error{HttpStatus: http.StatusConflict, Id: ERR_REGISTERED, Description: "The specified device cannot be deleted. It is being used by a user account."}
	ErrorDecodingBase64String = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "Unable to decode base64 string."}
)
