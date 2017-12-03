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
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

func (controller *Controller) ListStreams(context *gin.Context) {
	mlog.Info("ListStreams")

	// Get account id.
	accountId, mErr := controller.getAccountId(context)
	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrListStreams.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	jobs, mErr := controller.streamProvider.ListStreams(accountId)
	if mErr != nil {
		controller.RenderServiceError(context, mErr)
		utils.ErrListStreams.Incr()
		return
	}

	utils.ListStreams.Incr()
	context.JSON(http.StatusOK, jobs)
}

func (controller *Controller) GetStream(context *gin.Context) {
	mlog.Info("GetStream")
	jobId := context.Params.ByName("id")

	// Get account id.
	accountId, mErr := controller.getAccountId(context)
	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrGetStream.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	job, mErr := controller.streamProvider.GetStream(accountId, jobId)
	if mErr != nil {
		mlog.Error("Failed to get job with error: %v", mErr)
		utils.ErrGetStream.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.GetStream.Incr()
	context.JSON(http.StatusOK, job)
}

func (controller *Controller) RemoveStream(context *gin.Context) {
	mlog.Info("RemoveStream")
	jobId := context.Params.ByName("id")

	// Get account id.
	accountId, mErr := controller.getAccountId(context)
	if mErr != nil {
		mlog.Error("Failed to get account id with error: %v", mErr)
		utils.ErrGetStream.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	mErr = controller.streamProvider.RemoveStream(accountId, jobId)
	if mErr != nil {
		mlog.Error("Failed to get remove job with error: %v", mErr)
		utils.ErrRemoveStream.Incr()
		controller.RenderServiceError(context, mErr)
		return
	}

	utils.RemoveStream.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
