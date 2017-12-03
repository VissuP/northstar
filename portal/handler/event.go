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

	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/portal/model"
)

// EventCallback executes a specific event action.
func (controller *Controller) EventCallback(context *gin.Context) {
	mlog.Debug("EventActions")

	// Get path parameter.
	connectionID := context.Params.ByName("connectionId")
	eventID := context.Params.ByName("eventId")
	eventType := context.Params.ByName("type")

	// Read message body.
	request := context.Request

	defer request.Body.Close()
	payload, err := ioutil.ReadAll(request.Body)

	// Populate our payload. Basically the only requirement is valid JSON.
	if err != nil {
		mlog.Error("Failed to execute event action with error: %s", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Create an event based on the URL parameters
	event, mErr := controller.portalProvider.ProcessEvent(eventID, eventType, payload)
	if mErr != nil {
		mlog.Error(mErr.String())
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Get the channel corresponding to our websocket so that it can be sent to the user.
	writer, err := controller.writers.Get(connectionID)
	if err != nil {
		mlog.Error("Get connection with id %s returned error: %s", connectionID, err.Error())
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	// Convert to expected type.
	channel, ok := writer.(chan model.Event)
	if !ok {
		mlog.Error("Connection writer type is invalid.")
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	// Write our event to the channel so it can be sent over the websocket.
	channel <- *event

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
