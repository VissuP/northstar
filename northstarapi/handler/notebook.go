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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/config"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

// Creates a new notebook.
func (controller *Controller) CreateNotebook(context *gin.Context) {
	mlog.Info("CreateNotebook")

	// Get the resource.
	notebook := &model.Notebook{}

	// Validate request message
	if err := controller.Bind(context, notebook); err != nil {
		mlog.Error("Failed to create notebook with error: %v.", err)
		utils.ErrCreateNotebook.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	if err := notebook.Validate(); err != nil {
		mlog.Error("Failed to validate notebook with error %v.", err)
		utils.ErrCreateNotebook.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidNotebookModel)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrCreateNotebook.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	createdNotebook, mErr := controller.notebookProvider.Create(user, notebook)

	if mErr != nil {
		mlog.Error("Failed to create notebook with error: %v", mErr)
		utils.ErrCreateNotebook.Incr()
		controller.RenderServiceError(context, management.ErrorExternal)
		return
	}

	utils.CreateNotebook.Incr()
	context.JSON(http.StatusOK, createdNotebook)
}

// GetNotebook returns the notebook with the specified id.
func (controller *Controller) GetNotebook(context *gin.Context) {
	mlog.Info("GetNotebook")

	// Get resource id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	if notebookId == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid resource id.")
		utils.ErrGetNotebook.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrGetNotebook.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	notebook, mErr := controller.notebookProvider.Get(user, notebookId)

	if mErr != nil {
		mlog.Error("Failed to get notebook with error: %v", mErr)
		utils.ErrGetNotebook.Incr()
		controller.RenderServiceError(context, management.ErrorExternal)
		return
	}

	utils.GetNotebook.Incr()
	context.JSON(http.StatusOK, notebook)
}

// GetNotebook returns the users of the notebook with the specified id.
func (controller *Controller) GetNotebookUsers(context *gin.Context) {
	mlog.Info("GetNotebookUsers")

	// Get resource id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	if notebookId == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid resource id.")
		utils.ErrGetUsers.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrGetUsers.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Get collection of users associated with notebook id.
	users, mErr := controller.notebookProvider.GetUsers(user, notebookId)

	if mErr != nil {
		mlog.Error("Failed to get notebook users with error: %v", mErr)
		utils.ErrGetUsers.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Note that the collection returned by notebook provider does
	// not include the full user information. We need to request
	// additional information from ThingSpace.
	for index := 0; index < len(users); index++ {
		user, mErr := controller.accountProvider.GetUserById(users[index].Id)

		if mErr != nil {
			mlog.Error("Failed to get user information with error: %v", mErr)
			utils.ErrGetUsers.Incr()
			controller.RenderServiceError(context, mErr)
			return
		}

		users[index].DisplayName = user.DisplayName
		users[index].Email = user.Email
		users[index].ImageId = user.ImageId
	}

	utils.GetUsers.Incr()
	context.JSON(http.StatusOK, users)
}

// UpdateNotebookUsers updates the users associated with the specified notebook id.
func (controller *Controller) UpdateNotebookUsers(context *gin.Context) {
	mlog.Info("UpdateNotebookUsers")

	// Get resource id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	if notebookId == "" {
		mlog.Error("Failed to get notebook due to bad request - Invalid resource id.")
		utils.ErrUpdateUsers.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get the resource.
	users := []model.User{}

	// Validate request message
	if err := controller.Bind(context, &users); err != nil {
		mlog.Error("Failed to update notebook users with error: %v.", err)
		utils.ErrUpdateUsers.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrUpdateUsers.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	if mErr := controller.notebookProvider.UpdateUsers(user, notebookId, users); mErr != nil {
		mlog.Error("Failed to update notebook users with error: %v", mErr)
		utils.ErrUpdateUsers.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.UpdateUsers.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// ListNotebooks returns the notebooks for the authenticated user.
func (controller *Controller) ListNotebooks(context *gin.Context) {
	mlog.Info("ListNotebooks")

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrListNotebooks.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	notebooks, mErr := controller.notebookProvider.List(user)

	if mErr != nil {
		mlog.Error("Failed to get notebooks with error: %v", mErr)
		utils.ErrListNotebooks.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ListNotebooks.Incr()
	context.JSON(http.StatusOK, notebooks)
}

// UpdateNotebook updates the notebook with the specified id.
func (controller *Controller) UpdateNotebook(context *gin.Context) {
	mlog.Info("UpdateNotebook")

	// Get the resource.
	notebook := &model.Notebook{}

	// Validate request message
	if err := controller.Bind(context, notebook); err != nil {
		mlog.Error("Failed to update notebook with error: %v.", err)
		utils.ErrUpdateNotebook.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	if err := notebook.Validate(); err != nil {
		mlog.Error("Failed to validate notebook with error %v.", err)
		utils.ErrUpdateNotebook.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidNotebookModel)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrUpdateNotebook.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	if mErr := controller.notebookProvider.Update(user, notebook); mErr != nil {
		mlog.Error("Failed to update notebook with error: %v", mErr)
		utils.ErrUpdateNotebook.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.UpdateNotebook.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Deletes the notebook with the specified id.
func (controller *Controller) DeleteNotebook(context *gin.Context) {
	mlog.Info("DeleteNotebook")

	// Get resource id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	if notebookId == "" {
		mlog.Error("Failed to delete notebook due to bad request - Invalid resource id.")
		utils.ErrDeleteNotebook.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrDeleteNotebook.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Delete the resource with id.
	if mErr = controller.notebookProvider.Delete(user, notebookId); mErr != nil {
		mlog.Error("Failed to delete notebook with error: %v", mErr)
		utils.ErrDeleteNotebook.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.DeleteNotebook.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Executes a specified notebook cell.
func (controller *Controller) ExecuteNotebookCell(context *gin.Context) {
	mlog.Info("ExecuteNotebookCell")

	// Get resource id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	if notebookId == "" {
		mlog.Error("Failed to execute notebook cell due to bad request - Invalid resource id.")
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get resource id.
	cellId := strings.TrimSpace(context.Params.ByName("cellId"))

	if cellId == "" {
		mlog.Error("Failed to execute notebook cell due to bad request - Invalid resource id.")
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get the resource.
	cell := &model.Cell{}

	// Validate request message
	if err := controller.Bind(context, cell); err != nil {
		mlog.Error("Failed to execute notebook cell with error: %v.", err)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	if config.EnforceChecksum {
		mErr := controller.templateProvider.TemplateExists(cell.Input.Body)
		if mErr != nil {
			mlog.Error("Error, code does not match any template. Cannot execute.")
			controller.RenderServiceError(context, mErr)
			return
		}
	}

	if err := cell.Validate(); err != nil {
		mlog.Error("Failed to validate cell with error %v.", err)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, management.GetInternalError(err.Error()))
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Get callback url.
	callbackUrl, err := controller.getCallbackUrl(context)

	if err != nil {
		mlog.Error("Failed to execute notebook cell with error: %v.", err)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidCallbackUrl)
		return
	}

	// Get the notebook owner.
	executionUser, mErr := controller.notebookProvider.GetExecutionInformation(user, notebookId)

	if mErr != nil {
		mlog.Error("Get execution information returned error: %v", mErr)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Execute the cell in the context of the owner user.
	executionRequest, mErr := controller.executionProvider.ExecuteCell(executionUser, cell, callbackUrl)
	if mErr != nil {
		mlog.Error("Execute cell returned error: %v", mErr)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ExecuteNotebookCell.Incr()
	context.JSON(http.StatusOK, executionRequest)
}

// Executes a specific cell action.
func (controller *Controller) ExecuteCell(context *gin.Context) {
	mlog.Info("ExecuteCell")

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrExecuteCell.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Get callback url.
	callbackUrl, err := controller.getCallbackUrl(context)

	if err != nil {
		mlog.Error("Failed to execute cell with error: %v.", err)
		utils.ErrExecuteCell.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidCallbackUrl)
		return
	}

	// Get the resource.
	cell := &model.Cell{}

	// Validate request message
	if err := controller.Bind(context, cell); err != nil {
		mlog.Error("Failed to execute cell with error: %v.", err)
		utils.ErrExecuteCell.Incr()
		controller.RenderServiceError(context, management.GetBadRequestError(err.Error()))
		return
	}

	if config.EnforceChecksum {
		mErr = controller.templateProvider.TemplateExists(cell.Input.Body)
		if mErr != nil {
			mlog.Error("Error, code does not match any template. Cannot execute: %s", mErr.Description)
			controller.RenderServiceError(context, mErr)
			return
		}
	}

	if err := cell.Validate(); err != nil {
		mlog.Error("Failed to validate cell with error %v.", err)
		utils.ErrExecuteCell.Incr()
		controller.RenderServiceError(context, management.GetInternalError(err.Error()))
		return
	}

	executionRequest, mErr := controller.executionProvider.ExecuteCell(user, cell, callbackUrl)
	if mErr != nil {
		mlog.Error("Execute cell returned error: %v", mErr)
		utils.ErrExecuteCell.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ExecuteCell.Incr()
	context.JSON(http.StatusOK, executionRequest)
}
