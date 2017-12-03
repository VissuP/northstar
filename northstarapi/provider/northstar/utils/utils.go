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

package utils

import (
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/pkg/rte/repl"
)

// Defines the map used to translate output status.
var outputStatusCodeMap = map[string]model.OutputStatusCode{
	repl.SNIPPET_RUN_FINISHED:    model.OutputSuccessStatus,
	repl.SNIPPET_REPL_FAILED:     model.OutputFailedStatus,
	repl.SNIPPET_CODE_GET_FAILED: model.OutputInternalErrorStatus,
	repl.SNIPPET_RUN_TIMEDOUT:    model.OutputTimeoutStatus,
}

// Returns output status for the specified execution results.
func GetOutputStatus(status, message string) (code model.OutputStatusCode, description string) {
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

	// Else, use unknown error as the default. E.g., all execution should have
	// a status. Otherwise, this is a coding error.
	code = model.OutputUnknownStatus
	description = model.DefaultOutputStatusDescriptions[code]
	return
}
