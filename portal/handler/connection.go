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
	"net/http"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pborman/uuid"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/portal/config"
	"github.com/verizonlabs/northstar/portal/model"
)

// GetConnection sets up the websocket (creates a thread-safe map to store data as well as the websocket reads from
// the websocket until the socket is closed. Data that is read from the websocket is then sent off to
// the provider for processing of the event.
func (controller *Controller) GetConnection(context *gin.Context) {
	mlog.Debug("GetConnection")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get request token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  config.Configuration.ConnectionBufferSize,
		WriteBufferSize: config.Configuration.ConnectionBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := wsUpgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		mlog.Error("Failed to upgrade websocket: %+v", err)
		controller.RenderServiceError(context, management.GetInternalError("Failed to upgrade websocket"))
		return
	}

	connectionId := uuid.NewRandom().String()
	mlog.Info("Generated socket Id: %s", connectionId)
	writeChannel := make(chan model.Event)
	controller.writers.Set(connectionId, writeChannel)

	// Launch the go routing that hanndles writes to the websocket.
	go connectionWriter(conn, writeChannel)

	// Note that this look will process incoming events from the
	// client side until the socket is close.
	for {
		var event model.Event

		// Read and unmarshal JSON off of websocket. In case of error,
		// close the socket and remove from the list of connections.
		if err := conn.ReadJSON(&event); err != nil {
			mlog.Error("Error reading event from socket: %s", err.Error())
			close(writeChannel)
			controller.writers.Delete(connectionId)
			break
		}

		// If we don't have an ID, we won't have any way to connect it back to the generator of the event.
		if event.Id == "" {
			mlog.Error("Error, event ID not set.")
			writeChannel <- model.ErrorEventIdMissing
			continue
		}

		// We need the event type to figure out who processes the event. Is it an execution? Runtime adjustment? Something else?
		if event.Type == "" {
			mlog.Error("Error, event type not set.")
			writeChannel <- model.ErrorEventTypeMissing
			continue
		}

		//EventTypePing is an event type that can be used for keep-alive of connections. Requests of this type only serve to generate activity and keep the connection alive.
		if event.Type == model.EventTypePing {
			continue
		}

		// Process the event. In case of error, send error proper error event to client.
		// Note that results (if any) of the processing are expected to be return asynchronously.
		if mErr := controller.processEvent(token.AccessToken, connectionId, event); mErr != nil {
			mlog.Error("Process event failed with error: %+v", mErr)
			writeChannel <- model.NewErrorEvent("%s", mErr.Description)
		}
	}

}

// connectionWriter runs as a go routine who's entire job is to sit waiting to receive events on the write channel.
// When it receives an event, it writes it to the websocket. Note that reading from a channel is blocking */
func connectionWriter(connection *websocket.Conn, writeChannel chan model.Event) {
	mlog.Debug("connectionWriter")

	for {
		// Wait for data on the write channel. If the channel is closed, ok will be false and we will
		// automatically stop waiting. The reader will handle cleanup.
		event, ok := <-writeChannel
		if !ok {
			mlog.Info("Write channel closed. Returning")
			return
		}

		if err := connection.WriteJSON(event); err != nil {
			mlog.Error("Connection write JSON failed with error: %v", err)
			return
		}
	}
}

// processEvent is a helper method used to process events received from the client.
func (controller *Controller) processEvent(token string, connectionId string, event model.Event) *management.Error {
	mlog.Debug("ProcessEvent: event: %+v", event)

	switch event.Type {
	case model.EventTypeExecuteCell:
		{
			// Generate the callback url.
			callbackUrl := controller.getEventCallbackURL(connectionId, event.Id, model.EventTypeExecuteResult)

			// Parse expected payload.
			var cell model.Cell

			if err := json.Unmarshal([]byte(event.Payload), &cell); err != nil {
				mlog.Error("Event umarshal returned error: %+v", err)
				return management.GetBadRequestError("The event payload type is invalid. Expecting cell type.")
			}

			// Execute cell.
			if mErr := controller.portalProvider.ExecuteCell(token, callbackUrl, &cell); mErr != nil {
				return mErr
			}
		}
	case model.EventTypeExecuteTransformation:
		{
			// Generate the callback url.
			callbackUrl := controller.getEventCallbackURL(connectionId, event.Id, model.EventTypeExecuteResult)

			// Parse expected payload.
			var transformation model.Transformation

			if err := json.Unmarshal([]byte(event.Payload), &transformation); err != nil {
				mlog.Error("Event umarshal returned error: %+v", err)
				return management.GetBadRequestError("The event payload type is invalid. Expecting transformation type.")
			}

			// Execute transformation.
			if mErr := controller.portalProvider.ExecuteTransformation(token, callbackUrl, &transformation); mErr != nil {
				return mErr
			}
		}
	default:
		return management.GetBadRequestError("The event type is not supported.")
	}

	return nil
}

// getEventcallbackURL is a helper method used to generate callback url.
func (controller *Controller) getEventCallbackURL(connectionID, eventID string, eventType model.EventType) string {
	return fmt.Sprintf("http://%s/%s/%s/connections/%s/events/%s/callbacks/%s",
		config.Configuration.ServiceHostPort, model.InternalContext, model.Version, connectionID, eventID, eventType.ToString())
}
