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

package execution

import (
	"encoding/base64"
	"encoding/json"
	"time"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.MapExecution), NewNotebookMapExecutionTest)
}

type NotebookMapExecutionTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewNotebookMapExecutionTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &NotebookMapExecutionTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//Steps
	executeCellStep         utils.Step = "Execute map"
	createNotebookStep      utils.Step = "Create Notebook"
	executeNotebookCellStep utils.Step = "Execute notebook cell."

	//Common variables
	mapNotebookName = "Map Execution Example"
)

func (test *NotebookMapExecutionTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- MapExecutionTest")
	uuid, err := test.GetUUID()
	if err != nil {
		return logs.LogError("Failed to get callback ID with error: %s", err)
	}

	callbackURL := test.GetEventCallbackURL(uuid)

	writeChannel := make(chan json.RawMessage)
	callbacks.Set(uuid, writeChannel)
	logs.LogInfo("Registered callback URL: %s", callbackURL)

	cellUuid, err := test.GetUUID()
	if err != nil {
		return logs.LogStep(executeCellStep, false, "Failed to get cell ID with error: %s", err)
	}
	cell := northstarApiModel.Cell{
		Id:   cellUuid,
		Name: mapNotebookName,
		Input: northstarApiModel.Input{
			Type:       northstarApiModel.CodeCellType,
			Language:   "lua",
			EntryPoint: "main",
			Body: `local output = require("nsOutput")
			function main()
    			local map = {
        			type = "Map",
        			center = {latitude = 1, longitude = 1},
        			zoom = 8,
        			items = {
            			{
                			label = "My Bike",
                			locations = {
                    			{latitude = 0.4, longitude = 1.6},
                    			{latitude = 0.7, longitude = 1.3},
                    			{latitude = 1, longitude = 1},
                    			{latitude = 1.3, longitude = 0.7},
                    			{latitude = 1.6, longitude = 0.4}}
            			},
            			{
                			label = "My Car",
                			locations = {
                    			{latitude = 0.2, longitude = 1.6},
                    			{latitude = 0.5, longitude = 1.3},
                    			{latitude = 0.8, longitude = 1},
                    			{latitude = 1.1, longitude = 0.7},
                    			{latitude = 1.4, longitude = 0.4}}
            			}
        			}
    			}

    			local out, err = output.map(map)
    			if err ~= nil then
        			error(err)
    			end

    			return out
			end`,
		},
	}

	//Base64 encode the code section for execution
	cell.Input.Body = base64.StdEncoding.EncodeToString([]byte(cell.Input.Body))

	serviceErr := test.NorthstarApiClient.ExecuteCell(test.Users[0].Token.AccessToken, callbackURL, &cell)
	if serviceErr != nil {
		return logs.LogStep(executeCellStep, false, "Error executing cell. %s", serviceErr.Description)
	}

	select {
	case response := <-writeChannel:
		err, output := validateExecutionResult(logs, response, "application/vnd.vz.map")
		if err != nil {
			return logs.LogStep(executeCellStep, false, "Error executing cell. %s", err.Error())
		}
		logs.LogStep(executeCellStep, true, "Successful. Response: %+v", output)
	case <-time.After(config.Configuration.ExecutionResponseTimeout):
		return logs.LogStep(executeCellStep, false, "Error, no response received for cell execution. Timing out.")
	}

	notebook := northstarApiModel.Notebook{
		Name: mapNotebookName,
		Cells: []northstarApiModel.Cell{
			cell,
		},
	}

	createdNotebook, serviceErr := test.NorthstarApiClient.CreateNotebook(test.Users[0].Token.AccessToken, &notebook)
	if serviceErr != nil {
		return logs.LogStep(createNotebookStep, false, "Failed to create notebook. %s", serviceErr.Description)
	}
	defer test.DeleteNotebook(logs, test.Users[0].Token.AccessToken, createdNotebook.Id)
	logs.LogStep(createNotebookStep, true, "Success. Response: %+v", createdNotebook)

	//Sleep to let cassandra populate
	time.Sleep(config.DEFAULT_CASSANDRA_SLEEP)

	serviceErr = test.NorthstarApiClient.ExecuteNotebookCell(test.Users[0].Token.AccessToken, callbackURL, createdNotebook.Id, &cell)
	if serviceErr != nil {
		return logs.LogStep(executeNotebookCellStep, false, "Failed to execute notebook cell with error: %s", serviceErr.Description)
	}

	select {
	case response := <-writeChannel:
		err, output := validateExecutionResult(logs, response, "application/vnd.vz.map")
		if err != nil {
			return logs.LogStep(executeNotebookCellStep, false, "Failed to validate execution. %s", err.Error())
		}
		logs.LogStep(executeNotebookCellStep, true, "Response recieved: %+v", output)
	case <-time.After(config.Configuration.ExecutionResponseTimeout):
		return logs.LogStep(executeNotebookCellStep, false, "Error, no response received for notebook cell execution. Timing out.")
	}

	return nil
}

func (test *NotebookMapExecutionTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	notebooks, err := test.NorthstarApiClient.ListNotebooks(test.Users[0].Token.AccessToken)
	if err != nil {
		return logs.LogError("Cannot list notebooks. %s", err.Error())
	}

	for _, notebook := range notebooks {
		if notebook.Name == mapNotebookName {
			err := test.NorthstarApiClient.DeleteNotebook(test.Users[0].Token.AccessToken, notebook.Id)
			if err != nil {
				logs.LogError("Failed to delete notebook with ID: %s. %s", notebook.Id, err.Error())
			}
		}
	}

	return nil
}
