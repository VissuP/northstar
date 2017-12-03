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

// Note that changes to this page should update Dakota Provider API documentation
// for Device Resource. E.g., http://confluence.verizon.com/display/KSIPROD/Provider+Service+Client+API

const (
	// Defines the strings used to identify a resource as a Dakota Device
	// resource.
	DEVICE_ULTRATAG   string = "ts.device.ultratag"
	DEVICE_ULTRAMODEM string = "ts.device.ultramodem"

	// Defines the Device resource schema version. Note that the format
	// of this value is MAJOR.MINOR where:
	//	MAJOR - Matches the version of the Dakota RESTful API
	//	MINOR - Represents the resource schema verison.
	DEVICE_SCHEMA_VERSION string = "1.0"
)

const (
	// Defines the state of a device that is register to an account but not connected to the LTE network.
	DEVICE_STATE_REGISTER string = "Register"

	// Defines the state of a device that is register to an account and connected to the LTE network.
	DEVICE_STATE_READY string = "Ready"

	// Defines the state of a device that has been deactivated.
	DEVICE_STATE_DEACTIVATED string = "Deactivated"
)

// Defines the type used to represent Device resource.
type Device struct {
	// Thingspace Resource Identity
	Id          string `json:"id,omitempty"`
	Kind        string `json:"kind"`
	Version     string `json:"version"`
	CreatedOn   string `json:"createdon,omitempty"`
	LastUpdated string `json:"lastupdated,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	// Thingspace Resource Foreign Identity
	ForeignId string   `json:"foreignid,omitempty"`
	Tags      []string `json:"tagids,omitempty"`

	// Thingspace Device Resource. Note that this might not
	// be the complete property list. Only the onces needed
	// by dakota.
	ProviderId string `json:"providerid,omitempty"`
	State      string `json:"state,omitempty"`

	QrCode       string `json:"qrcode" binding:"required"`
	SerialNumber string `json:"serial,omitempty"`
	Imei         uint64 `json:"imei,omitempty"`
	Imsi         uint64 `json:"imsi,omitempty"`
	IccId        string `json:"iccid,omitempty"`
	BluetoothMac string `json:"bluetoothmac,omitempty"`

	Family                 string    `json:"family,omitempty"`
	Actions                string    `json:"actions,omitempty"`
	Model                  string    `json:"model,omitempty"`
	ProductModel           string    `json:"productmodel,omitempty"`
	HardwareVersion        string    `json:"hardwareversion,omitempty"`
	Firmware               *Firmware `json:"firmware"`
	GpsEnabled             int16     `json:"gpsenabled"`
	Configurable           bool      `json:"configurable"`
	MaxNumberOfPeripherals int       `json:"maxnumberofperipherals,omitempty"`
	Lifetime               int       `json:"lifetime,omitempty"`
	InitBatteryLevel       int       `json:"initbatterylevel,omitempty"`
	LowBatteryLevel        int       `json:"lowbatterylevel,omitempty"`
	CriticalBatteryLevel   int       `json:"criticalbatterylevel,omitempty"`
	BatteryLevelTestStatus string    `json:"batterylevelteststatus,omitempty"`
	ProvisionedOn          string    `json:"provisionedon,omitempty"`
	RegisteredOn           string    `json:"registeredon,omitempty"`
	ActivatedOn            string    `json:"activatedon,omitempty"`
	DeactivatedOn          string    `json:"deactivatedon,omitempty"`

	// Thingspace Device Resource Fields
	Fields map[string]interface{} `json:"fields,omitempty"`
}

// Defines the type used to represent device firmware information.
type Firmware struct {
	Version   string `json:"version,omitempty"`
	UpdatedOn string `json:"updatedOn,omitempty"`
}
