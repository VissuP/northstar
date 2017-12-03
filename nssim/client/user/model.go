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

package user

import (
	"time"
)

const (
	// Defines the strings used to identify a resource as a Dakota Device
	// resource.
	DEVICE_ULTRATAG   string = "ts.device.ultratag"
	DEVICE_ULTRAMODEM string = "ts.device.ultramodem"
	DEVICE_OTT        string = "ts.device.ott"
	DEVICE_SMSNAS     string = "ts.device.smsnas"
)

const (
	DEVICE_STATE_REGISTERED  = "registered"
	DEVICE_STATE_PROVISIONED = "provisioned"
	DEVICE_STATE_ACTIVATED   = "activated"
	DEVICE_STATE_DEACTIVATED = "deactivated"
	DEVICE_STATE_READY       = "ready"
)
const (
	// event states
	TS_EVENT_STATE_PENDING  = "pending"
	TS_EVENT_STATE_UPDATE   = "update"
	TS_EVENT_STATE_FAILED   = "failed"
	TS_EVENT_STATE_CANCELED = "canceled"
)

type shape string

const (
	SHAPE_CIRCLE   shape = "circle"
	SHAPE_LOCATION shape = "location"
)

type targetScheme string

const (
	DESTINATION_APNS targetScheme = "apns"
	DESTINATION_GCMS targetScheme = "gcms"
	DESTINATION_SMS  targetScheme = "sms"
	DESTINATION_SMTP targetScheme = "smtp"
	DESTINATION_REST targetScheme = "rest"
)

type Day string

const (
	DAY_MONDAY    Day = "Monday"
	DAY_TUESDAY   Day = "Tuesday"
	DAY_WEDNESDAY Day = "Wednesday"
	DAY_THURSDAY  Day = "Thursday"
	DAY_FRIDAY    Day = "Friday"
	DAY_SATURDAY  Day = "Saturday"
	DAY_SUNDAY    Day = "Sunday"
)

type Credentials struct {
	Email               string `json:"email,omitempty"`
	Password            string `json:"password,omitempty"`
	NoEmailVerification bool   `json:"noEmailVerification"`
}

type FieldValue struct {
	Value     interface{} `json:"value,omitempty"`
	StartTime string      `json:"starttime,omitempty"`
}

type Account struct {
	Identity
	CustomDataId string `json:"customdataid,omitempty"`
	OwnerId      string `json:"ownerid,omitempty"`
	State        string `json:"state,omitempty"`
	Email        string
	Password     string
}

type Device struct {
	Identity
	Foreign
	ProviderId string `json:"providerid,omitempty"`
	State      string `json:"state,omitempty"`
	Identifiers
	UltraTag UltraTag `json:"fields,omitempty"`
}

type Identity struct {
	Id          string    `json:"id,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Version     string    `json:"version,omitempty"`
	VersionId   string    `json:"versionid,omitempty"`
	CreatedOn   time.Time `json:"-"`
	LastUpdated time.Time `json:"-"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
}

type Foreign struct {
	ForeignId string   `json:"foreignid,omitempty"`
	Tags      []string `json:"tagids,omitempty"`
}

type Identifiers struct {
	QRCode          string `json:"qrcode,omitempty"`
	Mac             string `json:"mac,omitempty"`
	Imei            uint64 `json:"imei,omitempty"`
	Imsi            uint64 `json:"imsi,omitempty"`
	MsIsdn          string `json:"msisdn,omitempty"`
	MeId            string `json:"meid,omitempty"`
	IccId           string `json:"iccid,omitempty"`
	ProductModel    string `json:"productmodel,omitempty"`
	HardwareVersion string `json:"hardwareversion,omitempty"`
}

type Location struct {
	Kind         string  `json:"kind,omitempty"`
	Id           string  `json:"id,omitempty"`
	Version      string  `json:"version,omitempty"`
	CreatedAt    string  `json:"createdAt,omitempty"`
	Address      string  `json:"address,omitempty"`
	Altitude     float64 `json:"altitude,omitempty"`
	Longitude    float64 `json:"longitude" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Accuracy     float64 `json:"accuracy,omitempty"`
	Place        Place   `json:"place,omitempty"`
	LocationType string  `json:"type,omitempty"`
}

type Event struct {
	Identity
	Foreign
	event
	UltraTag
}

type event struct {
	Action   string `json:"action,omitempty"`
	State    string `json:"state,omitempty"`
	DeviceId string `json:"deviceid,omitempty"`
}

type Tag struct {
	Id            string    `json:"id,omitempty"`
	Kind          string    `json:"kind,omitempty"`
	Version       string    `json:"version,omitempty"`
	CreatedOn     time.Time `json:"-"`
	LastUpdatedOn time.Time `json:"-"`
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Color         string    `json:"color,omitempty"`
	CustomDataId  string    `json:"customdataid,omitempty"`
	IconId        string    `json:"iconid,omitempty"`
	ImageId       string    `json:"imageid,omitempty"`
}

type Place struct {
	Id          string    `json:"id,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Version     string    `json:"version,omitempty"`
	CreatedOn   time.Time `json:"-"`
	LastUpdated time.Time `json:"-"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Address
	Shape     shape    `json:"shape,omitempty"`
	Radius    float64  `json:"radius,omitempty"`
	Latitude  float64  `json:"latitude,omitempty"`
	Longitude float64  `json:"longitude,omitempty"`
	Tags      []string `json:"tagids,omitempty"`
}

type Address struct {
	AddressLine1 string `json:"addressline1,omitempty"`
	AddressLine2 string `json:"addressline2,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	State        string `json:"state,omitempty"`
	Zip          string `json:"zip,omitempty"`
}
type Schedule struct {
	Id             string    `json:"id,omitempty"`
	Kind           string    `json:"kind,omitempty"`
	Version        string    `json:"version,omitempty"`
	CreatedOn      time.Time `json:"-"`
	LastUpdated    time.Time `json:"-"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	Days           []Day     `json:"days,omitempty"`
	StartTime      time.Time `json:"starttime,omitempty"`
	EndTime        time.Time `json:"endtime,omitempty"`
	ScheduledStart time.Time `json:"scheduledstart,omitempty"`
	ScheduledStop  time.Time `json:"scheduledstop,omitempty"`
	Tags           []string  `json:"tagids,omitempty"`
}

type Target struct {
	Id            string       `json:"id,omitempty"`
	Kind          string       `json:"kind,omitempty"`
	Version       string       `json:"version,omitempty"`
	CreatedOn     time.Time    `json:"-"`
	LastUpdated   time.Time    `json:"-"`
	Name          string       `json:"name,omitempty"`
	Description   string       `json:"description,omitempty"`
	AddressScheme targetScheme `json:"addressscheme,omitempty"`
	Address       string       `json:"address,omitempty"`
	Tags          []string     `json:"tagids,omitempty"`
}

type Trigger struct {
	Id          string    `json:"id,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Version     string    `json:"version,omitempty"`
	CreatedOn   time.Time `json:"-"`
	LastUpdated time.Time `json:"-"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Condition   string    `json:"condition,omitempty"`
	Operation   string    `json:"operation,omitempty"`
	State       bool      `json:"-"`
	Context     string    `json:"context,omitempty"`
	ScheduleId  string    `json:"scheduleid,omitempty"`
	DeviceId    string    `json:"deviceid,omitempty"` //deprecated
	FieldId     string    `json:"fieldid,omitempty"`  //deprecated
	Tags        []string  `json:"tagids,omitempty"`
}

type Alert struct {
	Id          string    `json:"id,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Version     string    `json:"version,omitempty"`
	VersionId   string    `json:"versionid,omitempty"`
	CreatedOn   time.Time `json:"-"`
	LastUpdated time.Time `json:"-"`
	ForeignId   string    `json:"foreignid,omitempty"` // Specifies the Account ID
	TagIds      []string  `json:"-"`
	TriggerId   string    `json:"triggerid,omitempty"`
	DeviceId    string    `json:"deviceid,omitempty"`
	IsRead      bool      `json:"isread,omitempty"`
	Content     string    `json:"-"`
}

type HistoryEvent struct {
	Id          string    `json:"id,omitempty"`
	CreatedOn   time.Time `json:"-"`
	LastUpdated time.Time `json:"-"`
	ForeignId   string    `json:"foreignid,omitempty"`
	DeviceId    string    `json:"deviceid,omitempty"`
	ModelId     string    `json:"modelid,omitempty"`
	State       string    `json:"state,omitempty"`
	VersionId   string    `json:"versionid,omitempty"`
	Fields      UltraTag  `json:"fields,omitempty"`
}

type UltraTag struct {
	ReportingInterval FieldIntValue `json:"reportinginterval,omitempty"`
	Battery           int           `json:"battery,omitempty"`
	Location          *Location     `json:"location,omitempty"`
	AirplaneMode      FieldIntValue `json:"airplanemode,omitempty"`
	ActivationDate    string        `json:"activationdate,omitempty"`
	FirmwareVersion   string        `json:"firmwareversion,omitempty"`
}

type FieldIntValue struct {
	Value     int    `json:"value,omitempty"`
	StartTime string `json:"starttime,omitempty"`
}
