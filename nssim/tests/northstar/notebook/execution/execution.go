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
	"github.com/verizonlabs/northstar/nssim/utils"
)

//validateExecutionResult validates the execution result
func validateExecutionResult(logs *utils.Logger, response json.RawMessage, expectedType string) (error, northstarApiModel.Output) {
	output := northstarApiModel.Output{}
	err := json.Unmarshal([]byte(response), &output)
	if err != nil {
		return err, output
	}

	logs.LogInfo("Execution output: %+v results: %+v", output, output.ExecutionResults)
	if output.Status != northstarApiModel.OutputSuccessStatus {
		return logs.LogError("Failed to execute with error: %s.", output.StatusDescription), output
	}

	if output.ExecutionResults == nil {
		return logs.LogError("Error, invalid execution results returned: %+v", output.ExecutionResults), output
	}
	if string(output.ExecutionResults.Type) != expectedType {
		return logs.LogError("Wrong type (%s) returned from table execution. Expected %s", output.ExecutionResults.Type, expectedType), output
	}

	return nil, output
}
