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
	"encoding/json"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.Execution), NewExecutionTest)
}

type NotebookCRUDTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewExecutionTest() (tests.Test, error) {
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
	SubmitExecutionStep utils.Step = "Submitting execution request"
	ExecutionStatusStep utils.Step = "Retrieving execution status"
	ListExecutionStep   utils.Step = "List Executions"
)

func (test *NotebookCRUDTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- ExecutionTest")

	uuid, err := test.GetUUID()
	if err != nil {
		return logs.LogError("Failed to get callback ID with error: %s", err)
	}

	callbackURL := test.GetEventCallbackURL(uuid)

	writeChannel := make(chan json.RawMessage)
	callbacks.Set(uuid, writeChannel)
	logs.LogInfo("Registered callback URL: %s", callbackURL)

	executionRequest := &northstarApiModel.ExecutionRequest{
		Name:       "Test Cell 1",
		Language:   "lua",
		Arguments:  northstarApiModel.Arguments{},
		EntryPoint: "main",
		Code:       ``,
		Timeout:    90,
	}

	createdExecutionRequest, mErr := test.NorthstarApiClient.Execute(test.Users[0].Token.AccessToken, callbackURL, executionRequest)
	if mErr != nil {
		return logs.LogStep(SubmitExecutionStep, false, "Error, failed to submit execution request with error: %s", mErr.Description)
	}
	logs.LogStep(SubmitExecutionStep, true, "Created request: %+v", createdExecutionRequest)

	executionStatus, mErr := test.NorthstarApiClient.GetExecution(test.Users[0].Token.AccessToken, createdExecutionRequest.ExecutionId)
	if mErr != nil {
		return logs.LogStep(ExecutionStatusStep, false, "Failed to retrieve execution status. %s", mErr.Description)
	}
	logs.LogStep(ExecutionStatusStep, true, "Successfully retrieved execution status: %+v", executionStatus)

	executions, mErr := test.NorthstarApiClient.ListExecutions(test.Users[0].Token.AccessToken, 100)
	if mErr != nil {
		return logs.LogStep(ListExecutionStep, false, "Failed to list executions. %s", mErr.Description)
	}
	logs.LogStep(ListExecutionStep, true, "Executions: %+v", executions)

	return nil
}

func (test *NotebookCRUDTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	logs.LogInfo("No cleanup required for this test.")
	return nil
}
