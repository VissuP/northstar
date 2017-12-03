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
	NotebooksPath = "notebooks"
)

// CreateNotebook creates a new notebook object.
func (client *Client) CreateNotebook(accessToken string,
	notebook *model.Notebook) (*model.Notebook, *management.Error) {
	mlog.Debug("CreateNotebook")

	// Create notebook.
	path := client.getResourcePath(NotebooksPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.PostJSONWithHeaders(path, notebook, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	createdNotebook := &model.Notebook{}

	if err := json.Unmarshal(response, createdNotebook); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return createdNotebook, nil
}

// UpdateNotebook updates the specified notebook object,
func (client *Client) UpdateNotebook(accessToken string, notebook *model.Notebook) *management.Error {
	mlog.Debug("UpdateNotebook")

	// Update notebook.
	path := client.getResourcePath(NotebooksPath)
	headers := client.getRequestHeaders(accessToken)

	// If error, return.
	if _, mErr := client.lbClient.PutJSONWithHeaders(path, notebook, headers); mErr != nil {
		return mErr
	}

	return nil
}

// ListNotebooks returns all notebook associated with the specidied access token.
func (client *Client) ListNotebooks(accessToken string) ([]model.Notebook, *management.Error) {
	mlog.Debug("ListNotebooks")

	// Update notebook.
	path := client.getResourcePath(NotebooksPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resources.
	notebooks := make([]model.Notebook, 0)

	if err := json.Unmarshal(response, &notebooks); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return notebooks, nil
}

// GetNotebook gets the notebook with the specified id.
func (client *Client) GetNotebook(accessToken, notebookId string) (*model.Notebook, *management.Error) {
	mlog.Debug("GetNotebook")

	// Update notebook.
	path := client.getResourcePath(NotebooksPath) + "/" + notebookId
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	notebook := &model.Notebook{}

	if err := json.Unmarshal(response, notebook); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return notebook, nil
}

// DeleteNotebook deletes the notebook with the specified id.
func (client *Client) DeleteNotebook(accessToken, notebookId string) *management.Error {
	mlog.Debug("DeleteNotebook")

	// Update notebook.
	path := client.getResourcePath(NotebooksPath) + "/" + notebookId
	headers := client.getRequestHeaders(accessToken)

	return client.lbClient.DeleteWithHeaders(path, headers)
}

// GetNotebookUsers returns notebook users.
func (client *Client) GetNotebookUsers(accessToken, notebookId string) ([]model.User, *management.Error) {
	mlog.Debug("GetNotebookUsers")

	// Get notebook users.
	path := client.getResourcePath(NotebooksPath) + "/" + notebookId + "/users"
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	users := make([]model.User, 0)

	if err := json.Unmarshal(response, &users); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return users, nil
}

// UpdateNotebookUsers updates notebook users.
func (client *Client) UpdateNotebookUsers(accessToken, notebookId string, users []model.User) *management.Error {
	mlog.Debug("UpdateNotebookUsers")

	// Update notebook users.
	path := client.getResourcePath(NotebooksPath) + "/" + notebookId + "/users"
	headers := client.getRequestHeaders(accessToken)

	_, mErr := client.lbClient.PutJSONWithHeaders(path, users, headers)

	// If error, return.
	if mErr != nil {
		return mErr
	}

	return nil

}

// ExecuteNotebookCell executes the specified notebook cell returning the output through the callback URL.
func (client *Client) ExecuteNotebookCell(accessToken, callbackUrl string, notebookId string, cell *model.Cell) *management.Error {
	mlog.Debug("ExecuteCell: callbackUrl:%s", callbackUrl)

	// Execute cell.
	path := client.getResourcePath(NotebooksPath) + "/" + notebookId + "/cells/" + cell.Id + "/actions/execute"
	headers := client.getRequestHeaders(accessToken)
	headers[CallbackUrlHeader] = callbackUrl

	// If error, return.
	if _, mErr := client.lbClient.PostJSONWithHeaders(path, cell, headers); mErr != nil {
		return mErr
	}

	return nil
}
