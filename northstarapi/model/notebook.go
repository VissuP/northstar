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
	"time"

	"github.com/verizonlabs/northstar/northstarapi/config"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

const (
	// Defines the string used to identify a resource as a Notebook resource.
	NotebookKind string = "ts.notebook"

	// Defines the Notebook resource schema version. Note that the format
	// of this value is MAJOR.MINOR where:
	//	MAJOR - Matches the version of the Service RESTful API
	//	MINOR - Represents the resource schema version.
	NotebookSchemaVersion string = "1.0"
)

// Defines the type used to represent a notebook resource.
// Note that permissions represents the permissions of the
// current user. This property is ready only.
type Notebook struct {
	Kind        string `json:"kind,omitempty"`
	Id          string `json:"id,omitempty"`
	Etag        string `json:"etag,omitempty"`
	Version     string `json:"version,omitempty"`
	CreatedOn   string `json:"createdOn,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
	Name        string `json:"name,omitempty"`
	Permissions string `json:"permissions,omitempty"`
	Cells       []Cell `json:"cells,omitempty"`
}

// Defines internal event type for marshaling/unmarshaling.
type typeNotebook Notebook

// Helper method used to marshal runtime while translating fields.
func (notebook *Notebook) MarshalJSON() ([]byte, error) {
	var value typeNotebook

	value = typeNotebook(*notebook)

	// Make sure resource kind and version are set.
	value.Kind = NotebookKind
	value.Version = NotebookSchemaVersion

	return json.Marshal(&value)
}

// Validate notebook
func (notebook *Notebook) Validate() error {
	for _, cell := range notebook.Cells {
		if err := cell.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Defines the type used to represent a notebook cell.
type Cell struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Input       Input    `json:"input"`
	Output      Output   `json:"output"`
	Settings    Settings `json:"settings"`
	ExecutionId string   `json:"executionId"`
}

// Validate cell
func (cell *Cell) Validate() error {
	input := &cell.Input
	if input.Timeout <= 0 {
		mlog.Debug("No timeout provided, falling back to default")
		input.Timeout = config.MaxTimeout
	}

	if input.Timeout > config.MaxTimeout {
		return fmt.Errorf("The timeout value is out of range (%d of %d).",
			input.Timeout, config.MaxTimeout)
	}

	// Validate the body is not empty.
	if input.Body == "" {
		return fmt.Errorf("The input body is missing.")
	}

	//Validate that the code isn't too large
	if len(input.Body) > config.MaxCodeSize {
		return fmt.Errorf("The input body is too large.")
	}

	//Validate arg count
	if len(input.Arguments) > config.MaxArgCount {
		return fmt.Errorf("Number of arguments is too high (%d of %d).",
			len(input.Arguments), config.MaxArgCount)
	}

	// Validate the cell name is not empty.
	if cell.Name == "" {
		return fmt.Errorf("The cell name is empty.")
	}

	//Validate the settings for the cell
	if err := cell.Settings.Validate(); err != nil {
		return err
	}

	return nil
}

// Settings defines the user display settings of a notebook
type Settings struct {
	Status            string        `json:"status"`
	ShowCode          bool          `json:"showCode"`
	ShowOutput        bool          `json:"showOutput"`
	ShowConfiguration bool          `json:"showConfiguration"`
	Visualization     Visualization `json:"visualization, omitempty"`
	Memory            uint64        `json:"memory,omitempty"`
}

// Validate validates the settings
func (settings *Settings) Validate() error {
	if settings.Memory > config.Configuration.MaxMemory {
		return fmt.Errorf("Requested memory (%d) is greater than the max of: %d.", settings.Memory, config.Configuration.MaxMemory)
	}

	return nil
}

// Defines internal type for marshaling/unmarshaling.
type typeSettings Settings

// UnmarshalJSON is a helper method used to unmarshal runtime while validating required fields.
func (settings *Settings) UnmarshalJSON(data []byte) error {
	var value typeSettings

	// Unmarshal to the internal type
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	//If the memory is unset, set it to the default.
	if value.Memory == 0 {
		value.Memory = config.Configuration.DefaultMemory
	}

	*settings = Settings(value)

	return nil
}

// Visualization defines the visualization specific settings.
type Visualization struct {
	Parameters map[string]interface{} `json:"parameters"`
}

// Defines the type used to represent an input type.
type CellType string

const (
	StaticCellType CellType = "Static"
	CodeCellType   CellType = "Code"
	QueryCellType  CellType = "Query"
)

// Returns string representation.
func (cellType CellType) ToString() string {
	return string(cellType)
}

// Defines the type used to represent arguments to input cell.
type Arguments map[string]interface{}

// Defines the type used to represent notebook input.
type Input struct {
	Type       CellType  `json:"type"`
	Language   string    `json:"language"`
	Arguments  Arguments `json:"arguments,omitempty"`
	EntryPoint string    `json:"entryPoint"`
	Body       string    `json:"body"`
	Timeout    int       `json:"timeout"`
}

// Defines the Output status code.
type OutputStatusCode string

const (
	// Define the supported output status codes.
	OutputUnknownStatus       OutputStatusCode = "Unknown"
	OutputSuccessStatus       OutputStatusCode = "Success"
	OutputFailedStatus        OutputStatusCode = "Failed"
	OutputTimeoutStatus       OutputStatusCode = "ExecutionTimeout"
	OutputInternalErrorStatus OutputStatusCode = "InternalError"
	OutputOutOfMemoryStatus   OutputStatusCode = "OutOfMemory"
	OutputRunningStatus       OutputStatusCode = "Running"
	CodeGetFailedStatus       OutputStatusCode = "Unable to read code."
)

// Helper method used to get output status code type as string.
func (outputStatusCode OutputStatusCode) ToString() string {
	return string(outputStatusCode)
}

// Helper method used to translate strings to status code.
func (outputStatusCode OutputStatusCode) FromString(status string) OutputStatusCode {
	switch status {
	case OutputOutOfMemoryStatus.ToString():
		return OutputOutOfMemoryStatus
	case OutputSuccessStatus.ToString():
		return OutputSuccessStatus
	case OutputFailedStatus.ToString():
		return OutputFailedStatus
	case OutputTimeoutStatus.ToString():
		return OutputTimeoutStatus
	case OutputInternalErrorStatus.ToString():
		return OutputInternalErrorStatus
	case OutputRunningStatus.ToString():
		return OutputRunningStatus
	}

	return OutputUnknownStatus
}

// Defines the default output status descriptions.
// Defines the default status descriptions.
var DefaultOutputStatusDescriptions = map[OutputStatusCode]string{
	OutputSuccessStatus:       "The execution has succeeded.",
	OutputFailedStatus:        "The execution has failed.",
	OutputTimeoutStatus:       "The execution did not completed within the expected time.",
	OutputInternalErrorStatus: "The execution failed due to an unexpected condition.",
	OutputUnknownStatus:       "The execution failed due to an unknown condition.",
	OutputOutOfMemoryStatus:   "The execution failed due to insufficient memory.",
	OutputRunningStatus:       "The execution is currently running.",
}

// Defines the type used to represent notebook output.
type Output struct {
	Id                string           `json:"id,omitempty"`
	Status            OutputStatusCode `json:"status,omitempty"`
	StatusDescription string           `json:"statusDescription,omitempty"`
	ExecutionOutput   string           `json:"executionOutput,omitempty"`
	ExecutionResults  *CellResults     `json:"executionResults,omitempty"`
	ElapsedTime       int              `json:"elapsedTime,omitempty"`
	LastExecution     time.Time        `json:"lastExecution,omitempty"`
}

// Define the type used to represent result types.
type ResultTypes string

const (
	// Defines supported result types.
	JsonResultType  ResultTypes = "application/json"
	JpegResultType  ResultTypes = "image/jpeg"
	TextResultType  ResultTypes = "text/plain"
	HtmlResultType  ResultTypes = "text/html"
	TableResultType ResultTypes = "application/vnd.vz.table"
)

// Returns string representation.
func (resultTypes ResultTypes) ToString() string {
	return string(resultTypes)
}

// Defines the type used to represent notebook results.
type CellResults struct {
	Type    ResultTypes `json:"type"`
	Content interface{} `json:"content"`
}
