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

type InvocationData struct {
	Id          string                 `json:"id,omitempty"`
	RTEId       string                 `json:"rteId,omitempty"`
	SnippetId   string                 `json:"snippetId,omitempty"`
	Partition   int                    `json:"partition,omitempty"`
	CreatedOn   time.Time              `json:"createdOn,omitempty"`
	StartedOn   time.Time              `json:"startedOn,omitempty"`
	FinishedOn  time.Time              `json:"finishedOn,omitempty"`
	UpdatedOn   time.Time              `json:"updatedOn,omitempty"`
	ElapsedTime float64                `json:"elapsedTime,omitempty"`
	Runtime     string                 `json:"runtime,omitempty"`
	MainFn      string                 `json:"mainfn,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Code        string                 `json:"code,omitempty"`
	Timeout     int                    `json:"timeout,omitempty"`
	Memory      uint64                 `json:"memory,omitempty"`
	Callback    string                 `json:"callback,omitempty"`
	Args        map[string]interface{} `json:"args,omitempty"`
	Stdout      string                 `json:"stdout,omitempty"`
	Result      string                 `json:"result,omitempty"`
	Status      string                 `json:"status,omitempty"`
	ErrorDescr  string                 `json:"errorDescr,omitempty"`
}

func (invocation *InvocationData) ArgsToByteArr() ([]byte, error) {
	bytes, err := json.Marshal(invocation.Args)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (invocation *InvocationData) ByteArrToArgs(data []byte) error {
	args := make(map[string]interface{})
	if err := json.Unmarshal(data, &args); err != nil {
		return err
	}

	invocation.Args = args
	return nil
}

func (invoke *InvocationData) ValidateOnAdd() error {
	if invoke.SnippetId == "" {
		return fmt.Errorf("Snippet id is empty")
	}

	if invoke.Status == "" {
		return fmt.Errorf("Status is empty")
	}

	return nil
}

func (invoke *InvocationData) Print() string {
	return fmt.Sprintf("ID: %s\n"+
		"SnippetID: %s\n"+
		"Partition: %v\n"+
		"CreatedOn: %s\n"+
		"StartedOn: %s\n"+
		"FinishedOn: %s\n"+
		"UpdatedOn: %s\n"+
		"ElapsedTime: %f\n"+
		"Runtime: %s\n"+
		"MainFn: %s\n"+
		"URL: %s\n"+
		"Code: %s\n"+
		"Timeout: %d\n"+
		"Memory: %d\n"+
		"Callback: %s\n"+
		"Args: %s\n"+
		"Stdout: %s\n"+
		"Result: %s\n"+
		"Status: %s\n"+
		"ErrorDescr: %s",
		invoke.Id,
		invoke.SnippetId,
		invoke.Partition,
		invoke.CreatedOn,
		invoke.StartedOn,
		invoke.FinishedOn,
		invoke.UpdatedOn,
		invoke.ElapsedTime,
		invoke.Runtime,
		invoke.MainFn,
		invoke.URL,
		invoke.Code,
		invoke.Timeout,
		invoke.Memory,
		invoke.Callback,
		invoke.Args,
		invoke.Stdout,
		invoke.Result,
		invoke.Status,
		invoke.ErrorDescr)
}
