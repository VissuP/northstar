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
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/middleware"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/provider"
	"github.com/verizonlabs/northstar/northstarapi/provider/northstar"
	"github.com/verizonlabs/northstar/northstarapi/provider/northstar/object"
	"github.com/verizonlabs/northstar/northstarapi/provider/thingspace"
)

const (
	// Defines the supported custom headers.
	CallBackHeader = "x-vz-callback-url"
)

// Defines the service controller. This type implements the
// HTTP handler methods for the service API.
type Controller struct {
	accountProvider        provider.AccountProvider
	transformationProvider provider.TransformationProvider
	notebookProvider       provider.NotebookProvider
	executionProvider      provider.ExecutionProvider
	templateProvider       provider.TemplateProvider
	objectProvider         provider.ObjectProvider
	streamProvider         provider.StreamProvider
}

// Returns a new Controller.
func NewController() (*Controller, error) {
	mlog.Info("NewController")

	// Create the ThingSpace Account Provider.
	accountProvider, err := thingspace.NewThingSpaceAccountProvider()

	if err != nil {
		return nil, fmt.Errorf("Failed to create thingspace account provider with error: %+v", err)
	}

	// Create the NorthStar Providers.
	notebookProvider, err := northstar.NewNorthStarNotebooksProvider()

	if err != nil {
		return nil, fmt.Errorf("Failed to create notebook provider: %+v", err)
	}

	transformationProvider, err := northstar.NewNorthStarTransformationProvider()

	if err != nil {
		return nil, fmt.Errorf("Failed to create transformaton provider: %+v", err)
	}

	templateProvider, err := northstar.NewNorthStarTemplatesProvider()

	if err != nil {
		return nil, fmt.Errorf("Failed to create templates provider: %+v", err)
	}

	// Create execution provider.
	executionProvider, err := northstar.NewNorthStarExecutionProvider()

	if err != nil {
		return nil, fmt.Errorf("Failed to create execution provider with error: %+v", err)
	}

	//Create object provider
	objectProvider, err := object.NewS3ObjectProvider()
	if err != nil {
		return nil, fmt.Errorf("Failed to create object provider with error: %+v", err)
	}

	streamProvider, err := northstar.NewNorthstarStreamsProvider()
	if err != nil {
		return nil, fmt.Errorf("Failed to create jobs provider with error: %+v", err)
	}

	// Create the controller
	controller := &Controller{
		accountProvider:        accountProvider,
		transformationProvider: transformationProvider,
		notebookProvider:       notebookProvider,
		executionProvider:      executionProvider,
		templateProvider:       templateProvider,
		objectProvider:         objectProvider,
		streamProvider:         streamProvider,
	}

	return controller, nil
}

// Helper method used to get the account id of the authenticated user.
func (controller *Controller) getAccountId(context *gin.Context) (string, *management.Error) {
	mlog.Debug("getAccountId")

	// Get the loginname of the authenticated user.
	loginname, exists := context.Get(middleware.LOGINNAME_KEY)

	if exists == false {
		mlog.Error("Failed to get authenticated user loginname from context.")
		return "", model.ErrorLoginNameNotFound
	}

	// Get the account id for the loginname.
	accountId, mErr := controller.accountProvider.GetAccountIdForLoginname(loginname.(string))

	if mErr != nil {
		mlog.Error("Failed to get account id from loginname with error: %s", mErr.Description)
		return "", mErr
	}

	return accountId, nil
}

// Helper method used to get authenticated user.
func (controller *Controller) getUser(context *gin.Context) (*model.User, *management.Error) {
	mlog.Debug("getUser")

	// Get the loginname of the authenticated user.
	loginname, exists := context.Get(middleware.LOGINNAME_KEY)

	if exists == false {
		mlog.Error("Failed to get authenticated user loginname from context.")
		return nil, model.ErrorLoginNameNotFound
	}

	// Get user associated with loginame.
	user, mErr := controller.accountProvider.GetUser(loginname.(string))

	if mErr != nil {
		mlog.Error("Failed to get user from loginname with error: %s", mErr.Description)
		return nil, mErr
	}

	return user, nil
}

// TODO - Move this method to the management library.

// Helper method used to render http response from given management error object.
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

// Helper method used to bind body based on supported content types.
func (controller *Controller) Bind(context *gin.Context, resource interface{}) error {
	request := context.Request
	bind := controller.getBinding(request.Method, gin.MIMEJSON)

	if err := bind.Bind(request, resource); err != nil {
		return fmt.Errorf("Failed to bind request body with error: %v", err)
	}

	return nil
}

// Helper method used to get content binding from content type.
func (controller *Controller) getBinding(method, contentType string) binding.Binding {
	mlog.Debug("getBinding: method:%s, contentType:%s", method, contentType)

	if method == "GET" {
		return binding.Form
	}

	// TODO - Add return bu supported content types. For now, assuming JSON.
	return binding.JSON
}

// Helper method used to get callback url header.
func (controller *Controller) getCallbackUrl(context *gin.Context) (string, error) {
	mlog.Debug("getCallbackUrl")

	// Get header
	request := context.Request
	value := request.Header.Get(CallBackHeader)

	if value == "" {
		return "", fmt.Errorf("Callback url is empty")
	}

	// Validate the value.
	if _, err := url.Parse(value); err != nil {
		return "", fmt.Errorf("Parse url returned error: %+v", err)
	}

	return value, nil
}
