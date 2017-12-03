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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/notebooks/model"
	"github.com/verizonlabs/northstar/data/util"
)

const NotebooksURI = util.DataBasePath + "/notebooks"

// Creates a new notebook.
func (client *NotebooksClient) CreateNotebook(notebook *model.Notebook) (*model.Notebook, *management.Error) {
	mlog.Debug("CreateNotebook")

	response, mErr := client.lbClient.PostJSON(NotebooksURI, notebook)

	if mErr != nil {
		return nil, mErr
	}

	createdNotebook := &model.Notebook{}

	if err := json.Unmarshal(response, createdNotebook); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return createdNotebook, nil
}

// Updates an existing notebook.
func (client *NotebooksClient) UpdateNotebook(notebookId string, notebook *model.Notebook) (*model.Notebook, *management.Error) {
	mlog.Debug("UpdateNotebook")

	path := NotebooksURI + "/" + notebookId
	response, mErr := client.lbClient.PutJSON(path, notebook)

	if mErr != nil {
		return nil, mErr
	}

	updatedNotebook := &model.Notebook{}

	if err := json.Unmarshal(response, updatedNotebook); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return updatedNotebook, nil
}

// Returns notebook for the specified id.
func (client *NotebooksClient) GetNotebook(notebookId string) (*model.Notebook, *management.Error) {
	mlog.Debug("GetNotebook")

	path := NotebooksURI + "/" + notebookId
	response, mErr := client.lbClient.Get(path)

	if mErr != nil {
		return nil, mErr
	}

	notebook := &model.Notebook{}

	if err := json.Unmarshal(response, notebook); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return notebook, nil
}

// Deletes a notebook.
func (client *NotebooksClient) DeleteNotebook(notebookId string) *management.Error {
	mlog.Debug("DeleteNotebook")

	path := NotebooksURI + "/" + notebookId

	if mErr := client.lbClient.Delete(path); mErr != nil {
		return mErr
	}

	return nil
}
