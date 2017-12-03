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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

const (
	ExecutionsPath    = "executions"
	CallbackUrlHeader = "x-vz-callback-url"
)

// ExecuteCell executes the specified cell returning the output through the callback URL.
func (client *Client) ExecuteCell(accessToken, callbackUrl string, cell *model.Cell) *management.Error {
	mlog.Debug("ExecuteCell: callbackUrl:%s", callbackUrl)

	// Execute cell.
	path := client.getResourcePath(ExecutionsPath) + "/cell"
	headers := client.getRequestHeaders(accessToken)
	headers[CallbackUrlHeader] = callbackUrl

	// If error, return.
	if _, mErr := client.lbClient.PostJSONWithHeaders(path, cell, headers); mErr != nil {
		return mErr
	}

	return nil
}

// ExecuteTransformation executes the specified transformation returning the output through the callback URL.
func (client *Client) ExecuteTransformation(accessToken, callbackUrl string, transformation *model.Transformation) *management.Error {
	mlog.Debug("ExecuteTransformation: callbackUrl:%s", callbackUrl)

	// Execute transformation.
	path := client.getResourcePath(ExecutionsPath) + "/transformation"
	headers := client.getRequestHeaders(accessToken)
	headers[CallbackUrlHeader] = callbackUrl

	// If error, return.
	if _, mErr := client.lbClient.PostJSONWithHeaders(path, transformation, headers); mErr != nil {
		return mErr
	}

	return nil
}

// Execute executes the specified code, returning the output through the callback URL
func (client *Client) Execute(accessToken, callbackUrl string, executionRequest *model.ExecutionRequest) (*model.ExecutionRequest, *management.Error) {
	mlog.Debug("Execute: callbackURL: %s", callbackUrl)

	//Submit execution request
	path := client.getResourcePath(ExecutionsPath) + "/generic"
	headers := client.getRequestHeaders(accessToken)
	headers[CallbackUrlHeader] = callbackUrl

	//If errors, return
	response, mErr := client.lbClient.PostJSONWithHeaders(path, executionRequest, headers)

	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	createdExecution := &model.ExecutionRequest{}

	if err := json.Unmarshal(response, createdExecution); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return createdExecution, nil
}

//GetExecution returns the results/status of an execution
func (client *Client) GetExecution(accessToken string, executionId string) (*model.ExecutionResponse, *management.Error) {
	mlog.Debug("GetExecution")

	path := client.getResourcePath(ExecutionsPath) + "/generic/" + executionId
	headers := client.getRequestHeaders(accessToken)

	//If errors, return
	response, mErr := client.lbClient.GetWithHeaders(path, headers)
	if mErr != nil {
		return nil, mErr
	}

	executionResponse := &model.ExecutionResponse{}
	if err := json.Unmarshal(response, executionResponse); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return executionResponse, nil
}

//StopExecution stops the specified execution
func (client *Client) StopExecution(accessToken string, executionId string) *management.Error {
	mlog.Debug("StopExecution")

	path := client.getResourcePath(ExecutionsPath) + "/generic/" + executionId + "/actions/stop"
	headers := client.getRequestHeaders(accessToken)

	//If errors, return
	_, mErr := client.lbClient.PostJSONWithHeaders(path, nil, headers)
	if mErr != nil {
		return mErr
	}

	return nil
}

//ListExecutions lists the executions belonging to this account
func (client *Client) ListExecutions(accessToken string, limit int) ([]model.ExecutionResponse, *management.Error) {
	mlog.Debug("ListExecutions")

	path := client.getResourcePath(ExecutionsPath)
	headers := client.getRequestHeaders(accessToken)
	headers["Limit"] = fmt.Sprintf("%d", limit)

	//If errors, return
	response, mErr := client.lbClient.GetWithHeaders(path, headers)
	if mErr != nil {
		return nil, mErr
	}

	results := []model.ExecutionResponse{}
	if err := json.Unmarshal(response, results); err != nil {
		return nil, mErr
	}

	return results, nil

}
