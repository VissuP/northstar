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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/utils"
	"strings"
	"time"
)

var (
	pathApiV2         = "/api/v2"
	pathApiV2Devices  = pathApiV2 + "/devices"
	pathApiV2Accounts = pathApiV2 + "/accounts"
	pathApiV2Places   = pathApiV2 + "/places"
	pathApiV2Targets  = pathApiV2 + "/targets"
	pathApiV2Schedule = pathApiV2 + "/schedules"
	pathApiV2Trigger  = pathApiV2 + "/triggers"
	pathApiV2Tags     = pathApiV2 + "/tags"
	pathApiV2Alerts   = pathApiV2 + "/alerts"
)

// Defines the type used by Client applications to perform userapps operations.
type Client struct {
	httpClient *http.Client
	baseUrl    string
}

//////////////////////////////////////////////////////////////////////////////////
// Public Methods

// Returns a new UserApps client for a specific protocol (i.e., https, http), host and port.
func NewClient(protocol, hostAndPort string) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 5000,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				//HTTP methods do not have a timeout by default. Add one.
				Dial: func(netw, addr string) (net.Conn, error) {
					mlog.Event("Connecting to %s", addr)
					c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(config.Configuration.HttpTimeout))
					if err != nil {
						if strings.Compare(reflect.TypeOf(err).String(), "*net.OpError") == 0 {
							mlog.Event("Unable to connect to %s", addr)
						}
						mlog.Event("Disconnected from %s", addr)
						return nil, err
					}
					mlog.Event("Connected to %s", addr)
					return c, nil
				},
				//HTTP methods do not have a timeout by default. Add one for waiting for a response.
				ResponseHeaderTimeout: time.Second * time.Duration(config.Configuration.HttpTimeout),
			},
		},
		baseUrl: fmt.Sprintf("%s://%s", protocol, hostAndPort),
	}
}

// Returns a new device object.
func NewDeviceRegistrationRequest(logs *utils.Logger, qrCode string) *Device {
	configDevice := config.Configuration.Devices[0]
	regReq := &Device{
		Identifiers: Identifiers{
			QRCode: qrCode,
		},
		Identity: Identity{
			Kind: configDevice.Kind,
		},
		ProviderId: configDevice.ProviderId,
	}
	return regReq
}

// Creates a new account.
func (client *Client) CreateAccount(accessToken string, credentials *Credentials) (*Account, *management.Error) {
	mlog.Debug("CreateAccount - accessToken: %s, credentials: %v", accessToken, credentials)

	// Create account - this creates both the Account and the Owning User.
	// Email that we passed gets set as user.CredentialsId as well as user.Email
	// For getting auth token for this user, use user.CredentialsId as login name
	url := client.baseUrl + pathApiV2Accounts

	logs := utils.NewLogger(true)
	body, err := client.post(logs, url, accessToken, credentials)

	if err != nil {
		mlog.Error("CreateAccount Error - " + err.Error())
		return nil, err
	}

	account := &Account{}

	if goErr := json.Unmarshal(body, account); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return account, nil
}

// Deletes the account with the specified id.
func (client *Client) DeleteAccount(accessToken string) *management.Error {
	mlog.Debug("DeleteAccount - accessToken: %s", accessToken)
	url := client.baseUrl + pathApiV2Accounts + "/me"

	logs := utils.NewLogger(true)
	return client.delete(logs, url, accessToken)
}

// GetAccount retrieves the account with the specified id.
func (client *Client) GetAccount(accessToken string) (*Account, *management.Error) {
	mlog.Debug("DeleteAccount - accessToken: %s", accessToken)
	url := client.baseUrl + pathApiV2Accounts + "/me"

	logs := utils.NewLogger(true)
	body, mErr := client.get(logs, url, accessToken, nil)
	if mErr != nil {
		return nil, mErr
	}

	account := &Account{}

	if goErr := json.Unmarshal(body, account); goErr != nil {
		return nil, management.GetInternalError(goErr.Error())
	}

	return account, nil

}

//ListDevices lists devices belonging to the account.
func (client *Client) ListDevices(logs *utils.Logger, accessToken string) ([]Device, *management.Error) {
	logs.LogDebug("ListDevices - accessToken: %s", accessToken)

	url := client.baseUrl + pathApiV2Devices
	body, err := client.get(logs, url, accessToken, nil)
	if err != nil {
		logs.LogError("Failed to list devices. %s", err.Error())
		return nil, err
	}

	devices := []Device{}
	if goErr := json.Unmarshal(body, &devices); goErr != nil {
		return nil, management.GetInternalError(goErr.Error())
	}

	return devices, nil
}

// Register a device in the account associated with the authenticated user.
func (client *Client) RegisterDevice(logs *utils.Logger, accessToken string, req *Device) (*Device, *management.Error) {
	logs.LogDebug("RegisterDevice - accessToken: %s, device %+v", accessToken, req)

	url := client.baseUrl + pathApiV2Devices
	body, err := client.post(logs, url, accessToken, req)

	if err != nil {
		logs.LogError("Failed to register device %+v, err %+v", req, err)
		return nil, err
	}

	logs.LogDebug("Register response %s", string(body))
	registeredDevice := &Device{}

	if goErr := json.Unmarshal(body, registeredDevice); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return registeredDevice, nil
}

func (client *Client) UnregisterDevice(logs *utils.Logger, accessToken string, id string) *management.Error {
	logs.LogDebug("UnregisterDevice - accessToken: %s, id %s", accessToken, id)

	url := client.baseUrl + pathApiV2Devices + "/" + id
	err := client.delete(logs, url, accessToken)

	if err != nil {
		logs.LogError("Failed to unregister device %s, err %+v", id, err)
	}
	return err
}

// will be deprecated soon
// use GetApiDevice instead
func (client *Client) GetDevice(logs *utils.Logger, accessToken string, id string) (*Device, *management.Error) {
	logs.LogDebug("GetDevice - accessToken: %s, id %s", accessToken, id)

	url := client.baseUrl + pathApiV2Devices + "/" + id
	body, err := client.get(logs, url, accessToken, nil)

	if err != nil {
		logs.LogError("Failed to get device %s, err %+v", id, err)
		return nil, err
	}
	logs.LogDebug("GetDevice response %s", string(body))

	device := &Device{}
	if goErr := json.Unmarshal(body, device); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return device, nil
}

// similar to GetDevice above, except is returns api.Device
// all new tests should use this method
func (client *Client) GetApiDevice(logs *utils.Logger, accessToken string, id string) (*api.Device, *management.Error) {
	logs.LogDebug("GetApiDevice - accessToken: %s, id %s", accessToken, id)

	url := client.baseUrl + pathApiV2Devices + "/" + id
	body, err := client.get(logs, url, accessToken, nil)

	if err != nil {
		logs.LogError("Failed to get device %s, err %+v", id, err)
		return nil, err
	}
	logs.LogDebug("GetDevice response %s", string(body))

	device := &api.Device{}
	if goErr := json.Unmarshal(body, device); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return device, nil
}

// returns the response event when a get field action TS API is executed
func (client *Client) GetDeviceField(logs *utils.Logger, accessToken string, id string, field string) (*api.Event, *management.Error) {
	logs.LogDebug("GetDeviceField - accessToken: %s, id:%s, field:%s", accessToken, id, field)

	url := client.baseUrl + pathApiV2Devices + "/" + id + "/fields/" + field + "/actions/get"

	body, mErr := client.post(logs, url, accessToken, nil)

	if mErr != nil {
		logs.LogError("Failed to get field:%s, for device %s, err %+v", field, id, mErr)
		return nil, mErr
	}

	var event api.Event
	err := json.Unmarshal(body, &event)
	if err != nil {
		logs.LogInfo("user event = %s", err.Error())
		mErr := management.GetInternalError("GetDeviceField response UnMarshal Error")
		return nil, mErr
	}
	return &event, nil
}

// returns the response map when a get field action TS API is executed
func (client *Client) GetDeviceFieldValue(logs *utils.Logger, accessToken string, id string, field string) (map[string]interface{}, *management.Error) {
	logs.LogDebug("GetDeviceFieldValue - accessToken: %s, id:%s, field:%s", accessToken, id, field)

	url := client.baseUrl + pathApiV2Devices + "/" + id + "/fields/" + field + "/actions/get"

	body, mErr := client.post(logs, url, accessToken, nil)

	if mErr != nil {
		logs.LogError("Failed to get field:%s, for device %s, err %+v", field, id, mErr)
		return nil, mErr
	}

	logs.LogDebug("GetDeviceFieldValue %s, response %s", field, body)

	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: err.Error(),
		}

	}

	m := f.(map[string]interface{})

	return m, nil
}

func (client *Client) CancelDeviceFieldEvent(logs *utils.Logger, accessToken string, id string, field string) (*api.Event, *management.Error) {
	logs.LogDebug("CancelDeviceFieldEvent - accessToken: %s, id:%s, field:%s", accessToken, id, field)

	url := client.baseUrl + pathApiV2Devices + "/" + id + "/fields/" + field + "/actions/cancel"

	body, mErr := client.post(logs, url, accessToken, nil)

	if mErr != nil {
		logs.LogError("Failed to cancel field:%s, for device %s, err %+v", field, id, mErr)
		return nil, mErr
	}

	var event api.Event
	err := json.Unmarshal(body, &event)
	if err != nil {
		logs.LogInfo("CancelDeviceFieldEvent response UnMarshal Error", err.Error())
		mErr := management.GetInternalError("CancelDeviceFieldEvent response UnMarshal Error")
		return nil, mErr
	}
	return &event, nil
}

// this method will be deprecated. Use SetFieldEvent instead
func (client *Client) SetDeviceField(logs *utils.Logger, accessToken string, id string, field string, fieldValue interface{}) (map[string]interface{}, *management.Error) {
	logs.LogDebug("SetDeviceFieldValue - accessToken: %s, id:%s, field:%s, value:%+v", accessToken, id, field, fieldValue)

	url := client.baseUrl + pathApiV2Devices + "/" + id + "/fields/" + field + "/actions/set"

	body, mErr := client.post(logs, url, accessToken, fieldValue)

	if mErr != nil {
		logs.LogError("Failed to set field %s, device %s, err %+v", field, id, mErr)
		return nil, mErr
	}

	logs.LogDebug("SetDeviceFieldValue response %s", body)
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: err.Error(),
		}

	}

	m := f.(map[string]interface{})

	return m, nil
}

// similar to SetDeviceField above, except it returns a api.Event object.
// all new tests should use this method.
func (client *Client) SetFieldEvent(logs *utils.Logger, accessToken string, id string, field string, fieldValue interface{}) (*api.Event, *management.Error) {
	logs.LogDebug("SetFieldEvent - accessToken: %s, id:%s, field:%s, value:%+v", accessToken, id, field, fieldValue)

	url := client.baseUrl + pathApiV2Devices + "/" + id + "/fields/" + field + "/actions/set"
	body, mErr := client.post(logs, url, accessToken, fieldValue)
	if mErr != nil {
		logs.LogError("Failed to set field %s, device %s, err %+v", field, id, mErr)
		return nil, mErr
	}

	event := &api.Event{}
	if goErr := json.Unmarshal(body, event); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}
	return event, nil
}

func (client *Client) DeactivateDevice(logs *utils.Logger, accessToken string, id string) (*Device, *management.Error) {
	logs.LogDebug("DeactivateDevice - accessToken: %s, id %s", accessToken, id)

	url := client.baseUrl + pathApiV2Devices + "/" + id + "/actions/deactivate"
	if _, err := client.post(logs, url, accessToken, nil); err != nil {
		logs.LogError("Failed to get device %s, err %v", id, err)
		return nil, err
	}

	// Doesn't look like this API returns the device object back. Do a GET instead
	device, err := client.GetDevice(logs, accessToken, id)
	if err != nil {
		logs.LogError("Failed to get device after deactivation. deviceId:%s, err %v", id, err)
		return nil, err
	}

	return device, nil
}

func (client *Client) GetAlerts(logs *utils.Logger, accessToken string) ([]Alert, *management.Error) {
	logs.LogDebug("GetAlerts -- accessToken: %s", accessToken)
	url := client.baseUrl + pathApiV2Alerts

	body, serviceErr := client.get(logs, url, accessToken, nil)
	if serviceErr != nil {
		logs.LogError("Error, failed to get alerts with error %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedAlerts []Alert
	err := json.Unmarshal(body, &returnedAlerts)
	if err != nil {
		logs.LogError("Error, failed to unmarshal alerts JSON with error %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	return returnedAlerts, nil

}

func (client *Client) UpdateAlerts(logs *utils.Logger, accessToken string, alerts []Alert) *management.Error {
	logs.LogDebug("UpdateAlerts -- accessToken: %s", accessToken)
	url := client.baseUrl + pathApiV2Alerts

	// Note that the body of a patch alert is a filter and the attributes of the alert
	// to modify. On the other hand, they do not support filters that can modify more
	// than one field at a time. This might be supported at some point but for now
	// doing one request per alert.

	for _, alert := range alerts {

		// Create the filter.
		filter := map[string]interface{}{
			"$filter": map[string]interface{}{
				"id": alert.Id,
			},
			"isread": alert.IsRead,
		}

		if _, serviceErr := client.patch(logs, url, accessToken, filter); serviceErr != nil {
			logs.LogError("Error, failed to update alerts with error %s", serviceErr.Description)
			return serviceErr
		}
	}

	return nil
}
func (client *Client) CreatePlace(logs *utils.Logger, accessToken string, place Place) (*Place, *management.Error) {
	//http://confluence.verizon.com/display/NPDTHING/Creating+and+Managing+a+Place
	logs.LogDebug("CreatePlace -- accessToken: %s, place: %+v", accessToken, place)
	url := client.baseUrl + pathApiV2Places

	body, serviceErr := client.post(logs, url, accessToken, place)
	if serviceErr != nil {
		logs.LogError("Error, Failed to create place with error %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedPlace Place
	err := json.Unmarshal(body, &returnedPlace)
	if err != nil {
		logs.LogError("Error %s encountered while unmarshaling JSON", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	return &returnedPlace, nil
}

func (client *Client) DeletePlace(logs *utils.Logger, accessToken string, placeID string) *management.Error {
	logs.LogDebug("DeletePlace -- accessToken: %s, placeID: %s", accessToken, placeID)
	url := client.baseUrl + pathApiV2Places

	serviceErr := client.delete(logs, url, accessToken)
	if serviceErr != nil {
		logs.LogError("Failed to delete place with error %s", serviceErr.Description)
		return serviceErr
	}

	return nil
}
func (client *Client) CreateTarget(logs *utils.Logger, accessToken string, target Target) (*Target, *management.Error) {
	//http://confluence.verizon.com/display/NPDTHING/Creating+and+Managing+a+Target
	logs.LogDebug("CreateTarget -- accessToken:%s, target: %+v", accessToken, target)
	url := client.baseUrl + pathApiV2Targets

	body, serviceErr := client.post(logs, url, accessToken, target)
	if serviceErr != nil {
		logs.LogError("Failed to create target with error %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedTarget Target
	err := json.Unmarshal(body, &returnedTarget)
	if err != nil {
		logs.LogError("Error %s encountered while unmarshaling JSON", err)
		return nil, management.GetInternalError(err.Error())
	}

	return &returnedTarget, nil
}
func (client *Client) DeleteTarget(logs *utils.Logger, accessToken string, targetID string) *management.Error {
	logs.LogDebug("DeleteTarget -- access token: %s, targetID: %s", accessToken, targetID)
	url := client.baseUrl + pathApiV2Targets + "/" + targetID

	serviceErr := client.delete(logs, url, accessToken)
	if serviceErr != nil {
		logs.LogError("Failed to delete target with error %s", serviceErr.Description)
		return serviceErr
	}

	return nil
}

func (client *Client) CreateTag(logs *utils.Logger, accessToken string, tag Tag) (*Tag, *management.Error) {
	//http://confluence.verizon.com/display/NPDTHING/Creating+and+Managing+a+Tag
	logs.LogDebug("CreateTag -- accessToken: %s, tag:%+v", accessToken, tag)
	url := client.baseUrl + pathApiV2Tags

	body, serviceErr := client.post(logs, url, accessToken, tag)

	if serviceErr != nil {
		logs.LogError("Failed to create tag with error: %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedTag Tag
	err := json.Unmarshal(body, &returnedTag)
	if err != nil {
		logs.LogError("Error, %s encountered while unmarhaling JSON", err)
		return nil, management.GetInternalError(err.Error())
	}

	return &returnedTag, nil
}

func (client *Client) DeleteTag(logs *utils.Logger, accessToken string, tagID string) *management.Error {
	logs.LogDebug("DeleteTag -- accessToken: %s, tagID: %s", accessToken, tagID)
	url := client.baseUrl + pathApiV2Tags + "/" + tagID

	serviceErr := client.delete(logs, url, accessToken)
	if serviceErr != nil {
		logs.LogError("Failed to delete tag with error: %s", serviceErr.Description)
		return serviceErr
	}

	return nil
}

func (client *Client) CreateSchedule(logs *utils.Logger, accessToken string, schedule Schedule) (*Schedule, *management.Error) {
	//http://confluence.verizon.com/display/NPDTHING/Creating+and+Managing+a+Schedule
	logs.LogDebug("CreateSchedule -- accessToken: %s, schedule: %+v", accessToken, schedule)
	url := client.baseUrl + pathApiV2Schedule

	body, serviceErr := client.post(logs, url, accessToken, schedule)
	if serviceErr != nil {
		logs.LogError("Failed to create schedule with error: %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedSchedule Schedule
	err := json.Unmarshal(body, &returnedSchedule)
	if err != nil {
		logs.LogError("Error %s encountered while unmarshaling JSON", err)
		return nil, management.GetInternalError(err.Error())
	}

	return &returnedSchedule, nil
}

func (client *Client) DeleteSchedule(logs *utils.Logger, accessToken string, scheduleID string) *management.Error {
	logs.LogDebug("DeleteSchedule -- accessToken: %s, scheduleID: %s", accessToken, scheduleID)
	url := client.baseUrl + pathApiV2Schedule + "/" + scheduleID

	serviceErr := client.delete(logs, url, accessToken)
	if serviceErr != nil {
		logs.LogError("Failed to delete tag with error: %s", serviceErr.Description)
		return serviceErr
	}

	return nil
}

func (client *Client) CreateTrigger(logs *utils.Logger, accessToken string, trigger Trigger) (*Trigger, *management.Error) {
	//http://confluence.verizon.com/display/NPDTHING/Creating+and+Managing+a+Trigger
	logs.LogDebug("CreateTrigger -- accessToken: %s, trigger: %+v", accessToken, trigger)
	url := client.baseUrl + pathApiV2Trigger

	body, serviceErr := client.post(logs, url, accessToken, trigger)
	if serviceErr != nil {
		logs.LogError("Failed to create trigger with error: %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedTrigger Trigger
	err := json.Unmarshal(body, &returnedTrigger)
	if err != nil {
		logs.LogError("Error %s encountered while unmarshaling JSON", err)
		return nil, management.GetInternalError(err.Error())
	}
	return &returnedTrigger, nil
}

func (client *Client) GetTriggers(logs *utils.Logger, accessToken string) ([]Trigger, *management.Error) {
	//http://confluence.verizon.com/display/NPDTHING/Creating+and+Managing+a+Trigger
	logs.LogDebug("GetTriggers -- accessToken: %s", accessToken)
	url := client.baseUrl + pathApiV2Trigger

	body, serviceErr := client.get(logs, url, accessToken, nil)
	if serviceErr != nil {
		logs.LogError("Failed to create trigger with error: %s", serviceErr.Description)
		return nil, serviceErr
	}

	var returnedTriggers []Trigger
	err := json.Unmarshal(body, &returnedTriggers)
	if err != nil {
		logs.LogError("Error %s encountered while unmarshaling JSON", err)
		return nil, management.GetInternalError(err.Error())
	}
	return returnedTriggers, nil
}

func (client *Client) DeleteTrigger(logs *utils.Logger, accessToken string, triggerID string) *management.Error {
	mlog.Debug("DeleteTrigger -- accessToken: %s, trigger: %+v", accessToken, triggerID)
	url := client.baseUrl + pathApiV2Trigger + "/" + triggerID

	serviceErr := client.delete(logs, url, accessToken)
	if serviceErr != nil {
		logs.LogError("Failed to delete trigger with error: %s", serviceErr.Description)
		return serviceErr
	}

	return nil
}

func (client *Client) GetLocationHistory(logs *utils.Logger, accessToken string, deviceID string) ([]HistoryEvent, *management.Error) {
	logs.LogDebug("GetLocationHistory -- accessToken: %s", accessToken)
	url := client.baseUrl + pathApiV2Devices + "/" + deviceID + "/fields/location/actions/history"

	body, serviceErr := client.post(logs, url, accessToken, nil)
	if serviceErr != nil {
		logs.LogError("Failed to get locations with error: %s", serviceErr.Description)
		return nil, serviceErr
	}

	var historyEvents []HistoryEvent
	err := json.Unmarshal(body, &historyEvents)
	if serviceErr != nil {
		logs.LogError("Error %s encountered while unmarshaling JSON", err)
		return nil, management.GetInternalError(err.Error())
	}

	return historyEvents, nil
}

func (client *Client) GetFieldHistory(logs *utils.Logger, accessToken string, deviceID string, field string) ([]api.Event, *management.Error) {
	logs.LogDebug("GetFieldHistory field:%s -- accessToken: %s", field, accessToken)
	return client.QueryFieldHistory(logs, accessToken, deviceID, field, nil)
}

func (client *Client) QueryFieldHistory(logs *utils.Logger, accessToken string, deviceID string, field string, query *api.Query) ([]api.Event, *management.Error) {
	logs.LogDebug("GetFieldHistory field:%s -- accessToken: %s, query:%+v", field, accessToken, query)
	url := client.baseUrl + pathApiV2Devices + "/" + deviceID + "/fields/" + field + "/actions/history"

	body, serviceErr := client.post(logs, url, accessToken, query)
	if serviceErr != nil {
		logs.LogError("Failed to get field:%s history with error: %s", field, serviceErr.Description)
		return nil, serviceErr
	}

	var historyEvents []api.Event
	err := json.Unmarshal(body, &historyEvents)
	if serviceErr != nil {
		logs.LogError("Error %s encountered while unmarshaling field:%s history JSON", err, field)
		return nil, management.GetInternalError(err.Error())
	}

	return historyEvents, nil
}

//////////////////////////////////////////////////////////////////////////////////
// Private Methods

// Returns the body returned buy the GET on the specified url.
func (client *Client) get(logs *utils.Logger, url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	logs.LogDebug("GET %s", url)

	startTime := time.Now()
	var request *http.Request

	if resourceObject != nil {
		// Parse the resource object as request body.
		requestBody, goErr := json.Marshal(resourceObject)

		if goErr != nil {
			return nil, &management.Error{
				HttpStatus:  http.StatusInternalServerError,
				Id:          fmt.Sprintf("go_error"),
				Description: goErr.Error(),
			}
		}

		request, _ = http.NewRequest("GET", url, bytes.NewBuffer(requestBody))
	} else {
		logs.LogDebug("URL: %s", url)
		request, _ = http.NewRequest("GET", url, nil)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	// If token was provided, add to the request header
	if accessToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	response, goErr := client.httpClient.Do(request)

	// Check for go errors (e.g., Connectin refuse, etc)
	if goErr != nil {
		err := &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}

		return nil, err
	}

	// Check for UserApps Service errors.
	if response.StatusCode != http.StatusOK {
		err = &management.Error{
			HttpStatus:  response.StatusCode,
			Id:          fmt.Sprintf("%d", response.StatusCode),
			Description: response.Status,
		}

		if response.Body != nil {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, err)
		}

		logs.LogError("Error, GET failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Get successful resonse body.
	defer response.Body.Close()
	body, _ = ioutil.ReadAll(response.Body)

	// Collect some metrics
	latency := time.Since(startTime)
	logs.LogDebug("GET %s Latency: %s", url, latency.String())

	return body, nil
}

// Returns the body returned after updation the specified resource object on the specified url.
func (client *Client) post(logs *utils.Logger, url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	logs.LogDebug("POST %s", url)

	startTime := time.Now()

	// Parse the resource object as request body.
	requestBody, goErr := json.Marshal(resourceObject)
	logs.LogDebug("Request Body %s", requestBody)

	if goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	// If token was provided, add to the request header
	if accessToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	response, goErr := client.httpClient.Do(request)

	// Check go errors executing the request.
	if goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	// Check for HTTP client or service errors.
	if response.StatusCode >= 300 {
		err := &management.Error{
			HttpStatus:  response.StatusCode,
			Id:          fmt.Sprintf("%d", response.StatusCode),
			Description: response.Status,
		}

		if response.Body != nil {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, err)
		}

		logs.LogError("Error, POST failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Parse the body if and only if one was provided.
	if response.Body != nil {
		defer response.Body.Close()
		body, _ = ioutil.ReadAll(response.Body)
		logs.LogDebug("Response %s", body)

		// Collect some metrics
		latency := time.Since(startTime)
		logs.LogDebug("POST %s Latency: %s", url, latency.String())

		return body, nil
	}

	// Otherwisem we got a sucessfull response with no body.
	return nil, nil
}

// Returns the body returned after updation the specified resource object on the specified url.
func (client *Client) put(logs *utils.Logger, url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	logs.LogDebug("PUT %s", url)

	startTime := time.Now()

	// Parse the resource object as request body.
	requestBody, goErr := json.Marshal(resourceObject)

	if goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	// If token was provided, add to the request header
	if accessToken != "" {
		logs.LogDebug("Setting request authorization token.")
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	response, goErr := client.httpClient.Do(request)

	// Check go errors executing the request.
	if goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	// Check for HTTP client or service errors.
	if response.StatusCode >= 300 {
		err := &management.Error{
			HttpStatus:  response.StatusCode,
			Id:          fmt.Sprintf("%d", response.StatusCode),
			Description: response.Status,
		}

		if response.Body != nil {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, err)
		}

		logs.LogError("Error, PUT failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Parse the body if and only if one was provided.
	if response.Body != nil {
		defer response.Body.Close()
		body, _ = ioutil.ReadAll(response.Body)

		// Collect some metrics
		latency := time.Since(startTime)
		logs.LogDebug("PUT %s Latency: %s", url, latency.String())

		return body, nil
	}

	// Otherwisem we got a sucessfull response with no body.
	return nil, nil
}

func (client *Client) patch(logs *utils.Logger, url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	logs.LogDebug("PATCH %s", url)

	startTime := time.Now()

	// Parse the resource object as request body.
	requestBody, goErr := json.Marshal(resourceObject)
	logs.LogDebug("Request Body %s", requestBody)

	if goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	request, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	// If token was provided, add to the request header
	if accessToken != "" {
		logs.LogDebug("Setting request authorization token.")
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	response, goErr := client.httpClient.Do(request)

	// Check go errors executing the request.
	if goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	// Check for HTTP client or service errors.
	if response.StatusCode >= 300 {
		err := &management.Error{
			HttpStatus:  response.StatusCode,
			Id:          fmt.Sprintf("%d", response.StatusCode),
			Description: response.Status,
		}

		if response.Body != nil {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, err)
		}

		logs.LogError("Error, PATCH failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Parse the body if and only if one was provided.
	if response.Body != nil {
		defer response.Body.Close()
		body, _ = ioutil.ReadAll(response.Body)

		// Collect some metrics
		latency := time.Since(startTime)
		logs.LogDebug("PATCH %s Latency: %s", url, latency.String())

		return body, nil
	}

	// Otherwisem we got a sucessfull response with no body.
	return nil, nil
}

// Deletes the resource object associated with the specified url.
func (client *Client) delete(logs *utils.Logger, url string, accessToken string) (err *management.Error) {
	logs.LogDebug("DELETE %s", url)

	startTime := time.Now()
	request, _ := http.NewRequest("DELETE", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Close = true

	// If token was provided, add to the request header
	if accessToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	response, goErr := client.httpClient.Do(request)

	// Check go errors executing the request.
	if goErr != nil {
		return &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	// Check for HTTP client or service errors.
	if response.StatusCode >= 300 {
		err := &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("%d", response.StatusCode),
			Description: response.Status,
		}

		if response.Body != nil {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, err)
		}

		logs.LogError("Error, DELETE failed with error: (%s) %s", err.Id, err.Description)

		return err
	}

	ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	// Collect some metrics
	latency := time.Since(startTime)
	logs.LogDebug("DELETE %s Latency: %s", url, latency.String())

	return nil
}
