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
	"github.com/gambol99/go-marathon"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/stream/client"
	dataModel "github.com/verizonlabs/northstar/data/stream/model"
	"github.com/verizonlabs/northstar/dpe-stream/master/cluster"
	"github.com/verizonlabs/northstar/dpe-stream/master/connection"
	"github.com/verizonlabs/northstar/dpe-stream/master/model"
	"github.com/verizonlabs/northstar/dpe-stream/master/stats"
	"github.com/verizonlabs/northstar/dpe-stream/master/util"
)

const (
	JOB_STARTED      = "STARTED"
	JOB_STOP_FAILED  = "FAILED_TO_STOP"
	JOB_START_FAILED = "FAILED_TO_START"
)

type StreamService struct {
	marathonCluster cluster.Cluster
	dataClient      client.Client
}

func NewSteamService() (*StreamService, error) {
	cluster, err := cluster.NewMarathonCluster()
	if err != nil {
		return nil, err
	}

	dataClient, err := client.NewStreamClient()
	if err != nil {
		return nil, err
	}

	return &StreamService{marathonCluster: cluster, dataClient: dataClient}, nil
}

func (s *StreamService) AddRoutes() {
	grp := management.Engine().Group(util.StreamBasePath)
	g := grp.Group("jobs")
	g.POST(":accountId", s.startJob)
	g.DELETE(":accountId/:jobId", s.stopJob)
}

func (s *StreamService) startJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("accountId is empty"))
		stats.ErrCheckAccountId.Incr()
		return
	}

	var job = new(model.StreamJob)
	if err := c.Bind(job); err != nil {
		mlog.Error("Bind error: %v", err)
		stats.ErrBindJob.Incr()
		return
	}

	err := job.Validate()
	if err != nil {
		mlog.Error("Failed to validate: %v", err)
		stats.ErrValidateJob.Incr()
		c.JSON(http.StatusBadRequest, management.GetInternalError(err.Error()))
		return
	}

	jobId := uuid.NewV4().String()
	jobData := createJobData(accountId, jobId, job)
	mErr := s.dataClient.AddJob(accountId, jobData)
	if err != nil {
		mlog.Error("Failed to add job: %v", mErr)
		stats.ErrDataAddJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(mErr.Error()))
		return
	}

	numberOfWorkers, err := connection.GetNumberOfWorkers(&job.Source)
	if err != nil {
		stats.ErrGetNumberOfWorkers.Incr()
		c.JSON(http.StatusBadRequest, management.GetInternalError(err.Error()))
		return
	}

	startJob := &cluster.StartJob{AccountId: accountId,
		JobId:        jobId,
		InvocationId: job.InvocationId,
		Instances:    numberOfWorkers,
		Source:       job.Source,
		Functions:    job.Functions}
	err = s.marathonCluster.StartJob(startJob)
	if err != nil {
		jobData := dataModel.JobData{Status: JOB_START_FAILED, ErrorDescr: err.Error()}
		mErr := s.dataClient.UpdateJob(accountId, jobId, &jobData)
		if mErr != nil {
			mlog.Error("Failed to update job: %v", mErr)
			stats.ErrDataUpdateJob.Incr()
			c.JSON(http.StatusInternalServerError, management.GetInternalError(mErr.Error()))
			return
		}

		mlog.Error("Marathon error: %v", err)
		stats.ErrMarathonStartJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	stats.StartJob.Incr()
	c.String(http.StatusCreated, jobId)
}

func createJobData(accountId string, jobId string, job *model.StreamJob) *dataModel.JobData {
	jobData := dataModel.JobData{Id: jobId,
		AccountId:    accountId,
		InvocationId: job.InvocationId,
		Memory:       job.Memory,
		Status:       JOB_STARTED,
		Source:       dataModel.Source{Name: job.Source.Name, Connection: job.Source.Connection},
		Description:  job.Description,
	}

	var functions = make([]dataModel.Function, 0)
	for _, function := range job.Functions {
		functions = append(functions, dataModel.Function{Name: function.Name,
			Evaluator: function.Evaluator})
	}

	jobData.Functions = functions
	return &jobData
}

func (s *StreamService) stopJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("accountId is empty"))
		stats.ErrCheckAccountId.Incr()
		return
	}

	jobId := c.Params.ByName("jobId")
	if jobId == "" {
		stats.ErrCheckJobId.Incr()
		c.JSON(http.StatusBadRequest, management.GetInternalError("Job id is empty"))
		return
	}

	mlog.Debug("Stopping job %s from account %s", jobId, accountId)
	err := s.marathonCluster.StopJob(accountId, jobId)
	if err != nil {
		apiErr := err.(*marathon.APIError)
		if apiErr.ErrCode == marathon.ErrCodeNotFound {
			mlog.Debug("Job %s not found on account %v", jobId, accountId)
			s.deleteJob(accountId, jobId, c)
			return
		}

		jobData := dataModel.JobData{Status: JOB_STOP_FAILED, ErrorDescr: err.Error()}
		mErr := s.dataClient.UpdateJob(accountId, jobId, &jobData)
		if mErr != nil {
			mlog.Error("Failed to update job: %v", mErr)
			stats.ErrDataUpdateJob.Incr()
			c.JSON(http.StatusInternalServerError, management.GetInternalError(mErr.Error()))
			return
		}

		stats.ErrMarathonStopJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	stats.StopJob.Incr()
	s.deleteJob(accountId, jobId, c)
}

func (s *StreamService) deleteJob(accountId string, jobId string, c *gin.Context) {
	err := s.dataClient.DeleteJob(accountId, jobId)
	if err != nil {
		mlog.Error("Failed to delete job from data service: %v", err)
		stats.ErrDataDeleteJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	stats.DataDeleteJob.Incr()
	c.String(http.StatusOK, "")
}
