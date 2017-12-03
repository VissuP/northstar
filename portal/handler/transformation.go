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

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/portal/model"
)

// CreateTransformation creates a new transformation.
func (controller *Controller) CreateTransformation(context *gin.Context) {
	mlog.Debug("CreateTransformation")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	transformation := &model.Transformation{}

	// Get resource from request object.
	if err := controller.Bind(context, transformation); err != nil {
		mlog.Error("Failed to create transformation with error: %v.", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	transformation, mErr := controller.portalProvider.CreateTransformation(token.AccessToken, transformation)
	if mErr != nil {
		mlog.Error("Failed to create transformation with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, transformation)
}

// UpdateTransformation updates an existing transformation.
func (controller *Controller) UpdateTransformation(context *gin.Context) {
	mlog.Debug("UpdateTransformation")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	transformation := &model.Transformation{}

	// Get resource from request object.
	if err := controller.Bind(context, transformation); err != nil {
		mlog.Error("Failed to update transformation with error: %v.", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	if transformation.Id == "" {
		mlog.Error("Failed to update transformation due to bad request -- Invalid transformation ID.")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("transformation ID"))
		return
	}

	if mErr := controller.portalProvider.UpdateTransformation(token.AccessToken, transformation); mErr != nil {
		mlog.Error("Failed to update transformation with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// GetTransformation returns the transformation with the specified id.
func (controller *Controller) GetTransformation(context *gin.Context) {
	mlog.Debug("GetTransformation")

	id := context.Params.ByName("id")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	transformation, mErr := controller.portalProvider.GetTransformation(token.AccessToken, id)
	if mErr != nil {
		mlog.Error("Failed to get transformation by id with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, transformation)
}

// GetTransformationResults returns the execution results for the specified transformation ID.
func (controller *Controller) GetTransformationResults(context *gin.Context) {
	mlog.Debug("GetTransformationResults")

	transformationID := context.Params.ByName("id")
	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	results, mErr := controller.portalProvider.GetTransformationResults(token.AccessToken, transformationID)
	if mErr != nil {
		mlog.Error("Failed to get transformation results with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, results)
}

// ListTransformations returns transformations associated with authenticated user.
func (controller *Controller) ListTransformations(context *gin.Context) {
	mlog.Debug("ListTransformations")

	token, tokenExists := controller.getAuthToken(context)

	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	transformations, mErr := controller.portalProvider.ListTransformations(token.AccessToken)

	if mErr != nil {
		mlog.Error("Failed to get transformation table with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, transformations)
}

// DeleteTransformation deletes the transformation with the specified id.
func (controller *Controller) DeleteTransformation(context *gin.Context) {
	mlog.Debug("DeleteTransformation")

	id := context.Params.ByName("id")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	if mErr := controller.portalProvider.DeleteTransformation(token.AccessToken, id); mErr != nil {
		mlog.Error("Failed to delete transformation with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// GetEventSchemas returns the supported schedule event schemas
func (controller *Controller) GetScheduleEventSchemas(context *gin.Context) {
	mlog.Info("GetScheduleEventSchemas")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	schemas, mErr := controller.portalProvider.GetScheduleEventSchemas(token.AccessToken)
	if mErr != nil {
		mlog.Error("Failed to get schedule event schemas with error: %v", mErr)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, schemas)
}
