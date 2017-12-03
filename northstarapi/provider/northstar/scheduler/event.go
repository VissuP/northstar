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

package scheduler

import (
	"fmt"
	"time"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	eventDataClient "github.com/verizonlabs/northstar/data/events/client"
	eventDataModel "github.com/verizonlabs/northstar/data/events/model"
	mappingDataClient "github.com/verizonlabs/northstar/data/mappings/client"
	mappingDataModel "github.com/verizonlabs/northstar/data/mappings/model"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Defines the type used to support operations on NorthStar resources
// (e.g., snippets, events, etc.)
type NorthStarEventScheduler struct {
	eventDataClient   *eventDataClient.EventsClient
	mappingDataClient *mappingDataClient.MappingsClient
}

// Returns a new NorthStar timer schedule provider.
func NewNorthStarEventScheduler() (Scheduler, error) {
	mlog.Debug("NewNorthStarEventScheduler")

	eventsClient, err := eventDataClient.NewEventsClient()
	if err != nil {
		return nil, err
	}

	mappingsClient, err := mappingDataClient.NewMappingsClient()
	if err != nil {
		return nil, err
	}

	scheduler := &NorthStarEventScheduler{
		eventDataClient:   eventsClient,
		mappingDataClient: mappingsClient,
	}

	return scheduler, nil
}

// Creates a new event-base schedule for the specified transformation id.
func (scheduler *NorthStarEventScheduler) Create(accountId string, transformationId string, schedule *model.Schedule) (*model.Schedule, *management.Error) {
	mlog.Debug("Create")

	// TODO - Current implementation have a one to one relation between the
	// mapping, event and snippet. That means multiple snippets cannot be
	// associated with the same Event.Id. So, do we really need the Mapping resource?
	// Example. Here are the use cases:
	//	- Execute snippet on Event Name for all devices on Account Id.
	//	- Execute snippet on Event Name for ThingSpace Device Id.
	//  - Execute snippet on Event Name for all devices on all accounts.
	// So if we update the Event resource to have:
	// 	- Snippet Id
	//	- Account Id - optional.
	//  - Device Id  - optional.
	// During event execution, if an event happens, you query
	// the Event Data Service using "Event Name".
	// Then evaluate:
	//  If Account Id and Device Id empty - trigger snippet.
	//	If Account Id not empty - evaluate they match.
	//	If Device Id not empty  - evaluate they match.
	//
	// These are the queries we might need:
	//	- Get event by "name".
	//	- Get event by "id"
	//	- Get events by "Account Id"
	//	- Get events by "Snippet Id"

	// Create the event.
	event := &eventDataModel.EventData{
		Name:      schedule.Event.Name,
		CreatedOn: time.Now(),
	}

	eventId, err := scheduler.eventDataClient.AddEvent(accountId, event)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Create event returned error: %v", err))
	}

	// Create the mapping.
	mapping := &mappingDataModel.MappingsData{
		EventId:   eventId,
		SnippetId: transformationId,
		CreatedOn: time.Now(),
	}

	mappingId, err := scheduler.mappingDataClient.AddMapping(accountId, mapping)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Create mapping returned error: %v", err))
	}

	// Update schedule with mapping information.
	schedule.Id = mappingId

	return schedule, nil
}

// Returns the schedule associated with the specified transformation id.
func (scheduler *NorthStarEventScheduler) Get(accountId string, scheduleId string) (*model.Schedule, *management.Error) {
	mlog.Debug("Get")

	// Get the mapping.
	mapping, err := scheduler.mappingDataClient.GetMapping(accountId, scheduleId)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get mapping returned error: %v", err))
	}

	// TODO - We need to fix this. Event Data client needs to provide Get event.
	events, err := scheduler.eventDataClient.ListEvents(accountId)

	for _, event := range events {
		if event.Id == mapping.EventId {
			schedule := &model.Schedule{
				Id:        mapping.Id,
				CreatedOn: mapping.CreatedOn.Format(time.RFC3339),
				Event: model.Event{
					Category: model.DeviceEvent,
					Name:     event.Name,
				},
			}

			return schedule, nil
		}
	}

	return nil, management.ErrorNotFound
}

// Deletes the schedule associated with the specified transformation id.
func (scheduler *NorthStarEventScheduler) Delete(accountId string, scheduleId string) *management.Error {
	mlog.Debug("Delete")

	// Get the mapping.
	mapping, err := scheduler.mappingDataClient.GetMapping(accountId, scheduleId)

	if err != nil {
		return management.GetExternalError(fmt.Sprintf("Get mapping returned error: %v", err))
	}

	// Delete the event.
	if err := scheduler.eventDataClient.DeleteEvent(accountId, mapping.EventId); err != nil {
		return management.GetExternalError(fmt.Sprintf("Delete event returned error: %+v", err))
	}

	// Delete the mapping.
	if err := scheduler.mappingDataClient.DeleteMapping(accountId, mapping.Id); err != nil {
		return management.GetExternalError(fmt.Sprintf("Delete mapping returned error: %+v", err))
	}

	return nil
}
