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
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

// Creates a new schedule for specified transformation id.
func (controller *Controller) CreateSchedule(context *gin.Context) {
	mlog.Info("CreateSchedule")

	// Get the resource.
	var schedule model.Schedule

	// Validate request message
	if err := controller.Bind(context, &schedule); err != nil {
		mlog.Error("Failed to create schedule with error: %v.", err)
		utils.ErrCreateSchedule.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Get resource id.
	transformationId := strings.TrimSpace(context.Params.ByName("transformationId"))

	if transformationId == "" {
		mlog.Error("Failed to create schedule due to bad request - Invalid resource id.")
		utils.ErrCreateSchedule.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrCreateSchedule.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	if mErr := controller.transformationProvider.CreateSchedule(accountId, transformationId, &schedule); mErr != nil {
		mlog.Error("Failed to create schedule with error: %v", mErr)
		utils.ErrCreateSchedule.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.CreateSchedule.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Gets the schedule for specified transformation id.
func (controller *Controller) GetSchedule(context *gin.Context) {
	mlog.Info("GetSchedule")

	// Get resource id.
	transformationId := strings.TrimSpace(context.Params.ByName("transformationId"))

	if transformationId == "" {
		mlog.Error("Failed to get schedule due to bad request - Invalid resource id.")
		utils.ErrGetSchedule.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrGetSchedule.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	schedule, mErr := controller.transformationProvider.GetSchedule(accountId, transformationId)

	if mErr != nil {
		mlog.Error("Failed to get schedule with error: %v", mErr)
		utils.ErrGetSchedule.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.GetSchedule.Incr()
	context.JSON(http.StatusOK, schedule)
}

// Deletes the schedule for specified transformation id.
func (controller *Controller) DeleteSchedule(context *gin.Context) {
	mlog.Info("DeleteSchedule")

	// Get resource id.
	transformationId := strings.TrimSpace(context.Params.ByName("transformationId"))

	if transformationId == "" {
		mlog.Error("Failed to delete schedule due to bad request - Invalid resource id.")
		utils.ErrDeleteSchedule.Incr()
		controller.RenderServiceError(context, model.ErrorInvalidResourceId)
		return
	}

	// Get account id.
	accountId, mErr := controller.getAccountId(context)

	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrDeleteSchedule.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	if mErr := controller.transformationProvider.DeleteSchedule(accountId, transformationId); mErr != nil {
		mlog.Error("Failed to delete schedule with error: %v", mErr)
		utils.ErrDeleteSchedule.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.DeleteSchedule.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
