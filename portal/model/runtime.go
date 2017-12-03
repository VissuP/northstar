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

// Runtime defines the type used to represent runtime resource.
type Runtime struct {
	Model          string     `json:"model,omitempty"`
	Language       string     `json:"language,omitempty"`
	ImageTag       string     `json:"imageTag,omitempty"`
	Concurrency    int        `json:"concurrency,omitempty"`
	MaxConcurrency int        `json:"maxConcurrency,omitempty"`
	Memory         int        `json:"memory,omitempty"`
	MasterMemory   int        `json:"masterMemory,omitempty"`
	WorkerMemory   int        `json:"workerMemory,omitempty"`
	Partitions     int        `json:"partitions,omitempty"`
	Replication    int        `json:"replication,omitempty"`
	Instances      []Instance `json:"instances,omitempty"`
}

// Instance defines the type used to represent instance resource.
type Instance struct {
	Id            string `json:"id"`
	CreatedOn     string `json:"createdOn,omitempty"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
}
