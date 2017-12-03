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
	tests.Register(tests.TestId(config.ValueExecution), NewNotebookValueExecutionTest)
}

type NotebookValueExecutionTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewNotebookValueExecutionTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &NotebookValueExecutionTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//Steps
	executeValueStep utils.Step = "Execute value."

	//Common variables
	valueNotebookName = "NSSim Value Test notebook"
)

func (test *NotebookValueExecutionTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- ValueExecutionTest")
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
		return logs.LogStep(executeValueStep, false, "Failed to get cell ID with error: %s", err)
	}
	cell := northstarApiModel.Cell{
		Id:   cellUuid,
		Name: "NSSim Value Execution Example",
		Input: northstarApiModel.Input{
			Type:       northstarApiModel.CodeCellType,
			Language:   "lua",
			EntryPoint: "main",
			Body: `local output = require("nsOutput")
				function main()
    				local value = {
        				type = "int",
        				value = "5"
    				}

    				local out, err = output.value(value)
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
		return logs.LogStep(executeValueStep, false, "Error executing cell. %s", serviceErr.Description)
	}

	select {
	case response := <-writeChannel:
		err, output := validateExecutionResult(logs, response, "application/vnd.vz.value")
		if err != nil {
			return logs.LogStep(executeValueStep, false, err.Error())
		}
		logs.LogStep(executeValueStep, true, "Success. Response: %+v", output)
	case <-time.After(config.Configuration.ExecutionResponseTimeout):
		return logs.LogStep(executeValueStep, false, "Error, no response received for cell execution. Timing out.")
	}

	notebook := northstarApiModel.Notebook{
		Name: valueNotebookName,
		Cells: []northstarApiModel.Cell{
			cell,
		},
	}

	createdNotebook, serviceErr := test.NorthstarApiClient.CreateNotebook(test.Users[0].Token.AccessToken, &notebook)
	if serviceErr != nil {
		return logs.LogStep(createNotebookStep, false, "Failed to create notebook with error: %s", serviceErr.Description)
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
		err, output := validateExecutionResult(logs, response, "application/vnd.vz.value")
		if err != nil {
			return logs.LogStep(executeNotebookCellStep, false, err.Error())
		}
		logs.LogStep(executeNotebookCellStep, true, "Response received: %+v", output)
	case <-time.After(config.Configuration.ExecutionResponseTimeout):
		return logs.LogStep(executeNotebookCellStep, false, "Error, no response received for notebook cell execution. Timing out.")
	}

	return nil
}

func (test *NotebookValueExecutionTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	notebooks, err := test.NorthstarApiClient.ListNotebooks(test.Users[0].Token.AccessToken)
	if err != nil {
		return logs.LogError("Cannot list notebooks. %s", err.Error())
	}

	for _, notebook := range notebooks {
		if notebook.Name == valueNotebookName {
			err := test.NorthstarApiClient.DeleteNotebook(test.Users[0].Token.AccessToken, notebook.Id)
			if err != nil {
				logs.LogError("Failed to delete notebook with ID: %s. %s", notebook.Id, err.Error())
			}
		}
	}

	return nil
}
