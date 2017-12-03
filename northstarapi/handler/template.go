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

// Creates a new template.
func (controller *Controller) CreateTemplate(context *gin.Context) {
	mlog.Info("CreateTemplate")

	if config.EnforceChecksum {
		utils.ErrCreateTemplate.Incr()
		controller.RenderServiceError(context, model.ErrorOperationDisabled)
		return
	}

	// Get the resource
	template := &model.Template{}

	// Validate request message
	if err := controller.Bind(context, template); err != nil {
		mlog.Error("Failed to create template with error: %v.", err)
		utils.ErrCreateTemplate.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrCreateTemplate.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	createdTemplate, mErr := controller.templateProvider.Create(user, template)

	if mErr != nil {
		mlog.Error("Failed to create template with error: %v", mErr)
		utils.ErrCreateTemplate.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.CreateTemplate.Incr()
	context.JSON(http.StatusOK, createdTemplate)
}

// Returns the templates for the authenticated user.
func (controller *Controller) ListTemplates(context *gin.Context) {
	mlog.Info("ListTemplates")

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrListTemplates.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	templates, mErr := controller.templateProvider.List(user)

	if mErr != nil {
		mlog.Error("Failed to get templates with error: %v", mErr)
		utils.ErrListTemplates.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ListTemplates.Incr()
	context.JSON(http.StatusOK, templates)
}

// Returns the template with the specified id.
func (controller *Controller) GetTemplate(context *gin.Context) {
	mlog.Info("GetTemplate")

	// Get resource id.
	templateId := strings.TrimSpace(context.Params.ByName("templateId"))

	if templateId == "" {
		mlog.Error("Failed to get template due to bad request - Invalid resource id.")
		utils.ErrGetTemplate.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrGetTemplate.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	template, mErr := controller.templateProvider.Get(user, templateId)

	if mErr != nil {
		mlog.Error("Failed to get template with error: %v", mErr)
		utils.ErrGetTemplate.Incr()
		controller.RenderServiceError(context, management.ErrorExternal)
		return
	}

	utils.GetTemplate.Incr()
	context.JSON(http.StatusOK, template)
}

// Deletes the templaye with the specified id.
func (controller *Controller) DeleteTemplate(context *gin.Context) {
	mlog.Info("DeleteTemplate")

	// Get resource id.
	templateId := strings.TrimSpace(context.Params.ByName("templateId"))

	if templateId == "" {
		mlog.Error("Failed to delete template due to bad request - Invalid resource id.")
		utils.ErrDeleteTemplate.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrDeleteTemplate.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	// Delete the resource with id.
	if mErr = controller.templateProvider.Delete(user, templateId); mErr != nil {
		mlog.Error("Failed to delete template with error: %v", mErr)
		utils.ErrDeleteTemplate.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.DeleteTemplate.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
