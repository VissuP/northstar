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

package transformation

import (
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.Transformation), NewTransformationCRUDTest)
}

type TransformationCRUDTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewTransformationCRUDTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &TransformationCRUDTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//steps
	createTransformationStep              utils.Step = "Creating transformation"
	retrieveTransformationStep            utils.Step = "Retrieving created transformation"
	confirmTransformationNotScheduledStep utils.Step = "Confirm that transformation is not scheduled"
	getTransformationUnauthorizedStep     utils.Step = "Attempting to get transformation as an unauthorized user"
	updateTransformationStep              utils.Step = "Updating transformation"
	scheduleTransformationStep            utils.Step = "Schedule Transformation"
	retrieveScheduledTransformationStep   utils.Step = "Retrieving scheduled transformation"
	confirmTransformationScheduledStep    utils.Step = "Confirm that transformation is scheduled"
	listTransformationsStep               utils.Step = "List transformations"
	deleteScheduledTransformationStep     utils.Step = "Attempt to delete transformation without unscheduling first"
	unscheduleTransformationStep          utils.Step = "Unscheduling transformation"
	deleteTransformationStep              utils.Step = "Deleting transformation"

	//common variables
	transformationName = "NS Sim CRUD Transformation"
)

func (test *TransformationCRUDTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {

	logs.LogDebug("Execute -- TransformationCRUDTest")

	//Create transformation
	transformation := &northstarApiModel.Transformation{
		Name:        "TestTransformation",
		Description: transformationName,
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

	createdTransformation, serviceErr := test.NorthstarApiClient.CreateTransformation(test.Users[0].Token.AccessToken, transformation)
	if serviceErr != nil {
		return logs.LogStep(createTransformationStep, false, "Error, failed to create transformation with error: %s", serviceErr.Description)
	}
	logs.LogStep(createTransformationStep, true, "Success. Transformation: %+v", createdTransformation)

	retrievedTransformation, serviceErr := test.NorthstarApiClient.GetTransformation(test.Users[0].Token.AccessToken, createdTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(retrieveTransformationStep, false, "Error, failed to get created transformation")
	}
	logs.LogStep(retrieveTransformationStep, true, "Success. Transformation: %+v", retrievedTransformation)

	if retrievedTransformation.Scheduled {
		return logs.LogStep(confirmTransformationNotScheduledStep, false, "Error, unscheduled transformation reporting scheduled: %+v", retrievedTransformation)
	}
	logs.LogStep(confirmTransformationNotScheduledStep, true, "Success.")

	_, serviceErr = test.NorthstarApiClient.GetTransformation(test.Users[1].Token.AccessToken, createdTransformation.Id)
	if serviceErr == nil {
		return logs.LogStep(getTransformationUnauthorizedStep, false, "Error, retrieved transformation with unauthorized user: %s token: %s", test.Users[1].AccountId, test.Users[1].Token.AccessToken)
	}
	logs.LogStep(getTransformationUnauthorizedStep, true, "Successfully prevented unauthorized access.")

	retrievedTransformation.Description = "Transformation updated"
	serviceErr = test.NorthstarApiClient.UpdateTransformation(test.Users[0].Token.AccessToken, retrievedTransformation)
	if serviceErr != nil {
		return logs.LogStep(updateTransformationStep, false, "Error, failed to update transformation: %s", serviceErr.Description)
	}
	logs.LogStep(updateTransformationStep, true, "Successfully updated transformation.")

	//Schedule transformation
	schedule := northstarApiModel.Schedule{
		Event: northstarApiModel.Event{
			Category: northstarApiModel.TimerEvent,
			Name:     "NSSim Test Event",
			Value:    "0 */5 * * * *",
		},
	}
	serviceErr = test.NorthstarApiClient.CreateSchedule(test.Users[0].Token.AccessToken, retrievedTransformation.Id, &schedule)
	if serviceErr != nil {
		return logs.LogStep(scheduleTransformationStep, false, "Error, failed to schedule transformation. Error was: %s", serviceErr.Description)
	}
	logs.LogStep(scheduleTransformationStep, true, "Success.")

	//Get transformation. Confirm scheduled
	scheduledTransformation, serviceErr := test.NorthstarApiClient.GetTransformation(test.Users[0].Token.AccessToken, createdTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(retrieveScheduledTransformationStep, false, "Error, failed to get scheduled transformation: %s", serviceErr.Description)
	}
	logs.LogStep(retrieveScheduledTransformationStep, true, "Success. Transformation: %+v. Schedule: %+v.", scheduledTransformation, scheduledTransformation.Schedule)

	if !scheduledTransformation.Scheduled {
		return logs.LogStep(confirmTransformationScheduledStep, false, "Error, scheduled transformation reporting unscheduled: %+v", retrievedTransformation)
	}
	logs.LogStep(confirmTransformationScheduledStep, true, "Success.")

	//Get list of transformations
	transformationList, serviceErr := test.NorthstarApiClient.ListTransformations(test.Users[0].Token.AccessToken)
	if serviceErr != nil {
		return logs.LogStep(listTransformationsStep, false, "Failed to list transformations. Error was: %s", serviceErr.Description)
	}

	if len(transformationList) == 0 {
		return logs.LogStep(listTransformationsStep, false, "Error, transformation list was empty.")
	}
	logs.LogStep(listTransformationsStep, true, "Success. Transformations: %+v", transformationList)

	serviceErr = test.NorthstarApiClient.DeleteTransformation(test.Users[0].Token.AccessToken, retrievedTransformation.Id)
	if serviceErr == nil {
		return logs.LogStep(deleteScheduledTransformationStep, false, "Error, able to delete scheduled transformation")
	}
	logs.LogStep(deleteScheduledTransformationStep, true, "Success.")

	serviceErr = test.NorthstarApiClient.DeleteSchedule(test.Users[0].Token.AccessToken, retrievedTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(unscheduleTransformationStep, false, "Error, failed to unschedule transformation: %s", serviceErr.Description)
	}
	logs.LogStep(unscheduleTransformationStep, true, "Success.")

	serviceErr = test.NorthstarApiClient.DeleteTransformation(test.Users[0].Token.AccessToken, retrievedTransformation.Id)
	if serviceErr != nil {
		return logs.LogStep(deleteTransformationStep, false, "Error, failed to delete transformation: %s", serviceErr.Description)
	}
	logs.LogStep(deleteTransformationStep, true, "Success.")
	return nil
}

func (test *TransformationCRUDTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	transformations, err := test.NorthstarApiClient.ListTransformations(test.Users[0].Token.AccessToken)
	if err != nil {
		logs.LogError("Error, could not list transformations. %s", err.Error())
	}

	for _, transformation := range transformations {
		if transformation.Name == transformationName {
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
