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

package thingspace

import (
	"fmt"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"errors"
)

// Defines the notif client interface.
type NotifClient interface {
	SendMessage(message api.Message, accessToken string) error
	SendMessageViaSouth(message api.Message) error
}

// Defines the ThingSpace notif client.
type ThingSpaceNotifClient struct {
	BaseClient
}

// Returns a new ThingSpace Notif Service client.
func NewTSNotifClient(thingSpaceNotifHostAndPort string) NotifClient {
	return NewTSNotifClientWithProtocol("http", thingSpaceNotifHostAndPort)
}

// Returns a new ThingSpace Notif Service client.
func NewTSNotifClientWithProtocol(protocol, hostAndPort string) NotifClient {
	notifClient := &ThingSpaceNotifClient{
		BaseClient: BaseClient{
			httpClient: management.NewHttpClient(),
			baseUrl:    fmt.Sprintf("%s://%s", protocol, hostAndPort),
		},
	}

	return notifClient
}

// Send message via notifrest service via north interface.
func (client *ThingSpaceNotifClient) SendMessage(message api.Message, accessToken string) error {
	mlog.Debug("SendMessage")

	url := client.baseUrl + api.NOTIF_BASE_PATH
	_, mErr := client.executeHttpMethod("POST", url, accessToken, message)

	if mErr != nil {
		mlog.Error("Error, failed to send message via notification service with error: %s", url, mErr.Description)
		return errors.New(mErr.String())
	}

	return nil
}

// Send message via notifrest service via south interface.
func (client *ThingSpaceNotifClient) SendMessageViaSouth(message api.Message) error {
	mlog.Debug("SendMessageViaSouth")

	url := client.baseUrl + api.NOTIF_SOUTH_BASE_PATH
	_, mErr := client.executeHttpMethod("POST", url, "", message)

	if mErr != nil {
		mlog.Error("Error, failed to send message via south notification service with error: %s", mErr.Description)
		return errors.New(mErr.String())
	}

	return nil
}
