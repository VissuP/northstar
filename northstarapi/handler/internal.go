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
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

//ExecutionCallback processes the results of an execution
func (controller *Controller) ExecutionCallback(context *gin.Context) {
	mlog.Info("ExecutionCallback")

	var response model.ExecutionResponse
	err := controller.Bind(context, &response)
	if err != nil {
		mlog.Error("Failed to parse RTE event with error: %v.", err)
		utils.ErrExecutionCallback.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	controller.executionProvider.ExecutionCallback(&response)
	utils.ExecutionCallback.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
