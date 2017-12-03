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

package ott

import ()
import "encoding/json"

const (
	// Defines the OTT event field ids. Note that these need to match the
	// field ids in the OTT Device Model registered with ThingSpace.

	OttStatusFieldId        string = "status"
	OttPollingIntervalField string = "pollinginterval"
	OttSleepField                  = "sleep"
	OttUpdateField                 = "update"

	// Define other constants

	// 512 - 1 command-byte - 2 bytes encoding overhead for size
	OttUpdateFieldMaxBinaryLength = 509
)

const (
	DEVICE_OTT string = "ts.device.ott"
)

// Defines the value that represents Polling Interval field type.
type PollingIntervalField struct {
	Value     int    `json:"value"`
	StartTime string `json:"starttime,omitempty"`
}

// Defines the value that represents Sleep field type.
type SleepField struct {
	Value     int    `json:"value"`
	StartTime string `json:"starttime,omitempty"`
}

// Defines the value that represents Update field type
type UpdateField string

// Some background about following methods
//
// Notice that Fields in thingspace Event structure is defined as
// map[string]interface{}. When unmarshaling a JSON object into an
// interface{} type, JSON unmarshal creates a map[string] interface{}.
// Refer to https://golang.org/pkg/encoding/json/#Unmarshal
//
// So below are our convenience methods to convert these
// map[string] interface{} into specific structs.

func ToPollingIntervalField(t interface{}) (*PollingIntervalField, error) {
	retObj := &PollingIntervalField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}

func ToSleepField(t interface{}) (*SleepField, error) {
	retObj := &SleepField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}

func convert(from interface{}, to interface{}) error {
	body, err := json.Marshal(from)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, to); err != nil {
		return err
	}
	return nil
}
