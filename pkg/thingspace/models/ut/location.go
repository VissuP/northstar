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

package ut

import (
	"time"
)

const (
	// Defines the string used to identify a resource as a Location resource.
	LOCATION_KIND string = "ts.location"

	// Defines the Location resource schema version. Note that the format
	// of this value is MAJOR.MINOR where:
	//	MAJOR - Matches the version of the Dakota RESTful API
	//	MINOR - Represents the resource schema verison.
	LOCATION_SCHEMA_VERSION string = "1.1"
)

type LocationType string

const (
	// Defines the different location types (e.g., from GPS data,
	// estimated from Cell Id, etc.)
	LOCATION_TYPE_CELL_ID LocationType = "CellId"
	LOCATION_TYPE_GPS                  = "Gps"
)

// Defines the type used to represent UserApps Location resource.
type LocationField struct {
	Kind          string       `json:"kind,omitempty"`
	Id            string       `json:"id,omitempty"`
	Version       string       `json:"version,omitempty"`
	CreatedAt     string       `json:"createdAt,omitempty"`
	Address       string       `json:"address,omitempty"`
	Altitude      uint32       `json:"altitude,omitempty"`
	Longitude     float64      `json:"longitude" binding:"required"`
	Latitude      float64      `json:"latitude" binding:"required"`
	Shape         string       `json:"shape,omitempty"`
	Radius        uint32       `json:"radius,omitempty"`
	SemiMajorAxis uint32       `json:"semimajoraxis,omitempty"`
	SemiMinorAxis uint32       `json:"semiminoraxis,omitempty"`
	Place         Place        `json:"place,omitempty"`
	Type          LocationType `json:"type,omitempty"`
	EventTime     time.Time    `json:"eventtime,omitempty"`
}

type Place struct {
	Id        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Radius    float64 `json:"radius,omitempty"`
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

// Some background about following method
//
// Notice that Fields in thingspace Event structure is defined as
// map[string]interface{}. When unmarshaling a JSON object into an
// interface{} type, JSON unmarshal creates a map[string] interface{}.
// Refer to https://golang.org/pkg/encoding/json/#Unmarshal
//
// So below are our convenience methods to convert these
// map[string] interface{} into specific structs.

func ToLocationField(t interface{}) (*LocationField, error) {
	retObj := &LocationField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}
