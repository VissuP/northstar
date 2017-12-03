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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	. "github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/pkg/thingspace/models/ut"
)

// Defines the client interface.
type Client interface {
	CreateEvent(event *Event) *management.Error
	PatchEvent(event *Event) *management.Error
	PatchDevice(device *Device) *management.Error
	RegisterProvider(provider Provider) *management.Error
	QueryAccounts(query Query) ([]Account, *management.Error)
	QueryUsers(query Query) ([]User, *management.Error)
	QueryPlaces(query Query) ([]ut.Place, *management.Error)
	QueryFieldHistory(filter interface{}) ([]Event, *management.Error)
}

// Defines the ThingSpace client.
type ThingSpaceClient struct {
	BaseClient
	authClient   AuthClient
	clientId     string
	clientSecret string
}

// Returns a new ThingSpace client.
func NewThingSpaceClient(thingSpaceHostAndPort string, clientId string, clientSecret string) (Client, error) {
	mlog.Info("NewThingSpaceClient")

	// If no thingspace host and port provided, attempt to get default value.
	if thingSpaceHostAndPort == "" {
		mlog.Info("Getting ThingSpace Host and Port from environment variable.")

		if thingSpaceHostAndPort = os.Getenv("THINGSPACE_HOST_PORT"); thingSpaceHostAndPort == "" {
			mlog.Error("Error, failed to get Thingspace Host and Port.")
			return nil, errors.New("Error, ThingSpace Host and Port environment variable not found or invalid")
		}
	}

	// Create the ThingSpace User Client.
	thingSpaceClient := &ThingSpaceClient{
		BaseClient: BaseClient{
			httpClient: management.NewHttpClient(),
			baseUrl:    fmt.Sprintf("http://%s", thingSpaceHostAndPort),
		},
		authClient:   NewThingSpaceAuthClient(thingSpaceHostAndPort),
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	return thingSpaceClient, nil
}

// Creates a new event in ThingSpace.
func (client *ThingSpaceClient) CreateEvent(event *Event) *management.Error {
	mlog.Info("CreateEvent")

	// TODO: In the future, TS will require a client token in the request.
	// Code is in place but not use for now.

	// Get the service access token.
	//token, err := client.authClient.GetClientToken(client.clientId, client.clientSecret, "")

	// If error, return.
	//if err != nil {
	//	return err
	//}

	// workaround for NPDTHING-3366
	event.Version = EVENT_SCHEMA_VERSION
	event.Kind = EVENT_KIND

	url := client.baseUrl + "/south/v2/events"

	// Request the partial event update.
	if _, err := client.executeHttpMethod("POST", url, "", event); err != nil {
		mlog.Error("Error, failed to execute path event for url %s with error: %s", url, err.Description)
		return err
	}

	return nil
}

// Send a partial event update request to ThingSpace.
func (client *ThingSpaceClient) PatchEvent(event *Event) *management.Error {
	mlog.Debug("PatchEvent - event.Id: %s", event.Id)

	// TODO: In the future, TS will require a client token in the request.
	// Code is in place but not use for now.

	// Get the service access token.
	//token, err := client.authClient.GetClientToken(client.clientId, client.clientSecret, "")

	// If error, return.
	//if err != nil {
	//	return err
	//}

	// workaround for NPDTHING-3366
	// Once fixed, we should take it out from here.
	event.Version = EVENT_SCHEMA_VERSION
	event.Kind = EVENT_KIND

	url := client.baseUrl + BASE_PATH + "events/" + event.Id

	// Request the partial event update.
	if _, err := client.executeHttpMethod("PATCH", url, "", event); err != nil {
		mlog.Error("Error, failed to execute path event for url %s with error: %s", url, err.Description)
		return err
	}

	return nil
}

// Send a partial device update request to ThingSpace.
func (client *ThingSpaceClient) PatchDevice(device *Device) *management.Error {
	mlog.Info("PatchDevice - device.Id: %s", device.Id)

	// TODO: In the future, TS will require a client token in the request.
	// Code is in place but not use for now.

	// Get the service access token.
	//token, err := client.authClient.GetClientToken(client.clientId, client.clientSecret, "")

	// If error, return.
	//if err != nil {
	//	return err
	//}

	url := client.baseUrl + BASE_PATH + "devices/" + device.Id

	// Request the partial device update.
	if _, err := client.executeHttpMethod("PATCH", url, "", device); err != nil {
		mlog.Error("Error, failed to execute patch event for url %s with error: %s", url, err.Description)
		return err
	}

	return nil
}

// Registers the specified Provider with ThingSpace.
func (client *ThingSpaceClient) RegisterProvider(provider Provider) *management.Error {
	mlog.Debug("Register")

	// TODO: In the future, TS will require a client token in the request.
	// Code is in place but not use for now.

	// Get the service access token.
	//token, err := client.authClient.GetClientToken(client.clientId, client.clientSecret, "")

	// If error, return.
	//if err != nil {
	//	return err
	//}
	url := client.baseUrl + BASE_PATH + "providers/" + provider.Id

	if _, err := client.executeHttpMethod("PUT", url, "", provider); err != nil {
		mlog.Error("Error, failed to execute registration for url %s with error: %s", url, err.Description)
		return err
	}

	return nil
}

// Returns the account that meet the specified query.
func (client *ThingSpaceClient) QueryAccounts(query Query) ([]Account, *management.Error) {
	mlog.Debug("QueryAccounts: query:%#v", query)

	url := client.baseUrl + BASE_PATH + "accounts"
	body, err := client.executeHttpMethod("GET", url, "", query)

	if err != nil {
		mlog.Error("Error, failed to get accounts with error: %s", url, err.Description)
		return nil, err
	}

	accounts := make([]Account, 0)
	if err := json.Unmarshal(body, &accounts); err != nil {
		mlog.Error("Error, failed to unmarshal response to accounts with error: %s.", err.Error())
		return nil, management.ErrorInternal
	}

	return accounts, nil
}

// Returns the users that meet the specified query.
func (client *ThingSpaceClient) QueryUsers(query Query) ([]User, *management.Error) {
	mlog.Debug("QueryUsers: query:%#v", query)

	url := client.baseUrl + BASE_PATH + "users"
	body, err := client.executeHttpMethod("GET", url, "", query)

	if err != nil {
		mlog.Error("Error, failed to get users with error: %s", url, err.Description)
		return nil, err
	}

	users := make([]User, 0)
	if err := json.Unmarshal(body, &users); err != nil {
		mlog.Error("Error, failed to unmarshal response to users with error: %s.", err.Error())
		return nil, management.ErrorInternal
	}

	return users, nil
}

// Returns the places that meet the specified query.
func (client *ThingSpaceClient) QueryPlaces(query Query) ([]ut.Place, *management.Error) {
	mlog.Debug("QueryPlaces: query:%#v", query)

	url := client.baseUrl + BASE_PATH + "places"
	body, err := client.executeHttpMethod("GET", url, "", query)

	if err != nil {
		mlog.Error("Error, failed to get places with error: %s", url, err.Description)
		return nil, err
	}

	places := make([]ut.Place, 0)
	if err := json.Unmarshal(body, &places); err != nil {
		mlog.Error("Error, failed to unmarshal response to places with error: %s.", err.Error())
		return nil, management.ErrorInternal
	}

	return places, nil
}
