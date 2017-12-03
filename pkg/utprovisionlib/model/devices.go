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
	"fmt"
	"github.com/verizonlabs/northstar/pkg/management"
	"time"
)

const (
	// Defines supported hardware types.
	IdTypeMacAddress = "mac"
	IdTypeImsi       = "imsi"
	IdTypeImei       = "imei"
	IdTypeUuid       = "uuid"
)

const (
	// Defines supported credential types.
	CredentialsTypeNone    = "none"
	CredentialsTypeSecret  = "secret"
	CredentialsTypeCert    = "oem-certificate"
	CredentialsTypeDevCert = "dev-certificate"
)

const (
	KILL_REASON_VERIZON_REQUESTED = "Verizon Requested"
)

// Defines the type that represents an inventory device.
type Device struct {
	// Identity. Note that Kind should match the expected ThingSpace
	// device model.
	Id     string `json:"id" csv:"id"`
	IdType string `json:"idType" csv:"idType"`
	Kind   string `json:"kind" csv:"kind"`

	// Security
	CredentialsType string `json:"credentialsType,omitempty" csv:"credentialsType"`
	Secret          string `json:"secret,omitempty"  csv:"-"`
	Certificate     string `json:"certificate,omitempty" csv:"certificate"`
	CertPrivateKey  string `json:"certPrivateKey,omitempty" csv:"-"`
	CertPublicKey   string `json:"certPublicKey,omitempty" csv:"-"`

	// Device Hardware Information. Note that ModelId identifies the
	// Hardware Model (or SKU) Unique Id.
	QrCode                 string `json:"qrCode,omitempty" csv:"qrCode"`
	Family                 string `json:"family, omitempty" csv:"-"`
	Manufacturer           string `json:"manufacturer, omitempty" csv:"manufacturer"`
	SkuNumber              string `json:"skuNumber, omitempty" csv:"skuNumber"`
	ModelId                string `json:"modelId, omitempty" csv:"modelId"`
	SerialNumber           string `json:"serialNumber, omitempty" csv:"serialNumber"`
	Imei                   uint64 `json:"imei, omitempty" csv:"imei"`
	Imsi                   uint64 `json:"imsi, omitempty" csv:"imsi"`
	IccId                  string `json:"iccId, omitempty" csv:"iccId"`
	BluetoothMac           string `json:"bluetoothMac, omitempty" csv:"bluetoothMac"`
	MacAddress             string `json:"macAddress,omitempty" csv:"macAddress"`
	ChipSet                string `json:"chipSet,omitempty" csv:"chipSet"`
	DeviceManufactureDate  string `json:"deviceManufactureDate,omitempty" csv:"deviceManufactureDate"`
	BatteryManufactureDate string `json:"batteryManufactureDate,omitempty" csv:"batteryManufactureDate"`
	BatteryBatchId         string `json:"batteryBatchId,omitempty" csv:"batteryBatchId"`
	HardwareVersion        uint16 `json:"hardwareVersion, omitempty" csv:"hardwareVersion"`
	FirmwareVersion        uint16 `json:"firmwareVersion, omitempty" csv:"firmwareVersion"`

	// Registration Information.
	Registered      bool      `json:"registered"`
	State           string    `json:"state,omitempty" csv:"-"`
	ProvisionedOn   time.Time `json:"provisionedon,omitempty" csv:"-"`
	RegisteredOn    time.Time `json:"registeredon,omitempty" csv:"-"`
	ActivatedOn     time.Time `json:"activatedon,omitempty" csv:"-"`
	DeactivatedOn   time.Time `json:"deactivatedon,omitempty" csv:"-"`
	KillStatus      string    `json:"killstatus,omitempty" csv:"-"`
	KillInitiatorId string    `json:"killinitiatorid,omitempty" csv:"-"`
	KillCounter     uint16    `json:"killcounter, omitempty" csv:"-"`
	KillReason      string    `json:"killreason,omitempty" csv:"-"`
}

// Helper method used to validate device.
func (device Device) Validate() error {
	// If device id is not empty, validate id type.
	if device.Id != "" {
		switch device.IdType {
		case IdTypeMacAddress:
		case IdTypeImsi:
		case IdTypeImei:
		case IdTypeUuid:
		default:
			return fmt.Errorf("The id type is invalid.")
		}
	}

	return nil
}

// Defines the type used to represent device provisioning errors.
type DeviceError struct {
	Device *Device           `json:"device,omitempty"`
	Error  *management.Error `json:"error,omitempty"`
}

// Defines the type used to represent bulk provisioning requests.
type DeviceProvisioningResponse struct {
	Requested     int               `json:"requested"`
	Attempted     int               `json:"attempted"`
	Created       int               `json:"created"`
	NumErrors     int               `json:"numErrors"`
	Devices       []Device          `json:"devices,omitempty"`
	FailedDevices []DeviceError     `json:"failedDevices,omitempty"`
	LastError     *management.Error `json:"lastError,omitempty"`
}
