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
	TransformationsPath = "transformations"
)

// CreateTransformation creates a new transformation object.
func (client *Client) CreateTransformation(accessToken string, transformation *model.Transformation) (*model.Transformation, *management.Error) {
	mlog.Debug("CreateTransformation")

	// Create transformation.
	path := client.getResourcePath(TransformationsPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.PostJSONWithHeaders(path, transformation, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	createdTransformation := &model.Transformation{}

	if err := json.Unmarshal(response, createdTransformation); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return createdTransformation, nil
}

// UpdateTransformation updates the specified transformation object,
func (client *Client) UpdateTransformation(accessToken string, transformation *model.Transformation) *management.Error {
	mlog.Debug("UpdateTransformation")

	// Update transformation.
	path := client.getResourcePath(TransformationsPath)
	headers := client.getRequestHeaders(accessToken)

	// If error, return.
	if _, mErr := client.lbClient.PutJSONWithHeaders(path, transformation, headers); mErr != nil {
		return mErr
	}

	return nil
}

// GetTransformationResults returns the execution ID for a transformation ID.
func (client *Client) GetTransformationResults(accessToken string, transformationID string) ([]model.Output, *management.Error) {
	mlog.Debug("GetTransformationResults")

	path := client.getResourcePath(TransformationsPath) + "/" + transformationID + "/results"
	headers := client.getRequestHeaders(accessToken)

	//If error, return
	response, mErr := client.lbClient.GetWithHeaders(path, headers)
	if mErr != nil {
		return nil, mErr
	}

	results := []model.Output{}
	if err := json.Unmarshal(response, &results); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return results, nil

}

// ListTransformations returns all transformation associated with the specidied access token.
func (client *Client) ListTransformations(accessToken string) ([]model.Transformation, *management.Error) {
	mlog.Debug("ListTransformations")

	// Update transformation.
	path := client.getResourcePath(TransformationsPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	transformations := make([]model.Transformation, 0)

	if err := json.Unmarshal(response, &transformations); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return transformations, nil
}

// GetTransformation gets the transformation with the specified id.
func (client *Client) GetTransformation(accessToken, transformationId string) (*model.Transformation, *management.Error) {
	mlog.Debug("GetTransformation")

	// Update transformation.
	path := client.getResourcePath(TransformationsPath) + "/" + transformationId
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	transformation := &model.Transformation{}

	if err := json.Unmarshal(response, transformation); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return transformation, nil
}

// DeleteTransformation deletes the transformation with the specified id.
func (client *Client) DeleteTransformation(accessToken, transformationId string) *management.Error {
	mlog.Debug("DeleteTransformation")

	// Update transformation.
	path := client.getResourcePath(TransformationsPath) + "/" + transformationId
	headers := client.getRequestHeaders(accessToken)

	return client.lbClient.DeleteWithHeaders(path, headers)
}

// CreateSchedule creates a new schedule for the specified transformation id.
func (client *Client) CreateSchedule(accessToken, transformationId string, schedule *model.Schedule) *management.Error {
	mlog.Debug("CreateSchedule")

	// Create schedule.
	path := client.getResourcePath(TransformationsPath) + "/" + transformationId + "/schedule"
	headers := client.getRequestHeaders(accessToken)

	// If error, return.
	if _, mErr := client.lbClient.PostJSONWithHeaders(path, schedule, headers); mErr != nil {
		return mErr
	}

	return nil
}

// GetSchedule returns the schedule associated with the specified transformation id.
func (client *Client) GetSchedule(accessToken, transformationId string) (*model.Schedule, *management.Error) {
	mlog.Debug("GetSchedule")

	// Get schedule.
	path := client.getResourcePath(TransformationsPath) + "/" + transformationId + "/schedule"
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	schedule := &model.Schedule{}

	if err := json.Unmarshal(response, schedule); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return schedule, nil
}

// DeleteSchedule deletes the schedule associated with the specified transformation id.
func (client *Client) DeleteSchedule(accessToken, transformationId string) *management.Error {
	mlog.Debug("DeleteSchedule")

	// Delete schedule.
	path := client.getResourcePath(TransformationsPath) + "/" + transformationId + "/schedule"
	headers := client.getRequestHeaders(accessToken)

	// If error, return.
	if mErr := client.lbClient.DeleteWithHeaders(path, headers); mErr != nil {
		return mErr
	}

	return nil
}
