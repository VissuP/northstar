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

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/utprovisionlib/model"
)

const (
	DEVICES_PATH_V2 = "/api/v2/inventory/devices"
)

func (client *Client) CreateOttDeviceV2(accessToken string, devices []model.Device) (*model.DeviceProvisioningResponse, *management.Error) {
	mlog.Debug("CreateOttDeviceV2 - accessToken: %s", accessToken)

	// Create device
	url := client.baseUrl + DEVICES_PATH_V2
	body, err := client.post(url, accessToken, devices)

	if err != nil {
		return nil, err
	}

	rsp := &model.DeviceProvisioningResponse{}

	if goErr := json.Unmarshal(body, rsp); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return rsp, nil
}

func (client *Client) DeleteOttDeviceV2(accessToken string, deviceId string) *management.Error {
	mlog.Debug("DeleteOttDeviceV2 - accessToken: %s", accessToken)

	// Delete device
	url := client.baseUrl + DEVICES_PATH_V2 + "/" + fmt.Sprintf("%s", deviceId)
	return client.delete(url, accessToken)
}

func (client *Client) PatchOttDeviceV2(accessToken string, deviceId string, devices *model.Device) (*model.Device, *management.Error) {
	mlog.Debug("PatchOttDeviceV2 - accessToken: %s", accessToken)

	// Patch device
	url := client.baseUrl + DEVICES_PATH_V2 + "/" + fmt.Sprintf("%s", deviceId)
	body, err := client.patch(url, accessToken, devices)

	if err != nil {
		return nil, err
	}

	rsp := &model.Device{}

	if goErr := json.Unmarshal(body, rsp); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return rsp, nil
}
