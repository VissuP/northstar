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

package api

import ()

const (
	// Defines the strings used to identify a resource as a ThingSpace resource.
	EVENT_KIND string = "ts.event"

	// Defines the Event resource schema version. Note that the format
	// of this value is MAJOR.MINOR where:
	//	MAJOR - Matches the version of the Dakota RESTful API
	//	MINOR - Represents the resource schema verison.
	EVENT_SCHEMA_VERSION string = "1.0"
)

const (
	// Defines the supported actions.
	GetAction     string = "get"
	SetAction     string = "set"
	UpdateAction  string = "update"
	CancelAction  string = "cancel"
	HistoryAction string = "history"
)

const (
	// Define the supported states.
	PendingState  string = "pending"
	UpdatedState  string = "update"
	FailedState   string = "failed"
	CanceledState string = "canceled"
)

// Defines the type used to represent Event resource. Note that
// ThingSpace will send an event resource with one additional
// property for each supported field where the name of the
// property matches the field id in the model.
type Event struct {
	// Thingspace Resource Identity
	Id          string `json:"id,omitempty"`
	Kind        string `json:"kind"`
	Version     string `json:"version"`
	CreatedOn   string `json:"createdon,omitempty"`
	LastUpdated string `json:"lastupdated,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// Thingspace Resource Foreign Identity

	ForeignId string   `json:"foreignid,omitempty"`
	Tags      []string `json:"tagids,omitempty"`

	// Thingspace Event Resource
	Action   string `json:"action,omitempty"`
	State    string `json:"state,omitempty"`
	DeviceId string `json:"deviceid,omitempty"`
	ModelId  string `json:"modelid,omitempty"`

	// Thingspace Event Resource Fields
	Fields map[string]interface{} `json:"fields,omitempty"`
}

func (event *Event) GetFieldId() string {
	// only one entry is expected in the event.Fields map
	for fieldId, _ := range event.Fields {
		return fieldId
	}
	return ""
}
