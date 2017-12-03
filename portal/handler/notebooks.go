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
	"strings"
)

// CreateNotebook creates a new notebook.
func (controller *Controller) CreateNotebook(context *gin.Context) {
	mlog.Debug("CreateNotebook")

	// Get the token for current request.
	token, found := controller.getAuthToken(context)

	if !found {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	notebook := &model.Notebook{}

	// Get resource from request object.
	if err := controller.Bind(context, notebook); err != nil {
		mlog.Error("Failed to create runtime with error: %v.", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Create the notebook.
	notebook, mErr := controller.portalProvider.CreateNotebook(token.AccessToken, notebook)
	if mErr != nil {
		mlog.Error("Failed to create notebook with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, notebook)
}

// UpdateNotebook updates an existing notebook.
func (controller *Controller) UpdateNotebook(context *gin.Context) {
	mlog.Debug("UpdateNotebook")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	notebook := &model.Notebook{}

	// Get resource from request object.
	if err := controller.Bind(context, notebook); err != nil {
		mlog.Error("Failed to update notebook with error: %v.", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	if mErr := controller.portalProvider.UpdateNotebook(token.AccessToken, notebook); mErr != nil {
		mlog.Error("Failed to update notebook with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// GetNotebook returns the notebook with the specified id.
func (controller *Controller) GetNotebook(context *gin.Context) {
	mlog.Debug("GetNotebook")
	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get request token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	id := context.Params.ByName("notebookId")
	if id == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid notebook ID")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("notebook ID"))
		return
	}

	notebook, mErr := controller.portalProvider.GetNotebook(token.AccessToken, id)
	if mErr != nil {
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, notebook)
}

// ListNotebook returns notebook associated with authenticated user.
func (controller *Controller) ListNotebooks(context *gin.Context) {
	mlog.Debug("ListNotebook")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get request token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	notebooks, mErr := controller.portalProvider.ListNotebooks(token.AccessToken)
	if mErr != nil {
		mlog.Error("Couldn't get notebook list. Error was: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, notebooks)
}

// DeleteNotebook deletes the notebook with the specified id.
func (controller *Controller) DeleteNotebook(context *gin.Context) {
	mlog.Debug("DeleteNotebook")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	id := context.Params.ByName("notebookId")
	if id == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid notebook id")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("notebook ID"))
		return
	}

	if mErr := controller.portalProvider.DeleteNotebook(token.AccessToken, id); mErr != nil {
		mlog.Error("Failed to delete notebook with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

func (controller *Controller) GetNotebookUsers(context *gin.Context) {
	mlog.Debug("GetNotebookUsers")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	//get the notebook id
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))
	if notebookId == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid notebook ID")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("notebook ID"))
		return
	}

	users, mErr := controller.portalProvider.GetNotebookUsers(token.AccessToken, notebookId)
	if mErr != nil {
		mlog.Error("Failed to get notebook users with error: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, users)
}

func (controller *Controller) UpdateNotebookUsers(context *gin.Context) {
	mlog.Debug("UpdateNotebookUsers")

	//get our token
	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	//get the notebook id
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))
	if notebookId == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid notebook ID")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("notebook ID"))
		return
	}

	users := []model.User{}
	if err := controller.Bind(context, &users); err != nil {
		mlog.Error("Failed to parse users with error: %s.", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	serviceErr := controller.portalProvider.UpdateNotebookUsers(token.AccessToken, notebookId, users)
	if serviceErr != nil {
		mlog.Error("Failed to update users with error: %s", serviceErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

}
