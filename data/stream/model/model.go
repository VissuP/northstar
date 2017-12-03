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
)

type Source struct {
	Name       string      `json:"name,omitempty"`
	Connection interface{} `json:"connection,omitempty"`
}

type Function struct {
	Name       string      `json:"name,omitempty"`
	Parameters interface{} `json:"parameters,omitempty"`
	Evaluator  interface{} `json:"evaluator,omitempty"`
}

type JobData struct {
	Id           string     `json:"id,omitempty"`
	AccountId    string     `json:"accountId,omitempty"`
	InvocationId string     `json:"invocationId,omitempty"`
	Memory       uint64     `json:"memory,omitempty"`
	Source       Source     `json:"source,omitempty"`
	Functions    []Function `json:"functions,omitempty"`
	CreatedOn    time.Time  `json:"createdOn,omitempty"`
	UpdatedOn    time.Time  `json:"updatedOn,omitempty"`
	Status       string     `json:"status,omitempty"`
	ErrorDescr   string     `json:"errorDescr,omitempty"`
	Description  string     `json:"description:omitempty"`
}

func (j *JobData) ByteArrToSource(data []byte) error {
	var source Source
	if err := json.Unmarshal(data, &source); err != nil {
		return err
	}

	j.Source = source
	return nil
}

func (j *JobData) ByteArrToFunctions(data []byte) error {
	functions := make([]Function, 0)
	if err := json.Unmarshal(data, &functions); err != nil {
		return err
	}

	j.Functions = functions
	return nil
}

func (j *JobData) Validate() error {
	if j.InvocationId == "" {
		return fmt.Errorf("Invocation id is empty")
	}

	if len(j.Functions) < 1 {
		return fmt.Errorf("Number of functions less than one")
	}

	return nil
}

func (data *JobData) Print() string {
	return fmt.Sprintf("ID: %s, "+
		"AccountID: %s, "+
		"InvocationId: %s, "+
		"Memory: %v, "+
		"Source: %v, "+
		"Functions: %v, "+
		"CreatedOn: %s, "+
		"UpdatedOn: %s"+
		"Status: %s, "+
		"ErrorDescr: %s,"+
		"Description: %s",
		data.Id, data.AccountId, data.InvocationId, data.Memory, data.Source, data.Functions,
		data.CreatedOn, data.UpdatedOn, data.Status, data.ErrorDescr, data.Description)
}
