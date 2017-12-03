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
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/config"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
	"strconv"
)

//TriggerExecution triggers an execution
func (controller *Controller) TriggerExecution(context *gin.Context) {
	mlog.Info("Execute")

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrTriggerExecution.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	executionRequest := &model.ExecutionRequest{}
	if err := controller.Bind(context, executionRequest); err != nil {
		mlog.Error("Failed to trigger execution. Resource validation failed. %v", err.Error())
		utils.TriggerExecution.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	//Verify the checksum if it's required.
	if config.EnforceChecksum {
		mErr := controller.templateProvider.TemplateExists(executionRequest.Code)
		if mErr != nil {
			mlog.Error("Failed to trigger execution. Custom code execution is disabled.")
			controller.RenderServiceError(context, mErr)
			return
		}
	}

	//Execute the code under the current user's permissions
	updatedExecutionRequest, mErr := controller.executionProvider.Execute(user.AccountId, executionRequest)
	if mErr != nil {
		mlog.Error("Execute cell returned error: %v", mErr)
		utils.ErrExecuteNotebookCell.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	context.JSON(http.StatusOK, updatedExecutionRequest)
}

//ListExecutions lists all executions for the account
func (controller *Controller) ListExecutions(context *gin.Context) {
	mlog.Info("ListExecutions")

	limit, err := strconv.Atoi(context.GetHeader("Limit"))
	if err != nil {
		mlog.Error("Invalid limit %v.", err.Error())
		utils.ErrListExecutions.Incr()
		controller.RenderServiceError(context, management.GetBadRequestError("Invalid limit."))
		return
	}

	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrListExecutions.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	results, mErr := controller.executionProvider.ListExecutions(user.AccountId, limit)
	if mErr != nil {
		mlog.Error("Failed to list executions. %+v", mErr.Description)
		utils.ErrListExecutions.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.ListExecutions.Incr()
	context.JSON(http.StatusOK, results)
}

//GetExecution gets the requested execution
func (controller *Controller) GetExecution(context *gin.Context) {
	mlog.Info("GetExecution")

	id := context.Params.ByName("id")
	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrGetExecution.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	result, mErr := controller.executionProvider.GetExecution(user.AccountId, id)
	if mErr != nil {
		utils.ErrGetExecution.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.GetExecution.Incr()
	context.JSON(http.StatusOK, result)

}

//StopExecution stops the specified execution
func (controller *Controller) StopExecution(context *gin.Context) {
	mlog.Info("StopExecution")

	id := context.Params.ByName("id")
	// Get user information.
	user, mErr := controller.getUser(context)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		utils.ErrStopExecution.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	mErr = controller.executionProvider.StopExecution(user.AccountId, id)
	if mErr != nil {
		utils.ErrStopExecution.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.StopExecution.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
