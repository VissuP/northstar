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

package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/portal/config"
	"github.com/verizonlabs/northstar/portal/middleware"
	"github.com/verizonlabs/northstar/portal/provider"
	"github.com/verizonlabs/northstar/portal/provider/northstar"
	"github.com/verizonlabs/northstar/portal/utils"
	north "github.com/verizonlabs/northstar/portal/utils/thingspace"
)

// Controller defines the structure for a portal controller.
type Controller struct {
	authClient     thingspace.AuthClient
	userClient     *north.UserClient
	clientId       string
	clientSecret   string
	userScopes     string
	portalProvider provider.PortalProvider
	writers        *utils.ThreadSafeMap
}

// NewController returns a new controller
func NewController() (*Controller, error) {
	mlog.Debug("NewController")
	portalProvider, err := northstar.NewNorthStarPortalProvider(config.Configuration.NorthstarAPIProtocol, config.Configuration.NorthstarAPIHostPort)
	if err != nil {
		return nil, err
	}

	// Controller
	controller := &Controller{
		authClient:     thingspace.NewThingSpaceAuthClientWithProtocol(config.Configuration.ThingspaceProtocol, config.Configuration.ThingSpaceAuthHostPort),
		userClient:     north.NewUserClient(config.Configuration.ThingspaceProtocol, config.Configuration.ThingSpaceUserHostPort),
		clientId:       config.Configuration.ThingSpaceClientID,
		clientSecret:   config.Configuration.ThingSpaceClientSecret,
		userScopes:     config.Configuration.ThingSpaceUserScopes,
		portalProvider: portalProvider,
		writers:        utils.NewThreadSafeMap(),
	}

	return controller, nil
}

// NotFound is a handler used for not found errors.
func (controller *Controller) NotFound(context *gin.Context) {
	mlog.Info("NotFound")
	controller.RenderServiceError(context, management.ErrorNotFound)
}

// RenderServiceError is a helper method used to render http response from given management error object.
func (controller *Controller) RenderServiceError(context *gin.Context, serviceError *management.Error) {
	// per docs, headers need to be set before calling context.JSON method
	for k, v := range serviceError.Header {
		for _, v1 := range v {
			context.Writer.Header().Add(k, v1)
		}
	}
	// now serialize rest of the response
	context.JSON(serviceError.HttpStatus, serviceError)
}

// Bind is a helper method used to bind body based on supported content types.
func (controller *Controller) Bind(context *gin.Context, resource interface{}) error {
	request := context.Request
	bind := controller.getBinding(request.Method, gin.MIMEJSON)

	if err := bind.Bind(request, resource); err != nil {
		return fmt.Errorf("Failed to bind request body with error: %v", err)
	}

	return nil
}

// getBinding is a helper method used to get content binding from content type.
func (controller *Controller) getBinding(method, contentType string) binding.Binding {
	mlog.Debug("getBinding: method:%s, contentType:%s", method, contentType)

	if method == "GET" {
		return binding.Form
	}

	// TODO - Add return by supported content types. For now, assuming JSON.
	return binding.JSON
}

// getAuthToken is a helper method used to return token from the request.
func (controller *Controller) getAuthToken(context *gin.Context) (*api.Token, bool) {
	mlog.Debug("getAuthToken")

	tokenInfo, found := context.Get(middleware.AccessTokenKeyName)
	if !found {
		mlog.Error("Error, auth token was not found in gin context.")
		return nil, false
	}
	token, valid := tokenInfo.(*api.Token)
	if !valid {
		mlog.Error("Error, could not convert interface to auth token.")
		return nil, false
	}
	return token, true
}

func (controller *Controller) GetIndexHTML(context *gin.Context) {
	context.File("web/dist/index.html")
}
