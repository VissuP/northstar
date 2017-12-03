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

package notebook

import (
	"time"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.Notebook), NewNotebookCRUDTest)
}

type NotebookCRUDTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewNotebookCRUDTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &NotebookCRUDTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//Steps

	//Common variables
	notebookName = "NSSim Crud Test Notebook"
)

func (test *NotebookCRUDTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NotebookCRUDTest")

	//create a notebook
	newNotebook := &northstarApiModel.Notebook{
		Name:  notebookName,
		Cells: nil,
	}

	var createNotebookStep utils.Step = "Creating notebook"
	createdNotebook, mErr := test.NorthstarApiClient.CreateNotebook(test.Users[0].Token.AccessToken, newNotebook)
	if mErr != nil {
		return logs.LogStep(createNotebookStep, false, "Failed to create notebook with error: %s", mErr.Description)
	}
	defer test.DeleteNotebook(logs, test.Users[0].Token.AccessToken, createdNotebook.Id)

	//sleep to let cassandra populate
	time.Sleep(config.DEFAULT_CASSANDRA_SLEEP)

	//Verify notebook has an ID
	if createdNotebook.Id == "" {
		return logs.LogStep(createNotebookStep, false, "Notebook missing ID. Notebook: %+v", createdNotebook)
	}

	logs.LogStep(createNotebookStep, true, "Created notebook: %+v", createdNotebook)

	var retrieveNotebookStep utils.Step = "Retrieving notebook"
	retrievedNotebook, mErr := test.NorthstarApiClient.GetNotebook(test.Users[0].Token.AccessToken, createdNotebook.Id)
	if mErr != nil {
		return logs.LogStep(retrieveNotebookStep, false, "Failed to retrieve notebook with error: %s", mErr.Description)
	}
	logs.LogStep(retrieveNotebookStep, true, "Successfully retrieved notebook: %+v", retrievedNotebook)

	var retrieveNotebookUnauthorizedStep utils.Step = "Attempting to retrieve notebook as an unauthorized user"
	_, mErr = test.NorthstarApiClient.GetNotebook(test.Users[1].Token.AccessToken, createdNotebook.Id)
	if mErr == nil {
		return logs.LogStep(retrieveNotebookUnauthorizedStep, false, "Error, user allowed to access notebook without permission.")
	}
	logs.LogStep(retrieveNotebookUnauthorizedStep, true, "Successfully prevented unauthorized access..")

	var listNotebooksStep utils.Step = "Attempting to list notebooks."
	retrievedNotebooks, mErr := test.NorthstarApiClient.ListNotebooks(test.Users[0].Token.AccessToken)
	if mErr != nil {
		return logs.LogStep(listNotebooksStep, false, "Failed to list notebooks with error: %s", mErr.Description)
	}

	//We created a notebook, so at least one should exist. More may exist depending on other simulataneous tests.
	if len(retrievedNotebooks) < 1 {
		return logs.LogStep(listNotebooksStep, false, "Error, %d notebooks were found. Expected this number to be greater than 1.", len(retrievedNotebooks))
	}

	logs.LogStep(listNotebooksStep, true, "Listed notebooks successfully.")

	//update the notebook
	var updateNotebookStep utils.Step = "Update Notebook"
	createdNotebook.Name = notebookName
	createdNotebook.Cells = []northstarApiModel.Cell{
		{
			Name: "Test Cell 1",
			Input: northstarApiModel.Input{
				Type:       northstarApiModel.CodeCellType,
				Language:   "Lua",
				Arguments:  northstarApiModel.Arguments{},
				EntryPoint: "main",
				Body: `local output = require("nsOutput")
							function main()
    								local table = {
        								columns = {"column1", "column2"},
        								rows = {{1, 2}, {3, 4}, {5, 6}}
    								}

    								local out, err = output.table(table)
    								if err ~= nil then
        								error(err)
    								end

    								return out
							end`,
				Timeout: 90,
			},
			Output:   northstarApiModel.Output{},
			Settings: northstarApiModel.Settings{},
		},
	}

	logs.LogInfo("Updating notebook")
	//Update as a user with permissions
	mErr = test.NorthstarApiClient.UpdateNotebook(test.Users[0].Token.AccessToken, createdNotebook)
	if mErr != nil {
		return logs.LogStep(updateNotebookStep, false, "Error, failed to update notebook with error: %s", mErr.Description)
	}

	retrievedUpdatedNotebook, mErr := test.NorthstarApiClient.GetNotebook(test.Users[0].Token.AccessToken, createdNotebook.Id)
	if mErr != nil {
		return logs.LogStep(updateNotebookStep, false, "Failed to retrieve notebook with error: %s", mErr.Description)
	}
	logs.LogStep(updateNotebookStep, true, "Successfully updated notebook: %+v", retrievedUpdatedNotebook)

	var updateNotebookUnauthorizedStep utils.Step = "Attempting to update notebook as unauthorized user."
	//Now try as a user without permission
	mErr = test.NorthstarApiClient.UpdateNotebook(test.Users[1].Token.AccessToken, createdNotebook)
	if mErr == nil {
		return logs.LogStep(updateNotebookUnauthorizedStep, false, "Error, unauthorized user was able to update notebook")
	}
	logs.LogStep(updateNotebookUnauthorizedStep, true, "Successfully prevented unauthorized access.")

	var deleteNotebookStep utils.Step = "Delete notebook as unauthorized user."
	//delete the notebook as an unauthorized user
	mErr = test.NorthstarApiClient.DeleteNotebook(test.Users[1].Token.AccessToken, createdNotebook.Id)
	if mErr == nil {
		return logs.LogStep(deleteNotebookStep, false, "Error, unauthorized user was able to delete notebook")
	}
	logs.LogStep(deleteNotebookStep, true, "Successfully prevented unauthorized access.")

	return nil
}

func (test *NotebookCRUDTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	notebooks, err := test.NorthstarApiClient.ListNotebooks(test.Users[0].Token.AccessToken)
	if err != nil {
		logs.LogError("Cannot list notebooks. %s", err.Error())
	}

	for _, notebook := range notebooks {
		if notebook.Name == notebookName {
			err := test.NorthstarApiClient.DeleteNotebook(test.Users[0].Token.AccessToken, notebook.Id)
			if err != nil {
				logs.LogError("Failed to delete notebook with ID: %s. %s", notebook.Id, err.Error())
			}
		}
	}

	return nil
}
