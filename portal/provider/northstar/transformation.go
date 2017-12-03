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
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/portal/model"
)

// CreateTransformation creates a new transformation.
func (provider *NorthStarPortalProvider) CreateTransformation(token string, transformation *model.Transformation) (*model.Transformation, *management.Error) {
	mlog.Debug("CreateTransformation: transformation:%+v", transformation)

	// Create the portal api (external) transformation.
	externalTransformation, mErr := provider.northstarApiClient.CreateTransformation(token, provider.toExternalTransformation(transformation))
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Create transformation returned error: %v", mErr))
	}
	createdTransformation := provider.fromExternalTransformation(externalTransformation)

	if transformation.Schedule != nil && transformation.Schedule.Event.Category != model.NoneEvent {
		mErr = provider.CreateSchedule(token, externalTransformation.Id, transformation.Schedule)

		if mErr != nil {
			return nil, management.GetExternalError(fmt.Sprintf("Create schedule returned error: %v", mErr))
		}

		createdTransformation.Scheduled = true
		createdTransformation.Schedule.Id = transformation.Schedule.Id
		createdTransformation.Schedule.CreatedOn = transformation.Schedule.CreatedOn
		createdTransformation.Schedule.Event.Category = transformation.Schedule.Event.Category
		createdTransformation.Schedule.Event.Name = transformation.Schedule.Event.Name
		createdTransformation.Schedule.Event.Value = transformation.Schedule.Event.Value
	}

	return createdTransformation, nil
}

// UpdateTransformation updates an existing transformation.
func (provider *NorthStarPortalProvider) UpdateTransformation(token string, transformation *model.Transformation) *management.Error {
	mlog.Debug("UpdateTransformation")

	// 1. Schedule provided from request
	// 2. One already exists? Delete it
	// 3. Update transformation when not being scheduled
	// 4. Create schedule if it's not none

	// Get transformation to check if it's scheduled already
	extTransformation, mErr := provider.northstarApiClient.GetTransformation(token, transformation.Id)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Failed to get transformation with error: %v", mErr))
	}

	// If scheduled, delete it
	if extTransformation.Scheduled && transformation.Schedule != nil {
		mErr := provider.northstarApiClient.DeleteSchedule(token, transformation.Id)

		if mErr != nil {
			return management.GetExternalError(fmt.Sprintf("Failed to delete schedule with error: %v", mErr))
		}
	}

	// Update the portal api (external) transformation.
	if mErr := provider.northstarApiClient.UpdateTransformation(token, provider.toExternalTransformation(transformation)); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Update transformation returned error: %v", mErr))
	}

	if transformation.Schedule.Event.Category != model.NoneEvent && transformation.Schedule != nil {
		// Create a new schedule for the transformation
		mErr = provider.CreateSchedule(token, transformation.Id, transformation.Schedule)

		if mErr != nil {
			return management.GetExternalError(fmt.Sprintf("Failed to create schedule with error: %v", mErr))
		}
	}

	return nil
}

// ListTransformations returns transformations associated with authenticated user.
func (provider *NorthStarPortalProvider) ListTransformations(token string) ([]model.Transformation, *management.Error) {
	mlog.Debug("ListTransformations")

	// Get the portal api (external) transformations for user (i.e., from token).
	externalTransformations, mErr := provider.northstarApiClient.ListTransformations(token)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("List transformations returned error: %v", mErr))
	}

	// Translate to portal transformations.
	var transformations []model.Transformation

	for _, externalTransformation := range externalTransformations {
		transformation := provider.fromExternalTransformation(&externalTransformation)
		transformations = append(transformations, *transformation)
	}

	return transformations, nil
}

// GetTransformation returns transformation with the specified id.
func (provider *NorthStarPortalProvider) GetTransformation(token string, transformationId string) (*model.Transformation, *management.Error) {
	mlog.Debug("GetTransformation")

	// Get the portal api (external) transformation with id.
	externalTransformation, mErr := provider.northstarApiClient.GetTransformation(token, transformationId)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get transformation returned error: %v", mErr))
	}

	return provider.fromExternalTransformation(externalTransformation), nil
}

// GetTransformationResults returns the collection of execution results for a transformation.
func (provider *NorthStarPortalProvider) GetTransformationResults(token string, transformationID string) ([]model.Output, *management.Error) {
	mlog.Debug("GetTransformationHistory")

	// Get the external results with id.
	externalResults, mErr := provider.northstarApiClient.GetTransformationResults(token, transformationID)
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get transformation results returned error: %v", mErr))
	}

	var results []model.Output
	for _, invocation := range externalResults {
		resultsEntry := provider.fromExternalOutput(invocation)
		results = append(results, *resultsEntry)
	}
	return results, nil
}

// DeleteTransformation deletes the transformation with the specified id.
func (provider *NorthStarPortalProvider) DeleteTransformation(token string, transformationId string) *management.Error {
	mlog.Debug("DeleteTransformation: transformationId:%s", transformationId)

	extTransformation, mErr := provider.northstarApiClient.GetTransformation(token, transformationId)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Failed to get transformation with error: %v", mErr))
	}

	if extTransformation.Scheduled {
		mErr = provider.northstarApiClient.DeleteSchedule(token, transformationId)

		if mErr != nil {
			return management.GetExternalError(fmt.Sprintf("Failed to delete schedule with error: %v", mErr))
		}
	}

	if mErr := provider.northstarApiClient.DeleteTransformation(token, transformationId); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete transformation returned error: %v", mErr))
	}

	return nil
}

// Get Schemas that thingspace provides
func (provider *NorthStarPortalProvider) GetScheduleEventSchemas(token string) ([]model.ScheduleEventSchema, *management.Error) {
	mlog.Debug("GetScheduleEventSchemas")

	// Create NorthStar time-base event schemas.
	schemas := []model.ScheduleEventSchema{
		model.ScheduleEventSchema{
			Category:    model.TimerEvent,
			Description: model.TimerEventDescription,
			Type:        "string",
			Semantic:    "Cron expression",
			Fields: []model.ScheduleEventFieldSchema{
				//Rules are written by columns -- Minutes Hours DayOfMonth Month DayOfWeek
				model.ScheduleEventFieldSchema{
					Name:     "Every 5 Minutes",
					Constant: true,
					Value:    "0 */5 * * * *",
				}, model.ScheduleEventFieldSchema{
					Name:     "Every 10 Minutes",
					Constant: true,
					Value:    "0 */10 * * * *",
				},
				model.ScheduleEventFieldSchema{
					Name:     "Every 30 Minutes",
					Constant: true,
					Value:    "0 */30 * * * *",
				},
				model.ScheduleEventFieldSchema{
					Name:     "Every Hour",
					Constant: true,
					Value:    "0 0 */1 * * *",
				},
				model.ScheduleEventFieldSchema{
					Name:     "Every 6 Hours",
					Constant: true,
					Value:    "0 0 */6 * * *",
				},
				model.ScheduleEventFieldSchema{
					Name:     "Every 12 Hours",
					Constant: true,
					Value:    "0 0 */12 * * *",
				},
				model.ScheduleEventFieldSchema{
					Name:     "Everyday at Midnight",
					Constant: true,
					Value:    "0 0 0 * * *",
				},
			},
		},
		model.ScheduleEventSchema{
			Category:    model.NoneEvent,
			Description: model.NoneEventDescription,
		},
	}

	// Get ThingSpace device models needed to create device-base event schemas.
	deviceModels, mErr := provider.thingSpaceClient.GetModels(token)

	if mErr != nil {
		mlog.Error("Get model returned error: %v", mErr)
		return nil, management.GetExternalError(fmt.Sprintf("Failed to get models from thingspace with error: %s", mErr.Description))
	}

	// For every device model create a new event schema.
	for _, deviceModel := range deviceModels {
		// Create the event schema.
		schema := model.ScheduleEventSchema{
			Category:    model.DeviceEvent,
			DeviceKind:  deviceModel.ModelKind,
			Description: model.DeviceEventDescription,
			Fields:      make([]model.ScheduleEventFieldSchema, 0),
		}

		for deviceModelFieldName, deviceModelField := range deviceModel.Fields {
			schemaField := model.ScheduleEventFieldSchema{
				Name: deviceModelFieldName,
				Type: deviceModelField.Type,
			}
			schema.Fields = append(schema.Fields, schemaField)
		}

		if len(schema.Fields) > 0 {
			schemas = append(schemas, schema)
		}
	}
	return schemas, nil
}

// Creates a new schedule for the specified transformationId.
func (provider *NorthStarPortalProvider) CreateSchedule(token string, transformationId string, schedule *model.Schedule) *management.Error {
	mlog.Debug("CreateSchedule")

	// Create portal api schedule.
	externalSchedule := &northstarApiModel.Schedule{
		Event: northstarApiModel.Event{
			Category: schedule.Event.Category,
			Name:     schedule.Event.Name,
			Value:    schedule.Event.Value,
		},
	}

	if mErr := provider.northstarApiClient.CreateSchedule(token, transformationId, externalSchedule); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Create schedule returned error: %v", mErr))
	}

	return nil
}

// Returns the schedule associated with the specified transformationId.
func (provider *NorthStarPortalProvider) GetSchedule(token string, transformationId string) (*model.Schedule, *management.Error) {
	mlog.Debug("GetSchedule")

	externalSchedule, mErr := provider.northstarApiClient.GetSchedule(token, transformationId)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get schedule returned error: %v", mErr))
	}

	schedule := &model.Schedule{
		Id:          externalSchedule.Id,
		CreatedOn:   externalSchedule.CreatedOn,
		LastUpdated: externalSchedule.LastUpdated,
		Event: model.ScheduleEvent{
			Category: externalSchedule.Event.Category,
			Name:     externalSchedule.Event.Name,
			Value:    externalSchedule.Event.Value,
		},
	}

	return schedule, nil
}

// Deletes the schedule associated with the specified transformationId.
func (provider *NorthStarPortalProvider) DeleteSchedule(token string, transformationId string) *management.Error {
	mlog.Debug("DeleteSchedule")

	if mErr := provider.northstarApiClient.DeleteSchedule(token, transformationId); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete schedule returned error: %v", mErr))
	}

	return nil
}

// ExecuteTransformation submits an execution request for a transformation.
func (provider *NorthStarPortalProvider) ExecuteTransformation(token string, callbackURL string, transformation *model.Transformation) *management.Error {
	mlog.Debug("ExecuteTransformation")

	// Create portal api (external) transformation.
	externalTransformation := provider.toExternalTransformation(transformation)

	if mErr := provider.northstarApiClient.ExecuteTransformation(token, callbackURL, externalTransformation); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Execute transformation returned error: %v", mErr))
	}

	return nil
}

// ToExternalTransformation is a helper method used to translate portal model to portal api model.
func (provider *NorthStarPortalProvider) toExternalTransformation(transformation *model.Transformation) *northstarApiModel.Transformation {
	mlog.Debug("toExternalTransformation: transformation:%+v", transformation)

	// Create external transformation.
	externalTransformation := &northstarApiModel.Transformation{
		Id:          transformation.Id,
		Name:        transformation.Name,
		Description: transformation.Description,
		Timeout:     transformation.Timeout,
		Language:    transformation.Language,
		EntryPoint:  transformation.EntryPoint,
		Code: northstarApiModel.Code{
			Type:  northstarApiModel.SourceCodeType,
			Value: transformation.Code.Value,
		},
	}

	// Update default code type.
	if transformation.Code.Type == northstarApiModel.ObjectCodeType.ToString() {
		externalTransformation.Code.Type = northstarApiModel.ObjectCodeType
	}

	// Set the input arguments, if provided.

	if transformation.Arguments != nil {
		externalTransformation.Arguments = make(map[string]interface{})
		for key, value := range transformation.Arguments {
			externalTransformation.Arguments[key] = value
		}
	}

	return externalTransformation
}

// FromExternalTransformation is a helper method used to translate portal api model to portal model.
func (provider *NorthStarPortalProvider) fromExternalTransformation(externalTransformation *northstarApiModel.Transformation) *model.Transformation {
	mlog.Debug("fromExternalTransformation: externalTransformation:%+v", externalTransformation)

	transformation := &model.Transformation{
		Id:          externalTransformation.Id,
		CreatedOn:   externalTransformation.CreatedOn,
		LastUpdated: externalTransformation.LastUpdated,
		Name:        externalTransformation.Name,
		Description: externalTransformation.Description,
		Scheduled:   externalTransformation.Scheduled,
		Timeout:     externalTransformation.Timeout,
		EntryPoint:  externalTransformation.EntryPoint,
		Language:    externalTransformation.Language,
		Code: model.Code{
			Type:  externalTransformation.Code.Type.ToString(),
			Value: externalTransformation.Code.Value,
		},
	}

	//All events have an event type, by default, the event type is None (unscheduled).
	transformation.Schedule = &model.Schedule{
		Event: model.ScheduleEvent{
			Category: model.NoneEvent,
		},
	}

	if externalTransformation.Schedule != nil {
		transformation.Schedule.Id = externalTransformation.Schedule.Id
		transformation.Schedule.CreatedOn = externalTransformation.Schedule.CreatedOn
		transformation.Schedule.Event.Category = externalTransformation.Schedule.Event.Category
		transformation.Schedule.Event.Name = externalTransformation.Schedule.Event.Name
		transformation.Schedule.Event.Value = externalTransformation.Schedule.Event.Value
	}

	return transformation
}

func (provider *NorthStarPortalProvider) fromExternalOutput(externalOutput northstarApiModel.Output) *model.Output {
	mlog.Debug("fromExternalInvocation")

	output := &model.Output{
		State:         externalOutput.Status.ToString(),
		Stderr:        externalOutput.StatusDescription,
		Stdout:        externalOutput.ExecutionOutput,
		ElapsedTime:   externalOutput.ElapsedTime,
		LastExecution: externalOutput.LastExecution,
	}

	if externalOutput.ExecutionResults != nil {
		output.Results = &model.CellResults{
			Type:    string(externalOutput.ExecutionResults.Type),
			Content: externalOutput.ExecutionResults.Content,
		}
	}

	return output
}
