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
	"strconv"
)

const (
	DEVICES_PATH = "/ottps/v1/devices"
	MODELS_PATH  = "/ottps/v1/models"
)

func (client *Client) CreateOttDevice(accessToken string, device *model.Device, count int) (*model.DeviceProvisioningResponse, *management.Error) {
	mlog.Debug("CreateOttDevice - accessToken: %s", accessToken)

	// Create device
	url := client.baseUrl + DEVICES_PATH + "/" + strconv.Itoa(count)
	body, err := client.post(url, accessToken, device)

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

func (client *Client) UpdateOttDevice(accessToken string, device *model.Device) (*model.Device, *management.Error) {
	mlog.Debug("CreateOttDevice - accessToken: %s", accessToken)

	// Create device
	url := client.baseUrl + DEVICES_PATH
	body, err := client.put(url, accessToken, device)

	if err != nil {
		return nil, err
	}

	createdDevice := &model.Device{}

	if goErr := json.Unmarshal(body, createdDevice); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return createdDevice, nil
}

func (client *Client) DeleteOttDevice(accessToken string, deviceId string) *management.Error {
	mlog.Debug("DeleteOttDevice - accessToken: %s", accessToken)

	// Delete device
	url := client.baseUrl + DEVICES_PATH + "/" + fmt.Sprintf("%s", deviceId)
	return client.delete(url, accessToken)
}

func (client *Client) GetDeviceById(accessToken string, deviceId string) (*model.Device, *management.Error) {
	mlog.Debug("GetOttDeviceById - accessToken: %s", accessToken)

	// Delete device
	url := client.baseUrl + DEVICES_PATH + "/deviceId/" + fmt.Sprintf("%s", deviceId)

	data, serviceError := client.get(url, accessToken, nil)

	if serviceError != nil {
		return nil, serviceError
	}

	returnedDevice := &model.Device{}

	if goErr := json.Unmarshal(data, returnedDevice); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return returnedDevice, nil
}

func (client *Client) CreateOttModel(accessToken string, ottModel *model.Model) (*model.Model, *management.Error) {
	mlog.Debug("CreateOttModel - accessToken: %s", accessToken)

	url := client.baseUrl + MODELS_PATH
	body, err := client.post(url, accessToken, ottModel)

	if err != nil {
		return nil, err
	}

	createdDeviceModel := &model.Model{}

	if goErr := json.Unmarshal(body, createdDeviceModel); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return createdDeviceModel, nil
}

func (client *Client) ListOttModels(accessToken string) ([]model.Model, *management.Error) {
	mlog.Debug("ListOttModels - accessToken: %s", accessToken)

	url := client.baseUrl + MODELS_PATH
	body, err := client.get(url, accessToken, nil)

	if err != nil {
		return nil, err
	}

	receivedModels := make([]model.Model, 0)

	if goErr := json.Unmarshal(body, &receivedModels); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return receivedModels, nil
}

func (client *Client) GetOttModelById(accessToken string, modelId string) (*model.Model, *management.Error) {
	mlog.Debug("GetOttModelById - accessToken: %s", accessToken)

	url := client.baseUrl + MODELS_PATH + "/" + fmt.Sprintf("%s", modelId)

	data, serviceError := client.get(url, accessToken, nil)

	if serviceError != nil {
		return nil, serviceError
	}

	returnedModel := &model.Model{}

	if goErr := json.Unmarshal(data, returnedModel); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return returnedModel, nil
}
