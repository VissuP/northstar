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

package api

const (
	// resource kind identifier
	PROVIDER_KIND = "ts.provider"

	// resource schema version. note that the format
	// of this value is MAJOR.MINOR where:
	//	MAJOR - Matches the version of the RESTful API
	//	MINOR - Represents the resource schema verison.
	PROVIDER_VERSION = "1.0"
)

// Defines the type used to represent Provider resource.
type Provider struct {
	// Thingspace Resource Identity
	Id                string `json:"id,omitempty"`
	Kind              string `json:"kind,omitempty"`
	Version           string `json:"version,omitempty"`
	VersionId         string `json:"versionid,omitempty"`
	CreatedOn         string `json:"createdon,omitempty"`
	LastUpdated       string `json:"lastupdated,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	Historydelegation bool   `json:"historydelegation,omitempty"`
	// Thingspace Provider Resource
	Inventory InventoryType `json:"inventory"`
	Source    SourceType    `json:"source"`
	Sink      SinkType      `json:"sink"`
}

// Defines the Inventory Data Provider information.
type InventoryType struct {
	Alias        string `json:"alias,omitempty"`
	BasePath     string `json:"basePath,omitempty"`
	RelativePath string `json:"relativePath,omitempty"`
	HostAndPort  string `json:"hostAndPort,omitempty"`
}

// Defines the Source Channel Provider information.
type SourceType struct {
}

// Defines the Sink Channel Provider information.
type SinkType struct {
	Alias        string `json:"alias,omitempty"`
	BasePath     string `json:"basePath,omitempty"`
	RelativePath string `json:"relativePath,omitempty"`
	HostAndPort  string `json:"hostAndPort,omitempty"`
}
