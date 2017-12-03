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

	"bytes"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/portal/model"
)

// ProcessEvent is a helper method used to parse async event payload based on type. Note
// that this method has two main motivations:
//	1. Keep external types (i.e., Portal API Model) encapsulated behind the Portal Provider.
//	2. Keep the portal decoupled from external types (i.e., Portal API Model). This will
//     reduce the impact of Portal API Model changes at the expense of the extra cost and code.
func (provider *NorthStarPortalProvider) ProcessEvent(id, payloadType string, payloadData []byte) (*model.Event, *management.Error) {
	mlog.Debug("ProcessEvent: payloadType:%s", payloadType)

	var eventType model.EventType
	var payload interface{}

	switch payloadType {
	case model.EventTypeExecuteResult.ToString():
		{
			mlog.Debug("Processing executing results.")
			externalOutput := northstarApiModel.Output{}

			// Unmarshal the payload data into portal api, external, output.
			if err := json.Unmarshal(payloadData, &externalOutput); err != nil {
				return nil, management.GetBadRequestError(fmt.Sprintf("Failed to umarshal event payload with error: %v", err))
			}

			// Create portal output.
			output := model.Output{
				State:         externalOutput.Status.ToString(),
				Stdout:        externalOutput.ExecutionOutput,
				Stderr:        externalOutput.StatusDescription,
				ElapsedTime:   externalOutput.ElapsedTime,
				LastExecution: externalOutput.LastExecution,
			}

			if externalOutput.ExecutionResults != nil {
				output.Results = &model.CellResults{
					Type:    externalOutput.ExecutionResults.Type.ToString(),
					Content: externalOutput.ExecutionResults.Content,
				}
			}

			payload = &output
			eventType = model.EventTypeExecuteResult
		}
	default:
		return nil, management.GetBadRequestError(fmt.Sprintf("The callback type %s is invalid.", payloadType))
	}

	// Marshal the types payload back into raw message.
	rawMessage := new(bytes.Buffer)
	encoder := json.NewEncoder(rawMessage)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to marshal event payload with error: %v", err))
	}

	event := &model.Event{
		Id:      id,
		Type:    eventType,
		Payload: rawMessage.Bytes(),
	}

	return event, nil
}
