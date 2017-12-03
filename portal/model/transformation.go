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

const (
	// Defines the supported transformation codes.
	ObjectCodeType string = "Object"
	SourceCodeType string = "Source"
)

// Defines the type used to represent a Transformation resource.
type Transformation struct {
	Kind        string    `json:"kind,omitempty"`
	Id          string    `json:"id,omitempty"`
	Version     string    `json:"version,omitempty"`
	CreatedOn   string    `json:"createdOn,omitempty"`
	LastUpdated string    `json:"lastUpdated,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Scheduled   bool      `json:"scheduled"`
	Timeout     int       `json:"timeout,omitempty"`
	Arguments   Arguments `json:"arguments,omitempty"`
	EntryPoint  string    `json:"entryPoint"`
	Language    string    `json:"language"`
	Code        Code      `json:"code"`
	Schedule    *Schedule `json:"schedule,omitempty"`
}

// Defines the type used to represent code arguments.
type Arguments map[string]interface{}

// Defines the type used to represent a Transformation Code object.
type Code struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
