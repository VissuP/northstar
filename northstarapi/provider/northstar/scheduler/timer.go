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
	cronClient "github.com/verizonlabs/northstar/cron/client"
	cronModel "github.com/verizonlabs/northstar/cron/model"
	cronDataClient "github.com/verizonlabs/northstar/data/cron/client"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Defines the type used to support operations on NorthStar resources
// (e.g., snippets, events, etc.)
type NorthStarTimerScheduler struct {
	cronClient     *cronClient.CronClient
	cronDataClient *cronDataClient.CronClient
}

// Returns a new NorthStar timer schedule provider.
func NewNorthStarTimerScheduler() (Scheduler, error) {
	mlog.Debug("NewNorthStarTimerScheduler")

	cronClient, err := cronClient.NewCronClient()
	if err != nil {
		return nil, err
	}

	cronDataClient, err := cronDataClient.NewCronClient()
	if err != nil {
		return nil, err
	}

	scheduler := &NorthStarTimerScheduler{
		cronClient:     cronClient,
		cronDataClient: cronDataClient,
	}

	return scheduler, nil
}

// Creates a new time-based schedule for the specified transformation id.
func (scheduler *NorthStarTimerScheduler) Create(accountId string, transformationId string, schedule *model.Schedule) (*model.Schedule, *management.Error) {
	mlog.Debug("Create")

	// Create the job.
	job := cronModel.Job{
		Name:      schedule.Event.Name,
		SnippetId: transformationId,
		Schedule:  schedule.Event.Value,
	}

	jobId, err := scheduler.cronClient.AddJob(accountId, &job)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Create job returned error: %v", err))
	}

	// Set the schedule id to the job id.
	schedule.Id = jobId

	return schedule, nil
}

// Returns the schedule associated with the specified transformation id.
func (scheduler *NorthStarTimerScheduler) Get(accountId string, scheduleId string) (*model.Schedule, *management.Error) {
	mlog.Debug("Get")

	job, err := scheduler.cronDataClient.GetJob(accountId, scheduleId)

	if err != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get job returned error: %+v", err))
	}

	// TODO - Job needs a created on.

	schedule := &model.Schedule{
		Id:          job.Id,
		LastUpdated: job.UpdatedOn.Format(time.RFC3339),
		Event: model.Event{
			Category: model.TimerEvent,
			Name:     job.Name,
			Value:    job.Schedule,
		},
	}

	return schedule, nil
}

// Deletes the schedule associated with the specified transformation id.
func (scheduler *NorthStarTimerScheduler) Delete(accountId string, scheduleId string) *management.Error {
	mlog.Debug("Delete")

	if err := scheduler.cronClient.DeleteJob(accountId, scheduleId); err != nil {
		return management.GetExternalError(fmt.Sprintf("Delete job returned error: %+v", err))
	}

	return nil
}
