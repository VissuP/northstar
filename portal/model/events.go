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

package model

import (
	"encoding/json"
	"fmt"
)

// EventType defines the type used to represent event type.
type EventType string

// Define supported event types.
const (
	EventTypeExecuteCell           EventType = "CellExecution"
	EventTypeExecuteTransformation EventType = "TransformationExecution"
	EventTypeExecuteResult         EventType = "ExecutionResult"
	EventTypeError                 EventType = "Error"
	EventTypePing                  EventType = "Ping"
)

// Helper method used to translate event types to string.
func (eventType EventType) ToString() string {
	return string(eventType)
}

// Event defines the type used to represent asynchronous events.
type Event struct {
	Type    EventType       `json:"type"`
	Id      string          `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

// NewErrorEvent returns a new error event.
func NewErrorEvent(template string, message ...interface{}) Event {
	formattedMsg := fmt.Sprintf(template, message)
	payload := make(map[string]interface{})
	payload["Value"] = formattedMsg

	rawMessage, _ := json.Marshal(&payload)

	return Event{
		Type:    EventTypeError,
		Payload: json.RawMessage(rawMessage),
	}
}
