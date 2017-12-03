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
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/config"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

// Creates a new transformation.
func (controller *Controller) CreateTransformation(context *gin.Context) {
	mlog.Info("CreateTransformation")

	if config.EnforceChecksum {
		utils.ErrCreateTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorOperationDisabled)
		return
	}

	// Get the resource.
	var transformation model.Transformation

	// Validate request message
	if err := controller.Bind(context, &transformation); err != nil {
		mlog.Error("Failed to create transformation with error: %v.", err)
		utils.ErrCreateTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrCreateTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	if err := transformation.Validate(); err != nil {
		utils.ErrCreateTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Create the resource.
	createdTransformation, mErr := controller.transformationProvider.Create(accountId, &transformation)

	if mErr != nil {
		mlog.Error("Failed to create transformation with error: %v", mErr)
		utils.ErrCreateTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.CreateTransformation.Incr()
	context.JSON(http.StatusOK, createdTransformation)
}

// Returns the transformation with the specified id.
func (controller *Controller) GetTransformation(context *gin.Context) {
	mlog.Info("GetTransformation")

	// Get resource id.
	transformationId := strings.TrimSpace(context.Params.ByName("transformationId"))

	if transformationId == "" {
		mlog.Error("Failed to get transformation due to bad request - Invalid resource id.")
		utils.ErrGetTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrGetTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Get the resource.
	transformation, mErr := controller.transformationProvider.Get(accountId, transformationId)

	if mErr != nil {
		mlog.Error("Failed to get transformation with error: %v", mErr)
		utils.ErrGetTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.GetTransformation.Incr()
	context.JSON(http.StatusOK, transformation)
}

func (controller *Controller) GetTransformationResults(context *gin.Context) {
	mlog.Info("GetTransformationResults")

	//Get transformationID
	transformationId := strings.TrimSpace(context.Params.ByName("transformationId"))
	if transformationId == "" {
		mlog.Error("Failed to get transformation results due to bad request -- Invalid resource id.")
		utils.ErrTransformationResults.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	accountId, mErr := controller.getAccountId(context)
	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrTransformationResults.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	results, mErr := controller.transformationProvider.Results(accountId, transformationId)
	if mErr != nil {
		mlog.Error("Failed to get transformation results with error: %v", mErr)
		utils.ErrTransformationResults.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.TransformationResults.Incr()
	context.JSON(http.StatusOK, results)
}

// Returns the transformations associated with the authenticated user.
func (controller *Controller) ListTransformations(context *gin.Context) {
	mlog.Info("ListTransformations")

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrListTransformations.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Get resource for account id.
	transformations, mErr := controller.transformationProvider.List(accountId)

	if mErr != nil {
		mlog.Error("Failed to list transformations with error: %v", mErr)
		utils.ErrListTransformations.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ListTransformations.Incr()
	context.JSON(http.StatusOK, transformations)
}

// Updates the transformation with the specified id.
func (controller *Controller) UpdateTransformation(context *gin.Context) {
	mlog.Info("UpdateTransformation")

	if config.EnforceChecksum {
		utils.ErrUpdateTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorOperationDisabled)
		return
	}

	// Get the resource.

	var transformation model.Transformation

	// Validate request message
	if err := controller.Bind(context, &transformation); err != nil {
		mlog.Error("Failed to update transformation with error: %v.", err)
		utils.ErrUpdateTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)
	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrUpdateTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	if err := transformation.Validate(); err != nil {
		utils.ErrUpdateTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Update the resource.
	if mErr := controller.transformationProvider.Update(accountId, &transformation); mErr != nil {
		mlog.Error("Failed to update transformation with error: %v", mErr)
		utils.ErrUpdateTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.UpdateTransformation.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Deletes the transformation with the specified id.
func (controller *Controller) DeleteTransformation(context *gin.Context) {
	mlog.Info("DeleteTransformation")

	// Get the resource id.
	transformationId := strings.TrimSpace(context.Params.ByName("transformationId"))

	if transformationId == "" {
		mlog.Error("Failed to delete transformation due to bad request - Invalid resource id.")
		utils.ErrDeleteTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrDeleteTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Delete the resource with id.
	if mErr = controller.transformationProvider.Delete(accountId, transformationId); mErr != nil {
		mlog.Error("Failed to delete transformation with error: %v", mErr)
		utils.ErrDeleteTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.DeleteTransformation.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Updates the transformation with the specified id.
func (controller *Controller) ExecuteTransformation(context *gin.Context) {
	mlog.Info("ExecuteTransformation")

	if config.EnforceChecksum {
		utils.ErrExecuteTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorOperationDisabled)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrExecuteTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Get callback url.
	callbackUrl, err := controller.getCallbackUrl(context)

	if err != nil {
		mlog.Error("Failed to execute transformation action with error: %v.", err)
		utils.ErrExecuteTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidCallbackUrl)
		return
	}

	// Get the resource.
	transformation := &model.Transformation{}

	// Validate request message
	if err := controller.Bind(context, transformation); err != nil {
		mlog.Error("Failed to execute transformation action with error: %v.", err)
		utils.ErrExecuteTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	if err := transformation.Validate(); err != nil {
		utils.ErrExecuteTransformation.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	executionRequest, mErr := controller.executionProvider.ExecuteTransformation(user, transformation, callbackUrl)
	if mErr != nil {
		mlog.Error("Execute transformation returned error: %v", mErr)
		utils.ErrExecuteTransformation.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ExecuteTransformation.Incr()
	context.JSON(http.StatusAccepted, executionRequest)
}
