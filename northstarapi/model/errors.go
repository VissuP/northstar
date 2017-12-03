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

const (
	ERR_CONFLICT string = "resource_conflict"
)

var (
	// Defined user specific errors returned by the API.
	ErrorParseRequestBody     = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The request body is invalid. Failed to parse content."}
	ErrorInvalidResourceId    = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The resource id is missing or invalid."}
	ErrorInvalidCallbackUrl   = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The callback url is missing or invalid."}
	ErrorInvalidNotebookModel = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The notebook model is invalid or not supported."}
	ErrorInvalidEventCategory = &management.Error{HttpStatus: http.StatusBadRequest, Id: management.ERR_BAD_REQUEST, Description: "The event category is invalid or not supported."}

	ErrorLoginNameNotFound     = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service was unable to find information about the authenticated user."}
	ErrorToExternalNotebook    = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service found an error while marshaling notebook."}
	ErrorFromExternalNotebook  = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service found an error while unmarshaling internal notebook."}
	ErrorFromExternalSnippet   = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service found an error while unmarshaling internal transformation."}
	ErrorNotebookOwnerNotFound = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service found an error while getting notebook access information. Owner not found."}
	ErrorToExternalTemplate    = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service found an error while marshaling template."}
	ErrorFromExternalTemplate  = &management.Error{HttpStatus: http.StatusInternalServerError, Id: management.ERR_SERVICE_ERROR, Description: "The service found an error while unmarshaling template."}

	ErrorNoNotebookPermissions    = &management.Error{HttpStatus: http.StatusForbidden, Id: management.ERR_FORBIDDEN, Description: "The user is not authorized to access or update this notebook."}
	ErrorNoNotebookUsrPermission  = &management.Error{HttpStatus: http.StatusForbidden, Id: management.ERR_FORBIDDEN, Description: "The user is not authorized to access or update notebook users."}
	ErrorNoNotebookExecPermission = &management.Error{HttpStatus: http.StatusForbidden, Id: management.ERR_FORBIDDEN, Description: "The user is not authorized to execute notebook."}

	ErrorTransformationScheduled = &management.Error{HttpStatus: http.StatusConflict, Id: ERR_CONFLICT, Description: "The request could not be completed due to conflict with current transformation scheduled state. E.g., scheduled transformation can not be updated, deleted, etc."}
	ErrorOperationDisabled       = &management.Error{HttpStatus: http.StatusForbidden, Id: management.ERR_FORBIDDEN, Description: "This operation is forbidden in the current environment."}
)
