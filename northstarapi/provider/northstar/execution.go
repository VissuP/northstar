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
	"encoding/json"
	"fmt"
	"time"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/pkg/rte/events"
	"github.com/verizonlabs/northstar/pkg/rte/repl"

	"bytes"

	"code.cloudfoundry.org/bytefmt"
	invocations "github.com/verizonlabs/northstar/data/invocations/client"
	invocationModel "github.com/verizonlabs/northstar/data/invocations/model"
	snippets "github.com/verizonlabs/northstar/processing/snippets/client"
	snippetsModel "github.com/verizonlabs/northstar/processing/snippets/model"
)

// Defines the type used to represent an execution event.
type Execution struct {
	AccountId  string                 `json:"accountId"`
	Name       string                 `json:"name"`
	Language   string                 `json:"language"`
	EntryPoint string                 `json:"entryPoint"`
	Arguments  map[string]interface{} `json:"arguments"`
	Code       string                 `json:"code"`
	Url        string                 `json:"url"`
	Timeout    int                    `json:"timeout"`
	StartedOn  time.Time              `json:"startedOn"`
	Callback   string                 `json:"callback"`
	Memory     uint64                 `json:"memory,omitempty"`
}

// NorthStarExecutionProvider defines the type used to support code (.e.,g cell, transformation, etc.) executions.
type NorthStarExecutionProvider struct {
	snippetsClient    *snippets.SnippetsClient
	invocationsClient *invocations.InvocationClient
}

// NewNorthStarExecutionProvider returns a new NorthStar execution provider.
func NewNorthStarExecutionProvider() (*NorthStarExecutionProvider, error) {
	mlog.Info("NewNorthStarExecutionProvider")

	// Create data clients.
	snippetsClient, err := snippets.NewSnippetsClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create snippets client: %v", err)
	}

	invocationsClient, err := invocations.NewInvocationClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create invocations client: %v", err)
	}

	// Create the provider.
	provider := &NorthStarExecutionProvider{
		snippetsClient:    snippetsClient,
		invocationsClient: invocationsClient,
	}

	return provider, nil
}

//ExecutionCallback processes the response from an execution
func (provider *NorthStarExecutionProvider) ExecutionCallback(response *model.ExecutionResponse) {
	mlog.Info("ExecutionCallback")

	go provider.processExecutionCallback(response)

	return
}

func (provider *NorthStarExecutionProvider) processExecutionCallback(response *model.ExecutionResponse) {
	mlog.Debug("Processing response: %+v", response)

	output := &model.Output{}

	invocationData, mErr := provider.invocationsClient.GetInvocation(response.AccountID, response.InvocationID)
	if mErr != nil {
		mlog.Error("Get invocation data for id %s returned error: %v", response.InvocationID, mErr)
		output = &model.Output{
			Status:            model.OutputInternalErrorStatus,
			StatusDescription: "Could not retrieve execution results.",
		}
	}

	mlog.Debug("Serializing invocation data for id %s", response.InvocationID)
	output = fromExternalExecution(invocationData)

	//Note that we need to do manual encoding of the JSON because Go escapes HTML tags by default
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(output)
	if err != nil {
		mlog.Error("Error encoding execution output: %s", err.Error())

		//TODO: Fix this. Currently we drop it.
		output.Status = model.OutputInternalErrorStatus
		output.StatusDescription = "Could not serialize execution results."
		return
	}

	// Send results in callback. Note that our full URL is in the snippetOutput.Callback so there is no path
	if _, mErr := management.PostData(response.Callback, "", "application/json", buf.Bytes()); mErr != nil {
		mlog.Error("Send cell to callback url returned error: %v", mErr)
		return
	}
}

// Executes the specified cell.
func (provider *NorthStarExecutionProvider) ExecuteCell(user *model.User, cell *model.Cell, callback string) (*model.Cell, *management.Error) {
	mlog.Debug("ExecuteCell: cell:%v", cell)

	// Create execution. Note that we assume cells are always base64 encoded code.
	execution := model.ExecutionRequest{
		AccountId:  user.AccountId,
		Name:       cell.Name,
		Language:   cell.Input.Language,
		EntryPoint: cell.Input.EntryPoint,
		Arguments:  cell.Input.Arguments,
		Code:       cell.Input.Body,
		Timeout:    cell.Input.Timeout,
		Callback:   callback,
		Memory:     cell.Settings.Memory,
	}

	// Execute the cell.
	executionRequest, mErr := provider.Execute(user.AccountId, &execution)
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Execute returned error: %+v", mErr))
	}

	cell.ExecutionId = executionRequest.ExecutionId
	return cell, nil
}

// Executes the specified transformation.
func (provider *NorthStarExecutionProvider) ExecuteTransformation(user *model.User, transformation *model.Transformation, callback string) (*model.Transformation, *management.Error) {
	mlog.Debug("ExecuteTransformation: transformation:%v")

	// Create execution
	execution := model.ExecutionRequest{
		AccountId:  user.AccountId,
		Name:       transformation.Name,
		Language:   transformation.Language,
		EntryPoint: transformation.EntryPoint,
		Arguments:  transformation.Arguments,
		Code:       transformation.Code.Value,
		Timeout:    transformation.Timeout,
		Callback:   callback,
	}

	// Execute the cell.
	executionRequest, mErr := provider.Execute(user.AccountId, &execution)
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Execute returned error: %+v", mErr))
	}

	transformation.ExecutionId = executionRequest.ExecutionId
	return transformation, nil
}

// Execute is a helper method used to trigger executions.
func (provider *NorthStarExecutionProvider) Execute(accountId string, execution *model.ExecutionRequest) (*model.ExecutionRequest, *management.Error) {
	mlog.Debug("execute")

	options := snippetsModel.Options{
		Args:     execution.Arguments,
		Callback: execution.Callback,
		Memory:   execution.Memory * bytefmt.MEGABYTE,
	}
	// Create the snippet.
	snippet := &snippetsModel.Snippet{
		Runtime: execution.Language,
		MainFn:  execution.EntryPoint,
		Code:    execution.Code,
		URL:     "base64://code",
		Timeout: execution.Timeout * 1000, //Convert to milliseconds
		Options: options,
	}

	// Execute the snippet.
	invocationId, mErr := provider.snippetsClient.StartSnippet(accountId, snippet)
	if mErr != nil {
		return nil, mErr
	}

	mlog.Info("Submitted execution request: %s", invocationId)
	execution.ExecutionId = invocationId
	return execution, nil
}

//ListExecutions returns the list of all execution results that belong to the specified account.
func (provider *NorthStarExecutionProvider) ListExecutions(accountId string, limit int) ([]model.Output, *management.Error) {
	mlog.Debug("ListExecutions")

	invocations, mErr := provider.invocationsClient.GetInvocationsByAccountId(accountId, limit)
	if mErr != nil {
		return nil, mErr
	}

	results := []model.Output{}
	for _, invocation := range invocations {
		mlog.Info("Invocation result: %+v", invocation)
		results = append(results, *fromExternalExecution(invocation))
	}

	return results, nil
}

//GetExecution returns the specified execution
func (provider *NorthStarExecutionProvider) GetExecution(accountId string, executionId string) (*model.Output, *management.Error) {
	mlog.Debug("GetExecution")

	invocation, mErr := provider.invocationsClient.GetInvocation(accountId, executionId)
	if mErr != nil {
		return nil, mErr
	}
	result := fromExternalExecution(invocation)
	return result, nil
}

//StopExecution stops the specified execution
func (provider *NorthStarExecutionProvider) StopExecution(accountId string, executionId string) *management.Error {
	mlog.Debug("StopExecution")

	mErr := provider.snippetsClient.StopSnippet(accountId, executionId)
	if mErr != nil {
		return mErr
	}

	return nil
}

// Defines the map used to translate output status.
var outputStatusCodeMap = map[string]model.OutputStatusCode{
	repl.SNIPPET_RUN_FINISHED:    model.OutputSuccessStatus,
	repl.SNIPPET_REPL_FAILED:     model.OutputFailedStatus,
	repl.SNIPPET_CODE_GET_FAILED: model.CodeGetFailedStatus,
	repl.SNIPPET_RUN_TIMEDOUT:    model.OutputTimeoutStatus,
	repl.SNIPPET_OUT_OF_MEMORY:   model.OutputOutOfMemoryStatus,
	events.SNIPPET_START_EVENT:   model.OutputRunningStatus,
	events.SNIPPET_RUNNING_EVENT: model.OutputRunningStatus,
}

// Returns output status for the specified execution results.
func getOutputStatus(status, message string) (code model.OutputStatusCode, description string) {
	found := false

	// If status found, translate to code/description.
	if code, found = outputStatusCodeMap[status]; found {
		description = model.DefaultOutputStatusDescriptions[code]

		// If message provided, replace default description.
		if message != "" {
			description = message
		}

		return
	}

	// Else, use internal error as the default. E.g., all execution should have
	// a status. Otherwise, this is a coding error.
	code = model.OutputUnknownStatus
	description = model.DefaultOutputStatusDescriptions[code]
	return
}

//fromExternalExecution converts an execution result from the external format to the northstarapi format.
func fromExternalExecution(invocationData *invocationModel.InvocationData) *model.Output {
	mlog.Debug("fromExternalExecution")

	// Get the execution status.
	status, description := getOutputStatus(invocationData.Status, invocationData.ErrorDescr)

	output := model.Output{}
	output.Status = status
	output.StatusDescription = description
	output.Id = invocationData.Id

	// If the snippet status is successful, parse the results.
	if status == model.OutputSuccessStatus {
		// Set the execution output.
		output.ExecutionOutput = invocationData.Stdout
		output.ElapsedTime = int(invocationData.ElapsedTime)

		// Set the execution results (if any).
		if invocationData.Result != "" {
			var results model.CellResults

			if err := json.Unmarshal([]byte(invocationData.Result), &results); err != nil {
				mlog.Error("Unmarshal invocation data results for id %s returned error: %v", invocationData.Id, err)
				output.Status = "Could not parse execution results."
				output.StatusDescription = model.DefaultOutputStatusDescriptions[output.Status]
			} else {
				output.ExecutionResults = &results
			}
		}
	}
	return &output
}
