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

import (
	"time"
)

const (
	// Defines the supported message type ids.
	SmsTypeId    string = "sms"
	SmtpTypeId   string = "smtp"
	AppleTypeId  string = "apple"
	GoogleTypeId string = "google"
	RestTypeId   string = "rest"
	TestTypeId   string = "test"
)

// Defines the type used to represent ThingSpace notification
// service message.
type Message struct {
	ClientId         string `json:",omitempty"`
	AccountId        string
	DeviceId         string       `json:",omitempty"`
	TypeId           string       `json:",omitempty"`
	Resource         string       `json:",omitempty"`
	Message          interface{}  `json:",omitempty"`
	NotificationId   string       `json:",omitempty"`
	MessageTimestamp time.Time    `json:",omitempty"`
}

type EmailMessage struct {
	From        string
	To          string
	Subject     string
	MimeVersion string
	ContentType string
	Body        string
}
