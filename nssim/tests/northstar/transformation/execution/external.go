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
	"time"

	mappingClient "github.com/verizonlabs/northstar/data/mappings/client"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
	eventsClient "github.com/verizonlabs/northstar/processing/events/client"
	eventsModel "github.com/verizonlabs/northstar/processing/events/model"
)

func init() {
	tests.Register(tests.TestId(config.ExternalEventExecution), NewTransformationExternalEventTest)
}

type TransformationExternalEvent struct {
	*northstar.NorthstarApiBaseTest
	eventsClient  *eventsClient.EventsClient
	mappingClient *mappingClient.MappingsClient
}

func NewTransformationExternalEventTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	mappingClient, err := mappingClient.NewMappingsClient()
	if err != nil {
		return nil, err
	}

	eventsClient, err := eventsClient.NewEventsClient()
	if err != nil {
		return nil, err
	}

	return &TransformationExternalEvent{
		NorthstarApiBaseTest: nsapiBase,
		mappingClient:        mappingClient,
		eventsClient:         eventsClient,
	}, nil
}

var (
	//steps
	scheduleTransformation utils.Step = "Schedule Transformation"
	getTransformationStep  utils.Step = "Get transformation"
	getMappingStep         utils.Step = "Get mapping."
	triggerEventStep       utils.Step = "Triggering event"
	getResultsStep         utils.Step = "Getting transformation results."

	//common variables
	externalTransformationName = "External Transformation Test"
)

func (test *TransformationExternalEvent) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- TransformationExternalEventTest")

	//Create transformation
	transformation := &northstarApiModel.Transformation{
		Name:        externalTransformationName,
		Description: "NS Sim Test Transformation",
		Timeout:     180,
		EntryPoint:  "main",
		Language:    "lua",
		Code: northstarApiModel.Code{
			Type: northstarApiModel.SourceCodeType,
			Value: `local output = require('nsOutput')

				function main()
				output.printf("Device event triggered. Event: %s ID: %s", context.Args['name'], context.Args['id'])
				end`,
		},
	}

	//Base64 encode the code section for execution
	transformation.Code.Value = base64.StdEncoding.EncodeToString([]byte(transformation.Code.Value))

	createdTransformation, serviceErr := test.NorthstarApiClient.CreateTransformation(test.Users[0].Token.AccessToken, transformation)
	if serviceErr != nil {
		return logs.LogStep(createTransformationStep, false, "Error, failed to create transformation with error: %s", serviceErr.Description)
	}
	defer test.DeleteTransformation(logs, test.Users[0].Token.AccessToken, createdTransformation.Id)
	logs.LogStep(createTransformationStep, true, "Successful. Transformation: %s", createdTransformation)

	//Schedule transformation
	schedule := northstarApiModel.Schedule{
		Event: northstarApiModel.Event{
			Category: northstarApiModel.DeviceEvent,
			Name:     "energyBatteryCharge",
		},
	}

	serviceErr = test.NorthstarApiClient.CreateSchedule(test.Users[0].Token.AccessToken, createdTransformation.Id, &schedule)
	if serviceErr != nil {
		return logs.LogStep(scheduleTransformation, false, "Error, failed to schedule transformation. Error was: %s", serviceErr.Description)
	}
	defer test.DeleteSchedule(logs, test.Users[0].Token.AccessToken, createdTransformation.Id)
	logs.LogStep(scheduleTransformation, true, "Success.")

	retrievedTransformation, serviceErr := test.NorthstarApiClient.GetTransformation(test.Users[0].Token.AccessToken, createdTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(getTransformationStep, false, "Error, failed to get transformation. Error was: %s", serviceErr.Description)
	}

	if retrievedTransformation.Schedule == nil {
		return logs.LogStep(getTransformationStep, false, "Error, retrieved transformation was missing a schedule.")
	}

	logs.LogStep(getTransformationStep, true, "Success.Transformation: %+v. Schedule: %+v", retrievedTransformation, retrievedTransformation.Schedule)

	//CreateSchedule scheduled our event. Lets get the mapping
	mapping, serviceErr := test.mappingClient.GetMapping(test.Users[0].AccountId, retrievedTransformation.Schedule.Id)
	if serviceErr != nil {
		return logs.LogStep(getMappingStep, false, "Error, failed to get mapping. Error was: %s", serviceErr.Description)
	}
	logs.LogStep(getMappingStep, true, "Success. Mapping: %+v", mapping)

	eventOptions := &eventsModel.Options{
		Args: map[string]interface{}{
			"id":   retrievedTransformation.Id,
			"name": schedule.Event.Name,
		},
	}
	eventCount := 2
	for i := 0; i < eventCount; i++ {
		_, mErr := test.eventsClient.InvokeEvent(test.Users[0].AccountId, mapping.EventId, eventOptions)
		if mErr != nil {
			return logs.LogStep(triggerEventStep, false, "Failed to invoke event. %s", mErr.Description)
		}
		time.Sleep(config.Configuration.RepeatExecutionSleep)
	}

	results, serviceErr := test.NorthstarApiClient.GetTransformationResults(test.Users[0].Token.AccessToken, createdTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(getResultsStep, false, "Error, failed to retrieve execution results: %s", serviceErr.Description)
	}

	if len(results) != eventCount {
		return logs.LogStep(getResultsStep, false, "Error, transformation not running as scheduled. Received %d out of %d results. Results: %+v", len(results), eventCount, results)
	}
	logs.LogStep(getResultsStep, true, "Retrieved %d transformation results: %+v", len(results), results)

	return nil
}

func (test *TransformationExternalEvent) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	transformations, err := test.NorthstarApiClient.ListTransformations(test.Users[0].Token.AccessToken)
	if err != nil {
		logs.LogError("Error, could not list transformations. %s", err.Error())
	}

	for _, transformation := range transformations {
		if transformation.Name == externalTransformationName {
			if transformation.Scheduled {
				if err := test.NorthstarApiClient.DeleteSchedule(test.Users[0].Token.AccessToken, transformation.Id); err != nil {
					logs.LogError("Error, could not delete schedule for transformation %s. %s", transformation.Id, err.Error())
					continue
				}
			}
			if err := test.NorthstarApiClient.DeleteTransformation(test.Users[0].Token.AccessToken, transformation.Id); err != nil {
				logs.LogError("Error, could not delete transformation %s. %s", transformation.Id, err.Error())
			}
		}
	}

	return nil
}
