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

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.TimerExecution), NewTimerExecutionTest)
}

type TimerExecutionTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewTimerExecutionTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &TimerExecutionTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//Steps
	createTransformationStep   utils.Step = "Create transformation"
	scheduleTransformationStep utils.Step = "Schedule transformation"
	transformationResultsStep  utils.Step = "Check transformation results."

	//common variables
	timerTransformationName = "Timer Transformation"
)

func (test *TimerExecutionTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- TimerExecutionTest")

	//Create transformation
	transformation := &northstarApiModel.Transformation{
		Name:        timerTransformationName,
		Description: "NS Sim Test Transformation",
		Timeout:     180,
		EntryPoint:  "main",
		Language:    "lua",
		Code: northstarApiModel.Code{
			Type: northstarApiModel.SourceCodeType,
			Value: `local output = require('nsOutput')

				function main()
				output.print("Hello World!")
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
	logs.LogStep(createTransformationStep, true, "Success. Transformation: %+v", createdTransformation)

	//Schedule transformation
	schedule := northstarApiModel.Schedule{
		Event: northstarApiModel.Event{
			Category: northstarApiModel.TimerEvent,
			Name:     "NSSim Test Event",
			Value:    "*/45 * * * * *",
		},
	}
	serviceErr = test.NorthstarApiClient.CreateSchedule(test.Users[0].Token.AccessToken, createdTransformation.Id, &schedule)
	if serviceErr != nil {
		return logs.LogStep(scheduleTransformationStep, false, "Error, failed to schedule transformation. Error was: %s", serviceErr.Description)
	}
	defer test.DeleteSchedule(logs, test.Users[0].Token.AccessToken, createdTransformation.Id)
	logs.LogStep(scheduleTransformationStep, true, "Success.")

	logs.LogInfo("Sleep for 100 seconds to let transformation trigger")
	time.Sleep(100 * time.Second)

	results, serviceErr := test.NorthstarApiClient.GetTransformationResults(test.Users[0].Token.AccessToken, createdTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(transformationResultsStep, false, "Error, failed to retrieve execution results: %s", serviceErr.Description)
	}

	if len(results) < 2 {
		return logs.LogStep(transformationResultsStep, false, "Error, transformation not running as scheduled. Received %d out of %d results. Transformation: %+v", len(results), 2, transformation)
	}
	logs.LogStep(transformationResultsStep, true, "Retrieved %d transformation results: %+v", len(results), results)

	return nil
}

func (test *TimerExecutionTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	transformations, err := test.NorthstarApiClient.ListTransformations(test.Users[0].Token.AccessToken)
	if err != nil {
		logs.LogError("Error, could not list transformations. %s", err.Error())
	}

	for _, transformation := range transformations {
		if transformation.Name == timerTransformationName {
			if transformation.Scheduled {
				logs.LogInfo("Cleaning up schedule for transformation %s", transformation.Id)
				if err := test.NorthstarApiClient.DeleteSchedule(test.Users[0].Token.AccessToken, transformation.Id); err != nil {
					logs.LogError("Error, could not delete schedule for transformation %s. %s", transformation.Id, err.Error())
					continue
				}
			}
			logs.LogInfo("Cleaning up transformation %s", transformation.Id)
			if err := test.NorthstarApiClient.DeleteTransformation(test.Users[0].Token.AccessToken, transformation.Id); err != nil {
				logs.LogError("Error, could not delete transformation %s. %s", transformation.Id, err.Error())
			}
		}
	}

	return nil
}
