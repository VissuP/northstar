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
	"fmt"
	"time"
)

// Defines the supported event types.
const (
	NoneEventType   string = "None"
	TimerEventType  string = "Timer"
	DeviceEventType string = "Device"
)

type SnippetData struct {
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	CreatedOn   time.Time `json:"createdOn,omitempty"`
	UpdatedOn   time.Time `json:"updatedOn,omitempty"`
	Runtime     string    `json:"runtime,omitempty"`
	MainFn      string    `json:"mainfn,omitempty"`
	URL         string    `json:"url,omitempty"`
	Code        string    `json:"code,omitempty"`
	Timeout     int       `json:"timeout,omitempty"`
	Memory      uint64    `json:"memory,omitempty"`
	Callback    string    `json:"callback,omitempty"`
	Description string    `json:"description,omitempty"`
	EventType   string    `json:"eventType,omitempty"`
	EventId     string    `json:"eventId,omitempty"`
}

func (snippet *SnippetData) Validate() error {
	if snippet.Name == "" {
		return fmt.Errorf("Name is empty")
	}

	if snippet.Runtime == "" {
		return fmt.Errorf("Runtime type is empty")
	}

	if snippet.MainFn == "" {
		return fmt.Errorf("MainFN is empty")
	}

	if snippet.URL == "" {
		return fmt.Errorf("URL is empty")
	}

	if snippet.Timeout <= 0 {
		return fmt.Errorf("Timeout needs to be greater than zero")
	}

	switch snippet.EventType {
	case TimerEventType, DeviceEventType:
		if snippet.EventId == "" {
			return fmt.Errorf("EventId is empty.")
		}
	}

	return nil
}

func (snippet *SnippetData) Print() string {
	return fmt.Sprintf("ID: %s, "+
		"Name: %s, "+
		"CreatedOn: %s, "+
		"UpdatedOn: %s, "+
		"Runtime: %s, "+
		"MainFn: %s, "+
		"URL: %s, "+
		"Code: %s, "+
		"Timeout: %d, "+
		"Memory: %d, "+
		"Callback: %s, "+
		"Description: %s",
		snippet.Id,
		snippet.Name,
		snippet.CreatedOn,
		snippet.UpdatedOn,
		snippet.Runtime,
		snippet.MainFn,
		snippet.URL,
		snippet.Code,
		snippet.Timeout,
		snippet.Memory,
		snippet.Callback,
		snippet.Description)
}
