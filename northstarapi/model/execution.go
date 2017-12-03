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

//ExecutionResponse contains the output metadata of an execution
type ExecutionResponse struct {
	AccountID        string        `json:"accountId,omitempty"`
	Status           string        `json:"status,omitempty"`
	ErrorDescription string        `json:"errorDescription,omitempty"`
	InvocationID     string        `json:"invocationId,omitempty"`
	Callback         string        `json:"callback,omitempty"`
	SnippetID        string        `json:"snippetID,omitempty"`
	RteID            string        `json:"rteId,omitempty"`
	StartedOn        time.Time     `json:"startedOn,omitempty"`
	FinishedOn       time.Time     `json:"finishedOn,omitempty"`
	ElapsedTime      time.Duration `json:"elapsedTime,omitempty"`
}

type ExecutionRequest struct {
	AccountId   string                 `json:"-"`
	ExecutionId string                 `json:"executionId"`
	Name        string                 `json:"name"`
	Language    string                 `json:"language"`
	EntryPoint  string                 `json:"entryPoint"`
	Arguments   map[string]interface{} `json:"arguments"`
	Code        string                 `json:"code"`
	Timeout     int                    `json:"timeout"`
	Callback    string                 `json:"callback"`
	Memory      uint64                 `json:"memory,omitempty"`
}
