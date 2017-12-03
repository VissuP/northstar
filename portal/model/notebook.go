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

import "time"

// Notebook defines the overall resource for a notebook
type Notebook struct {
	Version     string `json:"version"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Cells       []Cell `json:"cells"`
	Permissions string `json:"permissions"`
}

// Cell defines the code and output for a snippet of code in a notebook.
type Cell struct {
	ID       string                 `json:"id,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Language string                 `json:"language,omitempty"`
	Code     string                 `json:"code,omitempty"`
	Output   Output                 `json:"output,omitempty"`
	Options  Options                `json:"options"`
	Inputs   map[string]interface{} `json:"inputs,omitempty"`
}

// Options defines the user display settings of a notebook
type Options struct {
	Status            string        `json:"status"`
	ShowCode          bool          `json:"showCode"`
	ShowOutput        bool          `json:"showOutput"`
	ShowConfiguration bool          `json:"showConfiguration"`
	Visualization     Visualization `json:"visualization,omitempty"`
	MainFunction      string        `json:"mainFunction,omitempty"`
	Timeout           int           `json:"timeout,omitempty"`
}

//Visualization defines the component-specific settings for a visualization (i.e. d3 settings)
type Visualization struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters"`
}

// Output defines the execution output structure for a cell.
type Output struct {
	State         string       `json:"state,omitempty"`
	Stdout        string       `json:"stdout,omitempty"`
	Stderr        string       `json:"stderr,omitempty"`
	Results       *CellResults `json:"results, omitempty"`
	ElapsedTime   int          `json:"elapsedTime,omitempty"`
	LastExecution time.Time    `json:"lastExecution,omitempty"`
}

// CellResults defines the visualization output structure for an execution
type CellResults struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
