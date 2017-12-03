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

package northstar

import (
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/portal/model"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
)

// CreateNotebook creates a new notebook.
func (provider *NorthStarPortalProvider) CreateNotebook(token string, notebook *model.Notebook) (*model.Notebook, *management.Error) {
	mlog.Debug("CreateNotebook")

	// Create the portal api (external) notebook.
	externalNotebook, mErr := provider.northstarApiClient.CreateNotebook(token, provider.toExternalNotebook(notebook))

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Create notebook returned error: %v", mErr))
	}

	// Update portal notebook fields.
	notebook.Id = externalNotebook.Id
	notebook.Version = externalNotebook.Version

	return notebook, nil
}

//UpdateNotebook updates an existing notebook.
func (provider *NorthStarPortalProvider) UpdateNotebook(token string, notebook *model.Notebook) *management.Error {
	mlog.Debug("UpdateNotebook")

	// Update the portal api (external) notebook.
	if mErr := provider.northstarApiClient.UpdateNotebook(token, provider.toExternalNotebook(notebook)); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Update notebook returned error: %v", mErr))
	}

	return nil
}

// ListNotebooks returns all notebooks associated with authenticated user. Note that this method only
// returns id and name of the notebook.
func (provider *NorthStarPortalProvider) ListNotebooks(token string) ([]model.Notebook, *management.Error) {
	mlog.Debug("ListNotebooks")

	// Get the portal api (external) notebooks for user (i.e. from token).
	externalNotebooks, mErr := provider.northstarApiClient.ListNotebooks(token)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("List notebooks returned error: %v", mErr))
	}

	// Translate to portal notebooks.
	var notebooks []model.Notebook

	for _, externalNotebook := range externalNotebooks {
		notebook := model.Notebook{
			Id:          externalNotebook.Id,
			Name:        externalNotebook.Name,
			Permissions: externalNotebook.Permissions,
		}

		notebooks = append(notebooks, notebook)
	}

	return notebooks, nil
}

// GetNotebook returns notebook with the specified id.
func (provider *NorthStarPortalProvider) GetNotebook(token string, notebookId string) (*model.Notebook, *management.Error) {
	mlog.Debug("GetNotebook")

	// Get the portal api (external) notebook with id.
	externalNotebook, mErr := provider.northstarApiClient.GetNotebook(token, notebookId)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get notebook with id %s returned error: %v", notebookId, mErr))
	}

	return fromExternalNotebook(externalNotebook), nil
}

// DeleteNotebook deletes the notbook with the specified id.
func (provider *NorthStarPortalProvider) DeleteNotebook(token string, notebookId string) *management.Error {
	mlog.Debug("DeleteNotebook")

	if mErr := provider.northstarApiClient.DeleteNotebook(token, notebookId); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete notebook with id %s returned error: %v", notebookId, mErr))
	}

	return nil
}

// GetNotebookUsers returns the permissions for the users of a notebook
func (provider *NorthStarPortalProvider) GetNotebookUsers(token string, notebookId string) ([]model.User, *management.Error) {
	mlog.Debug("GetNotebookUsers")

	externalUsers, serviceErr := provider.northstarApiClient.GetNotebookUsers(token, notebookId)
	if serviceErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("failed to get notebook users with error: %s", serviceErr.Description))
	}

	var users []model.User
	for _, externalUser := range externalUsers {
		user := provider.FromExternalUser(&externalUser)
		users = append(users, *user)
	}

	return users, nil
}

func (provider *NorthStarPortalProvider) UpdateNotebookUsers(token string, notebookId string, users []model.User) *management.Error {
	mlog.Debug("UpdateNotebookUsers")

	//TODO: Is this wrong? i.e. externalUsers :=
	var externalUsers []northstarApiModel.User
	for _, user := range users {
		externalUser := provider.ToExternalUser(&user)
		externalUsers = append(externalUsers, *externalUser)
	}

	serviceErr := provider.northstarApiClient.UpdateNotebookUsers(token, notebookId, externalUsers)
	if serviceErr != nil {
		return management.GetExternalError(fmt.Sprintf("error updating users: %s ", serviceErr.Description))
	}
	return nil
}

// ExecuteCell executes the input of the specified cell. Note that results will be returned, asynchronously, through the
// callback url.
func (provider *NorthStarPortalProvider) ExecuteCell(token string, callbackUrl string, cell *model.Cell) *management.Error {
	mlog.Debug("ExecuteCell: callbackUrl:%s, cell:%+v", callbackUrl, cell)

	// Get the portal api (external) cell.
	externalCell := toExternalCell(cell)

	if mErr := provider.northstarApiClient.ExecuteCell(token, callbackUrl, externalCell); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Execute cell returned error: %v", mErr))
	}

	return nil
}

// Helper method used to translate portal model to portal api model.
func (provider *NorthStarPortalProvider) toExternalNotebook(notebook *model.Notebook) *northstarApiModel.Notebook {
	mlog.Debug("ToExternalNotebook")

	// Create portal api (external) model notebook.
	externalNotebook := &northstarApiModel.Notebook{
		Id:   notebook.Id,
		Name: notebook.Name,
	}

	// Create the portal api (external) model cells.
	for _, cell := range notebook.Cells {
		externalCell := toExternalCell(&cell)
		externalNotebook.Cells = append(externalNotebook.Cells, *externalCell)
	}

	return externalNotebook
}

// toExternalCell is a helper method used to translate portal model to portal api model.
func toExternalCell(cell *model.Cell) *northstarApiModel.Cell {
	mlog.Debug("toExternalCell")

	// TODO(s)
	// - For now code cell only. E.g., R, Lua.
	//   This needs to be updated to support, potentially,
	//   queries.
	// - Need to pass (optional) arguments from the client side.
	var outputStatusCode northstarApiModel.OutputStatusCode

	// Create portal api (external) model cell.
	externalCell := &northstarApiModel.Cell{
		Id:   cell.ID,
		Name: cell.Name,
		Input: northstarApiModel.Input{
			Language:   cell.Language,
			EntryPoint: cell.Options.MainFunction,
			Body:       cell.Code,
			Timeout:    cell.Options.Timeout,
			Arguments:  cell.Inputs,
		},
		Output: northstarApiModel.Output{
			Status:            outputStatusCode.FromString(cell.Output.State),
			ExecutionOutput:   cell.Output.Stdout,
			StatusDescription: cell.Output.Stderr,
			ElapsedTime:       cell.Output.ElapsedTime,
			LastExecution:     cell.Output.LastExecution,
		},
		Settings: northstarApiModel.Settings{
			Status:            cell.Options.Status,
			ShowCode:          cell.Options.ShowCode,
			ShowOutput:        cell.Options.ShowOutput,
			ShowConfiguration: cell.Options.ShowConfiguration,
			Visualization: northstarApiModel.Visualization{
				Parameters: cell.Options.Visualization.Parameters,
			},
		},
	}

	if cell.Output.Results != nil {
		externalCell.Output.ExecutionResults = &northstarApiModel.CellResults{
			Type:    northstarApiModel.ResultTypes(cell.Output.Results.Type),
			Content: cell.Output.Results.Content,
		}
	}

	return externalCell
}

// fromExternalCell is a helper method used to translate portal model to portal api model.
func fromExternalCell(externalCell *northstarApiModel.Cell) *model.Cell {
	mlog.Debug("fromExternalCell")

	// Set template data to cell type.
	cell := &model.Cell{
		ID:       externalCell.Id,
		Name:     externalCell.Name,
		Language: externalCell.Input.Language,
		Code:     externalCell.Input.Body,
		Inputs:   externalCell.Input.Arguments,
		Output: model.Output{
			State:         externalCell.Output.Status.ToString(),
			Stdout:        externalCell.Output.ExecutionOutput,
			Stderr:        externalCell.Output.StatusDescription,
			ElapsedTime:   externalCell.Output.ElapsedTime,
			LastExecution: externalCell.Output.LastExecution,
		},
		Options: model.Options{
			Status:            externalCell.Settings.Status,
			MainFunction:      externalCell.Input.EntryPoint,
			Timeout:           externalCell.Input.Timeout,
			ShowCode:          externalCell.Settings.ShowCode,
			ShowOutput:        externalCell.Settings.ShowOutput,
			ShowConfiguration: externalCell.Settings.ShowConfiguration,
			Visualization: model.Visualization{
				Parameters: externalCell.Settings.Visualization.Parameters,
			},
		},
	}

	return cell
}

// Helper method used to translate portal api model to portal model.
func fromExternalNotebook(externalNotebook *northstarApiModel.Notebook) *model.Notebook {
	mlog.Debug("FromExternalNotebook")

	// Create portal model notebook.
	notebook := &model.Notebook{
		Id:          externalNotebook.Id,
		Version:     externalNotebook.Version,
		Name:        externalNotebook.Name,
		Permissions: externalNotebook.Permissions,
	}

	// Translate cells.
	for _, externalCell := range externalNotebook.Cells {
		cell := model.Cell{
			ID:       externalCell.Id,
			Name:     externalCell.Name,
			Language: externalCell.Input.Language,
			Code:     externalCell.Input.Body,
			Inputs:   externalCell.Input.Arguments,
			Output: model.Output{
				State:         externalCell.Output.Status.ToString(),
				Stdout:        externalCell.Output.ExecutionOutput,
				Stderr:        externalCell.Output.StatusDescription,
				ElapsedTime:   externalCell.Output.ElapsedTime,
				LastExecution: externalCell.Output.LastExecution,
			},
			Options: model.Options{
				Status:            externalCell.Settings.Status,
				MainFunction:      externalCell.Input.EntryPoint,
				Timeout:           externalCell.Input.Timeout,
				ShowCode:          externalCell.Settings.ShowCode,
				ShowOutput:        externalCell.Settings.ShowOutput,
				ShowConfiguration: externalCell.Settings.ShowConfiguration,
				Visualization: model.Visualization{
					Parameters: externalCell.Settings.Visualization.Parameters,
				},
			},
		}

		if externalCell.Output.ExecutionResults != nil {
			cell.Output.Results = &model.CellResults{
				Type:    string(externalCell.Output.ExecutionResults.Type),
				Content: externalCell.Output.ExecutionResults.Content,
			}
		}

		notebook.Cells = append(notebook.Cells, cell)
	}

	return notebook
}
