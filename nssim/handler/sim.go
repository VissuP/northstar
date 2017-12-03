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
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/config"
)

type systemInfo struct {
	Environment string `json:"environment"`
}

func (controller *Controller) GetInfo(context *gin.Context) {
	mlog.Debug("GetInfo")

	info := Summary{
		Environment: config.Configuration.Environment,
		Mode:        config.Configuration.Mode,
	}

	context.JSON(http.StatusOK, info)
}

//ExecutionEndpoint
func (controller *Controller) ReceiveCallback(context *gin.Context) {
	mlog.Debug("ReceiveCallback")
	id := context.Params.ByName("callbackID")
	writer, err := controller.callbacks.Get(id)
	if err != nil {
		mlog.Error("Error, callback %s not found. Error was: %s", id, err)
	}

	// Convert to expected type.
	callback, ok := writer.(chan json.RawMessage)
	if !ok {
		mlog.Error("Callback writer type is invalid.")
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	var data json.RawMessage
	if err := controller.Bind(context, &data); err != nil {
		mlog.Error("Failed to process callback with error: %+v", err)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	callback <- data
}
