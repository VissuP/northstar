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

// Configuration defines the type used to represent configuration object.
type Configuration struct {
	Runtimes []RuntimeConfiguration `json:"runtimes"`
}

// RuntimeConfiguration defines the type used to describe supported runtime configurations.
type RuntimeConfiguration struct {
	Language      string         `json:"language"`
	Version       string         `json:"version"`
	InstanceTypes []InstanceType `json:"instanceTypes"`
}

// InstanceType defines the type used to represent supported runtime instance types.
type InstanceType struct {
	Model        string `json:"model"`
	Instances    int    `json:"instances,omitempty"`
	MaxCores     int    `json:"maxCores,omitempty"`
	MasterCores  int    `json:"masterCores,omitempty"`
	MasterMemory int    `json:"masterMemory,omitempty"`
	WorkerCores  int    `json:"workerCores,omitempty"`
	WorkerMemory int    `json:"workerMemory,omitempty"`
	Partitions   int    `json:"partitions,omitempty"`
	Replication  int    `json:"replication,omitempty"`
}
