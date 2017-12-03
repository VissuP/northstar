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
	"encoding/json"
)

const (
	// Defines the dakota event field ids. Note that these need to match the
	// field ids in the Dakota Device Model registered with ThingSpace.
	BatteryFieldId           string = "battery"
	AirPlaneModeFieldId      string = "airplanemode"
	SleepStatusFieldId       string = "sleepstatus"
	LocationFieldId          string = "location"
	GpsEnabledFieldId        string = "gpsenabled"
	DeviceReadyFieldId       string = "deviceready"
	ReportingIntervalFieldId string = "reportinginterval"
	SounderFieldId           string = "sounder"
	ActivationDateFieldId    string = "activationdate"
	FirmwareVersionFieldId   string = "firmwareversion"

	StatusId                  string = "status"
	ConfigurableId            string = "configurable"
	MaxNumberPeripheralsId    string = "maxnumberperipherals"
	DataId                    string = "data"
	UmOemSpecificationFieldId string = "umoemspecification"

	NotSleeping int = 2
	Sleeping    int = 1
)

// Defines the value that represents Reporting  Interval field type.
type ReportingField struct {
	Value     int    `json:"value"`
	StartTime string `json:"starttime,omitempty"`
}

// Defines the value that represents AirPlane Mode field type.
type AirPlaneModeField struct {
	Value     uint16 `json:"value"`
	StartTime string `json:"starttime,omitempty"`
}

// Defines the type used to represent Status field.
type StatusField struct {
	Code   uint8  `json:"code"`
	Status string `json:"status"`
}

// Defines the type that represents Ultra Modem Data.
type UmDataField struct {
	Type                   string `json:"type"`
	SpecificationId        string `json:"specificationid"`
	SpecificationEventTime string `json:"specificationeventtime"`
	Value                  string `json:"value"`
}

// Defines the value that represents Sounder field type.
type SounderField struct {
	SoundProfile int `json:"soundprofile"`
	Duration     int `json:"duration"`
}

// Some background about following methods
//
// Notice that Fields in thingspace Event structure is defined as
// map[string]interface{}. When unmarshaling a JSON object into an
// interface{} type, JSON unmarshal creates a map[string] interface{}.
// Refer to https://golang.org/pkg/encoding/json/#Unmarshal
//
// So below are our convenience methods to convert these
// map[string] interface{} into specific structs.

func ToReportingField(t interface{}) (*ReportingField, error) {
	retObj := &ReportingField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}

func ToAirPlaneModeField(t interface{}) (*AirPlaneModeField, error) {
	retObj := &AirPlaneModeField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}

func ToSleepStatus(t interface{}) (int, error) {
	var i int
	if err := convert(t, &i); err != nil {
		return 0, err
	}
	return i, nil
}

// json unmarshal converts all JSON numbers to float64. But we need
// an int. So using same technique to do this conversion.
func ToBatteryLevel(t interface{}) (int, error) {
	var i int
	if err := convert(t, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func ToUmDataField(t interface{}) (*UmDataField, error) {
	retObj := &UmDataField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}

func ToSounderField(t interface{}) (*SounderField, error) {
	retObj := &SounderField{}
	if err := convert(t, retObj); err != nil {
		return nil, err
	}
	return retObj, nil
}

func ToGpsEnabled(t interface{}) (bool, error) {
	var b bool
	if err := convert(t, &b); err != nil {
		return false, err
	}
	return b, nil
}

func ToDeviceReady(t interface{}) (bool, error) {
	var b bool
	if err := convert(t, &b); err != nil {
		return false, err
	}
	return b, nil
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
