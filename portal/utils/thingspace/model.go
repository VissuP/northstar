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

package thingspace

import (
	"encoding/json"
	"fmt"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"time"
)

// Model defines the type used to represent ThingSpace Model Resource.
type Model struct {
	Id           string    `json:"id,omitempty"`
	Version      string    `json:"version,omitempty"`
	ModelKind    string    `json:"modelKind,omitempty"`
	ModelVersion string    `json:"modelVersion,omitempty"`
	CreatedOn    time.Time `json:"createdon,omitempty"`
	LastUpdated  time.Time `json:"lastupdated,omitempty"`
	Fields       Fields    `json:"fields,omitempty"`
}

// Fields defines the type that represents the dictionary of device fields.
type Fields map[string]Field

// Field defines the type that represents a device field.
type Field struct {
	Type     string `json:"type"`
	Semantic string `json:"semantic"`
	Source   bool   `json:"source"`
	Sink     bool   `json:"sink"`
}

// GetModels returns the collection of device models registered with ThingSpace.
func (userClient *UserClient) GetModels(accessToken string) ([]Model, *management.Error) {
	mlog.Debug("GetModels")

	// Generate request headers.
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	// Get response body.
	body, mErr := management.GetWithHeaders(userClient.hostAndPort, "/api/v2/models", headers)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Failed to get model with error: %v", mErr))
	}

	// Unmarshal model object.
	var models []Model

	if err := json.Unmarshal(body, &models); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to unmarshal model with error: %v", err))
	}

	return models, nil
}
