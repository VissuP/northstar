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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/utprovisionlib/model"
	"time"
)

// Defines the type used by Client applications to perform utprovision operations.
type Client struct {
	httpClient *http.Client
	baseUrl    string
}

//////////////////////////////////////////////////////////////////////////////////
// Public Methods

// Returns a new UtProvision client for a specific protocol (i.e., https, http), host and port.
func NewClient(protocol, hostAndPort string) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 5000,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		},
		baseUrl: fmt.Sprintf("%s://%s", protocol, hostAndPort),
	}
}

// Returns the service health.
func (client *Client) GetHealth() (health *management.Health, err *management.Error) {
	mlog.Debug("GetHealth")

	// Get the service health. Note that no access token needed for this request.
	url := client.baseUrl + "/management/health"
	body, err := client.get(url, "", nil)

	// If error, return.
	if err != nil {
		mlog.Error("ERROR: " + err.Error())
		return nil, err
	}

	// Otherwise, parse the health object.
	health = &management.Health{}

	if goErr := json.Unmarshal(body, health); goErr != nil {
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          fmt.Sprintf("go_error"),
			Description: goErr.Error(),
		}
	}

	return health, nil
}

// Creates a new tag in the dakota inventory.
func (client *Client) CreateDevice(accessToken string, device *model.Device) (*model.Device, *management.Error) {
	mlog.Debug("CreateDevice - accessToken: %s", accessToken)

	// Create device
	url := client.baseUrl + "/up/" + model.VERSION + "/devices"
	body, err := client.post(url, accessToken, device)

	// If error, return.
	if err != nil {
		return nil, err
	}

	// Otherwise, return the created device
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

// Start deactivation of a tag in the dakota inventory.
func (client *Client) DeactivateDevice(accessToken string, imsi uint64) *management.Error {
	mlog.Debug("DeactivateDevice - accessToken: %s", accessToken)

	// Deactivate device
	url := client.baseUrl + "/up/" + model.VERSION + "/devices/" + fmt.Sprintf("%d", imsi) + "/actions/deactivate"
	_, mErr := client.post(url, accessToken, nil)
	return mErr
}

// Deletes the tag with the specified imsi.
func (client *Client) DeleteDevice(accessToken string, imsi int) *management.Error {
	mlog.Debug("DeleteDevice - accessToken: %s", accessToken)

	// Delete device
	url := client.baseUrl + "/up/" + model.VERSION + "/devices/" + fmt.Sprintf("%d", imsi)
	return client.delete(url, accessToken)
}

// TODO - Add additional TAG methods in here.

// Create a new model in the dakota inventory.
func (client *Client) CreateModel(accessToken string, tagModel *model.Model) (*model.Model, *management.Error) {
	mlog.Debug("CreateModel - accessToken: %s", accessToken)

	// Create tag model.
	url := client.baseUrl + "/up/" + model.VERSION + "/models"
	body, err := client.post(url, accessToken, tagModel)

	// If error, return.
	if err != nil {
		return nil, err
	}

	// Otherwise, return the created tag model.
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

func (client *Client) GetDeviceByImsi(accessToken string, deviceId uint64) (*model.Device, *management.Error) {
	mlog.Debug("GetDeviceByImsi - accessToken: %s", accessToken)

	// Get device
	url := client.baseUrl + "/up/" + model.VERSION +  "/devices/" + fmt.Sprintf("%d", deviceId)

	data, serviceError := client.get(url, accessToken, nil)

	if serviceError != nil {
		mlog.Error("Error: %+v", serviceError)
		return nil, serviceError
	}

	returnedDevice := &model.Device{}

	mlog.Debug("Ready to unmarshall data %+v", data)
	if goErr := json.Unmarshal(data, returnedDevice); goErr != nil {
		mlog.Error("Error: %+v", goErr)
		return nil, &management.Error{
			HttpStatus:  http.StatusInternalServerError,
			Id:          "go_error",
			Description: goErr.Error(),
		}
	}

	mlog.Debug("Ready to return device %+v", returnedDevice)
	return returnedDevice, nil
}

//////////////////////////////////////////////////////////////////////////////////
// Private Methods

// Returns the body returned buy the GET on the specified url.
func (client *Client) get(url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	mlog.Debug("GET %s", url)

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
		mlog.Debug("URL: %s", url)
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

		mlog.Error("Error, GET failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Get successful resonse body.
	defer response.Body.Close()
	body, _ = ioutil.ReadAll(response.Body)

	// Collect some metrics
	latency := time.Since(startTime)
	mlog.Debug("GET %s Latency: %s", url, latency.String())

	return body, nil
}

// Returns the body returned after updation the specified resource object on the specified url.
func (client *Client) post(url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	mlog.Debug("POST %s", url)

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

		mlog.Error("Error, POST failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Parse the body if and only if one was provided.
	if response.Body != nil {
		defer response.Body.Close()
		body, _ = ioutil.ReadAll(response.Body)

		// Collect some metrics
		latency := time.Since(startTime)
		mlog.Debug("POST %s Latency: %s", url, latency.String())

		return body, nil
	}

	// Otherwisem we got a sucessfull response with no body.
	return nil, nil
}

// Returns the body returned after updation the specified resource object on the specified url.
func (client *Client) put(url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	mlog.Debug("PUT %s", url)

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
		mlog.Debug("Setting request authorization token.")
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

		mlog.Error("Error, PUT failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Parse the body if and only if one was provided.
	if response.Body != nil {
		defer response.Body.Close()
		body, _ = ioutil.ReadAll(response.Body)

		// Collect some metrics
		latency := time.Since(startTime)
		mlog.Debug("PUT %s Latency: %s", url, latency.String())

		return body, nil
	}

	// Otherwisem we got a sucessfull response with no body.
	return nil, nil
}

// Deletes the resource object associated with the specified url.
func (client *Client) delete(url string, accessToken string) (err *management.Error) {
	mlog.Debug("DELETE %s", url)

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

		mlog.Error("Error, DELETE failed with error: (%s) %s", err.Id, err.Description)

		return err
	}

	// Collect some metrics
	latency := time.Since(startTime)
	mlog.Debug("DELETE %s Latency: %s", url, latency.String())

	return nil
}

// Returns the body returned after updation the specified resource object on the specified url.
func (client *Client) patch(url string, accessToken string, resourceObject interface{}) (body []byte, err *management.Error) {
	mlog.Debug("PATCH %s", url)

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

	request, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
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

		mlog.Error("Error, PATCH failed with error: (%s) %s", err.Id, err.Description)

		return nil, err
	}

	// Parse the body if and only if one was provided.
	if response.Body != nil {
		defer response.Body.Close()
		body, _ = ioutil.ReadAll(response.Body)

		// Collect some metrics
		latency := time.Since(startTime)
		mlog.Debug("PATCH %s Latency: %s", url, latency.String())

		return body, nil
	}

	// Otherwise we got a sucessfull response with no body.
	return nil, nil
}
