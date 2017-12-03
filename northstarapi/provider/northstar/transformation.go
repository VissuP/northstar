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
	"net/url"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	invocation "github.com/verizonlabs/northstar/data/invocations/client"
	invocationModel "github.com/verizonlabs/northstar/data/invocations/model"
	snippets "github.com/verizonlabs/northstar/data/snippets/client"
	snippetsModel "github.com/verizonlabs/northstar/data/snippets/model"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/provider/northstar/scheduler"
	"github.com/verizonlabs/northstar/northstarapi/provider/northstar/utils"
)

// Defines the type used to support operations on NorthStar resources
// (e.g., snippets, events, etc.)
type NorthStarTransformationProvider struct {
	snippetClient    *snippets.SnippetsClient
	invocationClient *invocation.InvocationClient
}

// Returns a new NorthStar transformation provider.
func NewNorthStarTransformationProvider() (*NorthStarTransformationProvider, error) {
	mlog.Info("NewNorthStarTransformationProvider")

	timerScheduler, err := scheduler.NewNorthStarTimerScheduler()
	if err != nil {
		return nil, err
	}

	eventsScheduler, err := scheduler.NewNorthStarEventScheduler()
	if err != nil {
		return nil, err
	}

	// Register Schedulers.
	scheduler.Register(model.TimerEvent, timerScheduler)
	scheduler.Register(model.DeviceEvent, eventsScheduler)

	snippetsClient, err := snippets.NewSnippetClient()
	if err != nil {
		return nil, err
	}

	invocationsClient, err := invocation.NewInvocationClient()
	if err != nil {
		return nil, err
	}

	// Create the provider.
	provider := &NorthStarTransformationProvider{
		snippetClient:    snippetsClient,
		invocationClient: invocationsClient,
	}

	return provider, nil
}

// Creates a new snippet, etc., from the specified transformation object.
func (provider *NorthStarTransformationProvider) Create(accountId string,
	transformation *model.Transformation) (*model.Transformation, *management.Error) {
	mlog.Debug("Create")

	// Create the snippet
	transformation.Scheduled = false
	snippetData := provider.toExternal(transformation)
	snippetData.CreatedOn = time.Now()

	mlog.Debug("Creating snippet: %+v", snippetData)

	// Create the snippet.
	snippetId, mErr := provider.snippetClient.AddSnippet(accountId, snippetData)

	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("Create snippet returned error: %v", mErr))
	}

	// Update transformation snippet id.
	transformation.Id = snippetId

	return transformation, nil
}

// Updates an existing snippet.
func (provider *NorthStarTransformationProvider) Update(accountId string,
	transformation *model.Transformation) *management.Error {
	mlog.Debug("Update")

	// Get the current snippet.
	mlog.Info("Updating transformation:", transformation)
	currentSnippet, mErr := provider.snippetClient.GetSnippet(accountId, transformation.Id)
	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Get snippet returned error: %v", mErr))
	}

	// Note that before we update we verify the snippet has not been schedule.
	// If the transformation has been scheduled, return error.
	currentTransformation, err := provider.fromExternal(currentSnippet)

	if err != nil {
		return model.ErrorFromExternalSnippet
	}

	if currentTransformation.Scheduled == true {
		mlog.Error("Update validation error: Transformation %s has already been scheduled.",
			transformation.Id)
		return model.ErrorTransformationScheduled
	}

	// Note that scheduled is read only value. So, we force it to false to avoid
	// clients from changing it.
	transformation.Scheduled = false
	snippet := provider.toExternal(transformation)

	// TODO - This should not be required. Service should manage the updated on.
	snippet.UpdatedOn = time.Now()

	mlog.Debug("Updating snippet: %+v", snippet)

	if mErr := provider.snippetClient.UpdateSnippet(accountId, snippet.Id, snippet); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Update snippet returned error; %v", mErr))
	}

	return nil
}

// Get transformation results for the specified transformation ID
func (provider *NorthStarTransformationProvider) Results(accountID string,
	transformationID string) ([]model.Output, *management.Error) {
	mlog.Debug("Results -- transformationID: %s", transformationID)

	invocations, mErr := provider.invocationClient.GetInvocationResults(accountID,
		transformationID,
		100)
	if mErr != nil {
		erMsg := fmt.Sprintf("Get results for transformation %s returned error: %v",
			transformationID, mErr)
		return nil, management.GetInternalError(erMsg)
	}

	var results []model.Output
	for _, invocation := range invocations {
		resultEntry := provider.fromExternalInvocation(invocation)
		results = append(results, *resultEntry)
	}

	return results, nil

}

// Get snippet for the specified account and transformation id.
func (provider *NorthStarTransformationProvider) Get(accountId string,
	transformationId string) (*model.Transformation, *management.Error) {
	mlog.Debug("Get: transformationId:%s", transformationId)

	// Get the snippet by name.
	snippet, mErr := provider.snippetClient.GetSnippet(accountId, transformationId)
	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("Get snippet returned error: %v", mErr))
	}

	// Get the transformation.
	transformation, err := provider.fromExternal(snippet)

	if err != nil {
		return nil, model.ErrorFromExternalSnippet
	}

	// If the transformation has been schedule, attempt to get the schedule.
	if transformation.Scheduled {
		mlog.Debug("Getting schedule for transformation id %s.", transformation.Id)

		// Get the scheduler for the specified event category.
		eventScheduler, err := scheduler.Get(transformation.SchedulerType)

		if err != nil {
			return nil, model.ErrorInvalidEventCategory
		}

		// Get the schedule associated with the transformation.
		schedule, mErr := eventScheduler.Get(accountId, transformation.SchedulerId)

		if mErr != nil {
			return nil,
				management.GetExternalError(fmt.Sprintf("Get schedule returned error: %v", mErr))
		}

		transformation.Schedule = schedule
	}

	mlog.Debug("Returning transformation: %+v", transformation)

	return transformation, nil
}

// Returns all the snippets for an account id.
func (provider *NorthStarTransformationProvider) List(accountId string) ([]model.Transformation,
	*management.Error) {
	mlog.Debug("List")

	// Get snippets for account id.
	snippets, mErr := provider.snippetClient.GetSnippets(accountId)
	if mErr != nil {
		return nil,
			management.GetExternalError(fmt.Sprintf("List snippets returned error: %v", mErr))
	}

	transformations := make([]model.Transformation, 0)

	for _, snippet := range snippets {
		mlog.Debug("List snippet retrieved: %v", snippet)

		transformation, err := provider.fromExternal(snippet)

		if err != nil {
			return nil, model.ErrorFromExternalSnippet
		}

		// If the transformation has been schedule, attempt to get the schedule.
		if transformation.Scheduled {
			mlog.Debug("Getting schedule for transformation id %s.", transformation.Id)

			// Get the scheduler for the specified event category.
			eventScheduler, err := scheduler.Get(transformation.SchedulerType)

			if err != nil {
				return nil, model.ErrorInvalidEventCategory
			}

			// Get the schedule associated with the transformation.
			schedule, mErr := eventScheduler.Get(accountId, transformation.SchedulerId)

			if mErr != nil {
				mlog.Error("Get schedule returned error: %v", mErr)
				continue
			}

			transformation.Schedule = schedule
		}

		transformations = append(transformations, *transformation)
	}

	mlog.Debug("Returning transformations: %+v", transformations)

	return transformations, nil
}

// Deletes transformation with the specified id.
func (provider *NorthStarTransformationProvider) Delete(accountId string,
	transformationId string) *management.Error {
	mlog.Debug("Delete: transformationId: %s", transformationId)

	// Get the snippet for the transformation id.
	snippet, mErr := provider.snippetClient.GetSnippet(accountId, transformationId)
	mlog.Debug("Delete snippet retrieved: %v", snippet)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Get snippet %s returned error: %v",
			transformationId, mErr))
	}

	// Get transformation.
	transformation, err := provider.fromExternal(snippet)

	if err != nil {
		return model.ErrorFromExternalSnippet
	}

	// Note that before we update we verify the snippet has not been schedule.
	// If the transformation has been scheduled, return error.
	if transformation.Scheduled == true {
		mlog.Error("Delete validation error: Transformation %s has already been scheduled.",
			transformation.Id)
		return model.ErrorTransformationScheduled
	}

	// Otherwise, delete the snippet
	if mErr := provider.snippetClient.DeleteSnippet(accountId, snippet.Id); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete snippet %s returned error: %v",
			snippet.Name, mErr))
	}

	return nil
}

// Creates a new schedule for the specified transformation id.
func (provider *NorthStarTransformationProvider) CreateSchedule(accountId string,
	transformationId string,
	schedule *model.Schedule) *management.Error {
	mlog.Debug("CreateSchedule")

	// Get the snippet for the transformation id.
	snippet, mErr := provider.snippetClient.GetSnippet(accountId, transformationId)
	mlog.Debug("Create schedule retrieved snippet: %v", snippet)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Get snippet %s returned error: %v",
			transformationId, mErr))
	}

	// Get transformation.
	transformation, err := provider.fromExternal(snippet)
	mlog.Debug("Create schedule generated transformation: %v", transformation)
	if err != nil {
		return model.ErrorFromExternalSnippet
	}

	// Note that before we create a new schedule we verify the snippet has not been schedule.
	// If the transformation has been scheduled, return error.
	if transformation.Scheduled == true {
		mlog.Error("Create schedule validation error: Transformation %s has already been scheduled.",
			transformation.Id)
		return model.ErrorTransformationScheduled
	}

	// Get the scheduler for the specified event category.
	eventScheduler, err := scheduler.Get(schedule.Event.Category)

	if err != nil {
		return model.ErrorInvalidEventCategory
	}

	// Schedule the transformation.
	mlog.Debug("Schedule is: %v", schedule)
	schedule, mErr = eventScheduler.Create(accountId, transformationId, schedule)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Create schedule returned error: %v", mErr))
	}

	// Note that if a schedule is properly created. We update the transformation.
	transformation.SchedulerType = schedule.Event.Category
	transformation.SchedulerId = schedule.Id
	transformation.Scheduled = true
	snippet = provider.toExternal(transformation)

	mlog.Debug("Updated snippet: %v", snippet)
	if mErr := provider.snippetClient.UpdateSnippet(accountId, snippet.Id, snippet); mErr != nil {
		// In case of error, attempt to delete the schedule.
		if mErr := eventScheduler.Delete(accountId, schedule.Id); mErr != nil {
			mlog.Error("Delete schedule error: %v", mErr)
		}

		return management.GetExternalError(fmt.Sprintf("Update snippet returned error: %v", mErr))
	}

	return nil
}

// Returns the schedule associated with the specified transformation id.
func (provider *NorthStarTransformationProvider) GetSchedule(accountId string,
	transformationId string) (*model.Schedule, *management.Error) {
	mlog.Debug("GetSchedule")

	// Get the snippet for the transformation id.
	snippet, mErr := provider.snippetClient.GetSnippet(accountId, transformationId)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get snippet %s returned error: %v",
			transformationId, mErr))
	}

	// Get transformation.
	transformation, err := provider.fromExternal(snippet)

	if err != nil {
		return nil, model.ErrorFromExternalSnippet
	}

	// If no scheduled, return not found.
	if transformation.Scheduled == false {
		mlog.Error("Get schedule validation error: Transformation %s has NOT been scheduled.",
			transformation.Id)
		return nil, management.ErrorNotFound
	}

	// Get the scheduler for the specified event category.
	eventScheduler, err := scheduler.Get(transformation.SchedulerType)

	if err != nil {
		return nil, model.ErrorInvalidEventCategory
	}

	// Get the schedule associated with the transformation.
	schedule, mErr := eventScheduler.Get(accountId, transformation.SchedulerId)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get schedule returned error: %v", mErr))
	}

	return schedule, nil
}

// Deletes the schedule associated with the specified transformation id.
func (provider *NorthStarTransformationProvider) DeleteSchedule(accountId string,
	transformationId string) *management.Error {
	mlog.Debug("DeleteSchedule")

	// Get the snippet for the transformation id.
	snippet, mErr := provider.snippetClient.GetSnippet(accountId, transformationId)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Get snippet %s returned error: %v",
			transformationId, mErr))
	}

	// Get transformation.
	transformation, err := provider.fromExternal(snippet)

	if err != nil {
		return model.ErrorFromExternalSnippet
	}

	// If no scheduled, return not found.
	if transformation.Scheduled == false {
		mlog.Error("Delete schedule validate error: Transformation %s has NOT been scheduled.",
			transformation.Id)
		return management.ErrorNotFound
	}

	// Get the scheduler for the specified schedule category.
	eventScheduler, err := scheduler.Get(transformation.SchedulerType)

	if err != nil {
		return model.ErrorInvalidEventCategory
	}

	if mErr := eventScheduler.Delete(accountId, transformation.SchedulerId); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete schedule returned error: %v", mErr))
	}

	// Note that if a schedule is properly created. We update the transformation.
	transformation.SchedulerType = ""
	transformation.SchedulerId = ""
	transformation.Scheduled = false
	snippet = provider.toExternal(transformation)

	// TODO - This should not be required. Service should manage the updated on.
	snippet.UpdatedOn = time.Now()

	if mErr = provider.snippetClient.UpdateSnippet(accountId, snippet.Id, snippet); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Update snippet returned error: %v", mErr))
	}

	return nil
}

// Helper method used to translate data service model to portal api model.
func (provider *NorthStarTransformationProvider) fromExternal(externalTransformation *snippetsModel.SnippetData) (*model.Transformation, error) {
	mlog.Debug("fromExternal")

	// Create the transformation.
	transformation := &model.Transformation{
		Id:          externalTransformation.Id,
		Name:        externalTransformation.Name,
		Description: externalTransformation.Description,
		CreatedOn:   externalTransformation.CreatedOn.Format(time.RFC3339),
		LastUpdated: externalTransformation.UpdatedOn.Format(time.RFC3339),
		Language:    externalTransformation.Runtime,
		EntryPoint:  externalTransformation.MainFn,
		Timeout:     int(externalTransformation.Timeout / 1000), // Convert to Seconds.
		Code: model.Code{
			Type:  model.SourceCodeType,                                       // default
			Url:   fmt.Sprintf("%s://code", model.SourceCodeType.GetScheme()), // default
			Value: externalTransformation.Code,
		},
		Memory: externalTransformation.Memory / bytefmt.MEGABYTE,
	}

	// Get the code type from the snippet url.
	if externalTransformation.URL != "" {
		// Decode the url.
		url, err := url.Parse(externalTransformation.URL)

		if err != nil {
			return nil, fmt.Errorf("Parse transformation url returned error: %+v", err)
		}

		// Set the code type if different from default (source) value. Note
		// that in this case the URL should be the actual value.
		if url.Scheme == model.ObjectCodeType.GetScheme() {
			transformation.Code.Type = model.ObjectCodeType
			transformation.Code.Url = transformation.Code.Value
		}
	}

	// Set scheduled properties based on event type.
	switch externalTransformation.EventType {
	case snippetsModel.TimerEventType:
		transformation.SchedulerType = model.TimerEvent
		transformation.Scheduled = true
		transformation.SchedulerId = externalTransformation.EventId
	case snippetsModel.DeviceEventType:
		transformation.SchedulerType = model.DeviceEvent
		transformation.Scheduled = true
		transformation.SchedulerId = externalTransformation.EventId
	default:
		transformation.Scheduled = false
	}

	return transformation, nil
}

// Helper method used to translate portal api model to data service model.
func (provider *NorthStarTransformationProvider) toExternal(transformation *model.Transformation) *snippetsModel.SnippetData {
	mlog.Debug("toExternal: %v", transformation)

	// Set event type.
	eventType := snippetsModel.NoneEventType

	if transformation.Scheduled == true {
		switch transformation.SchedulerType {
		case model.DeviceEvent:
			eventType = snippetsModel.DeviceEventType
		case model.TimerEvent:
			eventType = snippetsModel.TimerEventType
		}
	}

	// Set the code url.
	codeUrl := fmt.Sprintf("%s://code", model.SourceCodeType.GetScheme())

	if transformation.Code.Type == model.ObjectCodeType {
		codeUrl = transformation.Code.Value
	}

	externalTransformation := &snippetsModel.SnippetData{
		Id:          transformation.Id,
		Name:        transformation.Name,
		Description: transformation.Description,
		Runtime:     transformation.Language,
		MainFn:      transformation.EntryPoint,
		URL:         codeUrl,
		Code:        transformation.Code.Value,
		Timeout:     transformation.Timeout * 1000, // Convert to Milliseconds
		EventType:   eventType,
		EventId:     transformation.SchedulerId,
		Memory:      transformation.Memory * bytefmt.MEGABYTE,
	}

	return externalTransformation
}

// fromExternalInvocation is a helper method used to translate data service invocation models to portal api models.
func (provider *NorthStarTransformationProvider) fromExternalInvocation(externalInvocation *invocationModel.InvocationData) *model.Output {
	mlog.Debug("fromExternalInvocation")

	status, description := utils.GetOutputStatus(externalInvocation.Status, externalInvocation.ErrorDescr)

	output := &model.Output{
		Status:            status,
		StatusDescription: description,
		ExecutionOutput:   externalInvocation.Stdout,
		ElapsedTime:       int(externalInvocation.ElapsedTime),
		LastExecution:     externalInvocation.CreatedOn,
	}

	// If the status is successful. Process the results.
	if status == model.OutputSuccessStatus && externalInvocation.Result != "" {
		var results model.CellResults
		if err := json.Unmarshal([]byte(externalInvocation.Result), &results); err != nil {
			mlog.Error("Unmarshal invocation data results returned error: %v", err)
			output.Status = model.OutputInternalErrorStatus
			output.StatusDescription = "Execution was successful but results could not be parsed."
		} else {
			output.ExecutionResults = &results
		}
	}

	return output
}
