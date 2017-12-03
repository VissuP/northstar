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
	"errors"
	"fmt"

	"github.com/verizonlabs/northstar/northstarapi/config"
)

const (
	// Defines the string used to identify a resource as a Transformation resource.
	TransformationKind string = "ts.transformation"

	// Defines the Transformation resource schema version. Note that the format
	// of this value is MAJOR.MINOR where:
	//	MAJOR - Matches the version of the Service RESTful API
	//	MINOR - Represents the resource schema version.
	TransformationSchemaVersion string = "1.0"
)

// Defines the type used to identify the code type.
type CodeType string

const (
	// Defines the supported transformation codes.
	ObjectCodeType CodeType = "Object"
	SourceCodeType CodeType = "Source"
)

// Helper method used to get Code type from string.
func (codeType CodeType) ToString() string {
	return string(codeType)
}

// Helper method used to get scheme from code type.
func (codeType CodeType) GetScheme() string {
	if codeType == ObjectCodeType {
		return "s3"
	}

	return "base64"
}

// Defines the type used to represent a Transformation resource.
type Transformation struct {
	Kind        string    `json:"kind,omitempty"`
	Id          string    `json:"id,omitempty"`
	ExecutionId string    `json:"executionId, omitempty"`
	Version     string    `json:"version,omitempty"`
	CreatedOn   string    `json:"createdOn,omitempty"`
	LastUpdated string    `json:"lastUpdated,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Scheduled   bool      `json:"scheduled"`
	Schedule    *Schedule `json:"schedule"`
	Timeout     int       `json:"timeout,omitempty"`
	Arguments   Arguments `json:"arguments,omitempty"`
	EntryPoint  string    `json:"entryPoint"`
	Language    string    `json:"language"`
	Code        Code      `json:"code"`
	Memory      uint64    `json:"memory"`

	// Internal fields.
	SchedulerType string `json:"-"`
	SchedulerId   string `json:"-"`
}

// Defines the type used to represent a Transformation Code object.
type Code struct {
	Type  CodeType `json:"type"`
	Value string   `json:"value"`

	// Internal fields.
	Url string `json:"-"`
}

// Define internal type used for serialization.
type typeTransformation Transformation

// Helper method used to marshal transformation while translating fields.
func (transformation *Transformation) MarshalJSON() ([]byte, error) {
	var value typeTransformation

	value = typeTransformation(*transformation)

	// Make sure resource kind and version are set.
	value.Kind = TransformationKind
	value.Version = TransformationSchemaVersion

	return json.Marshal(&value)
}

//Validate validates the values of the transformation
func (transformation *Transformation) Validate() error {
	// Validate required fields are not empty.

	if transformation.Name == "" {
		return errors.New("The name is missing.")
	}

	if transformation.Language == "" {
		return fmt.Errorf("The language is missing or invalid.")
	}

	if transformation.EntryPoint == "" {
		return fmt.Errorf("The entry point is missing.")
	}

	if transformation.Memory > config.Configuration.MaxMemory {
		return fmt.Errorf("The requested memory %d is greater than the max of: %d", transformation.Memory, config.Configuration.MaxMemory)
	}

	if transformation.Memory == 0 {
		transformation.Memory = config.Configuration.DefaultMemory
	}

	// Validate the Code required fields.

	switch transformation.Code.Type {
	case ObjectCodeType:
	case SourceCodeType:
	default:
		return fmt.Errorf("The code type is missing or invalid.")
	}

	if transformation.Code.Value == "" {
		return fmt.Errorf("The code value is empty.")
	}
	return nil
}
