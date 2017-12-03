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

package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/cron/model"
	"github.com/verizonlabs/northstar/cron/scheduler"
	"github.com/verizonlabs/northstar/cron/util"
	cronDataClient "github.com/verizonlabs/northstar/data/cron/client"
	cronDataModel "github.com/verizonlabs/northstar/data/cron/model"
	snippetsDataClient "github.com/verizonlabs/northstar/data/snippets/client"
	processingClient "github.com/verizonlabs/northstar/processing/snippets/client"
)

type CronService struct {
	snippetsDataClient snippetsDataClient.Client
	processingClient   processingClient.Client
	cronDataClient     cronDataClient.Client
	scheduler          scheduler.Scheduler
}

func NewCronService(processingClient *processingClient.SnippetsClient,
	snippetsDataClient *snippetsDataClient.SnippetsClient,
	cronDataClient *cronDataClient.CronClient) *CronService {
	return &CronService{processingClient: processingClient,
		snippetsDataClient: snippetsDataClient,
		cronDataClient:     cronDataClient,
		scheduler:          scheduler.NewScheduler(),
	}
}

func (cron *CronService) AddRoutes() {
	grp := management.Engine().Group(util.CronBasePath)
	g := grp.Group("jobs")
	g.POST(":accountId", cron.addJob)
	g.PUT(":accountId/:jobId", cron.updateJob)
	g.DELETE(":accountId/:jobId", cron.deleteJob)
}

func (cron *CronService) addJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var job = new(model.Job)
	c.Bind(job)

	err := job.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, ValidationFailed, err.Error()))
		ErrInsertJob.Incr()
		return
	}

	job.AccountId = accountId
	job.Id = uuid.NewV4().String()
	job.ProcessingClient = cron.processingClient

	_, mErr := cron.snippetsDataClient.GetSnippet(accountId, job.SnippetId)
	if mErr != nil {
		mlog.Error("Failed to get snippet: %v", mErr)
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError("Unable to retrieve snippet"))
		ErrInsertJob.Incr()
		return
	}

	err = cron.scheduler.Start([]*model.Job{job})
	if err != nil {
		mlog.Error("Failed to start job: %v", err)
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, StartJobFailed, err.Error()))
		ErrInsertJob.Incr()
		return
	}

	dJob := &cronDataModel.JobData{AccountId: accountId,
		Id:          job.Id,
		Name:        job.Name,
		SnippetId:   job.SnippetId,
		Schedule:    job.Schedule,
		Description: job.Description}
	mErr = cron.cronDataClient.AddJob(accountId, dJob)
	if mErr != nil {
		mlog.Error("Failed to add job: %v", err)
		c.JSON(http.StatusInternalServerError, mErr)
		ErrInsertJob.Incr()
		return
	}

	mlog.Debug("Job %s added", job.Id)
	InsertJob.Incr()
	c.String(http.StatusOK, job.Id)
}

func (cron *CronService) updateJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to update Job due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateJob.Incr()
		return
	}

	jobId := c.Params.ByName("jobId")
	var update = new(model.Job)
	c.Bind(update)

	mlog.Info("Updating job %s, account id: %s, name: %s, disabled: %t, snippet id: %s, "+
		"schedule: %s, description: %s", jobId, accountId, update.Name, update.Disabled,
		update.SnippetId, update.Schedule, update.Description)

	job := &cronDataModel.JobData{UpdatedOn: time.Now(),
		Disabled:    update.Disabled,
		Name:        update.Name,
		SnippetId:   update.SnippetId,
		Schedule:    update.Schedule,
		Description: update.Description}
	mErr := cron.cronDataClient.UpdateJob(accountId, jobId, job)
	if mErr != nil {
		mlog.Error("Failed to add update: %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		ErrUpdateJob.Incr()
		return
	}

	err := cron.StartScheduler()
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.NewError(http.StatusInternalServerError,
			UpdateJobFailed, err.Error()))
		ErrUpdateJob.Incr()
		return
	}

	mlog.Debug("Job %s updated", job.Id)
	UpdateJob.Incr()
	c.String(http.StatusOK, "")
}

func (cron *CronService) deleteJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete Job due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateJob.Incr()
		return
	}
	jobId := c.Params.ByName("jobId")

	mlog.Info("Deleting job %s, account id: %s", jobId, accountId)
	mErr := cron.cronDataClient.DeleteJob(accountId, jobId)
	if mErr != nil {
		mlog.Error("Failed to delete job %v", mErr)
		c.JSON(http.StatusInternalServerError, mErr)
		ErrDelJob.Incr()
		return
	}

	err := cron.StartScheduler()
	if err != nil {
		mlog.Error("Failed to start scheduler %v", err)
		c.JSON(http.StatusInternalServerError,
			management.NewError(http.StatusInternalServerError, DeleteJobFailed, err.Error()))
		ErrDelJob.Incr()
		return
	}

	mlog.Debug("Job %s deleted", jobId)
	DelJob.Incr()
	c.String(http.StatusOK, "")
}

// Start or restart scheduler
func (cron *CronService) StartScheduler() error {
	mlog.Info("Starting/restarting scheduler")
	jobs, err := cron.cronDataClient.GetAllJobs()
	if err != nil {
		return err
	}
	return cron.scheduler.Restart(ConvertJobDataArr(jobs, cron.processingClient))
}
