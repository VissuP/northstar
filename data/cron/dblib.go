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

package cron

import (
	"net/http"

	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/config"
	"github.com/verizonlabs/northstar/data/cron/model"
	"github.com/verizonlabs/northstar/data/util"
)

var (
	jobsColumns = "datacenter, accountid, id, name, snippetid, schedule, disabled, updatedon, description"
	sess        *gocql.Session
	lock        sync.Mutex
)

// Helper method used to get/create database session.
func getSession() (*gocql.Session, error) {
	var err error

	if sess == nil || sess.Closed() {
		lock.Lock()
		defer lock.Unlock()

		if sess == nil || sess.Closed() {
			sess, err = util.NewDB(Keyspace).GetSessionWithError()
		}
	}

	return sess, err
}

type CronService struct{}

func (s *CronService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("cron")
	g.POST("/by-accountid/:accountId", addJob)
	g.GET("/all/jobs", getAllJobs)
	g.GET("/by-accountid/:accountId", getJobs)
	g.GET("/by-accountid/:accountId/:jobId", getJob)
	g.PUT("/by-accountid/:accountId/:jobId", updateJob)
	g.DELETE("/by-accountid/:accountId/:jobId", deleteJob)
}

func addJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var job = new(model.JobData)

	if err := c.Bind(job); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertJob.Incr()
		return
	}

	err := job.ValidateOnAdd()
	if err != nil {
		ErrInsertJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		ErrInsertJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		return
	}

	if _, err := database.Insert(Keyspace, JobsTable).
		Param("datacenter", config.CassandraDatacenter).
		Param("accountId", accountId).
		Param("id", job.Id).
		Param("name", job.Name).
		Param("snippetid", job.SnippetId).
		Param("schedule", job.Schedule).
		Param("disabled", job.Disabled).
		Param("updatedon", job.UpdatedOn).
		Param("description", job.Description).
		Exec(session); err != nil {
		ErrInsertJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	InsertJob.Incr()
	c.String(http.StatusCreated, "")
}

func getJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get job due to bad request. Account Id is missing")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetJob.Incr()
		return
	}

	jobId := c.Params.ByName("jobId")
	mlog.Info("Retrieving job %s for account %s", jobId, accountId)

	job, err := getJobQuery(accountId, jobId)
	if err != nil {
		errorMessage := ""
		if err == gocql.ErrNoConnections {
			errorMessage = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway, management.GetInternalError(errorMessage))
		} else {
			errorMessage = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetInternalError(errorMessage))
		}
		mlog.Error("Failed to retrieve job %s for account %s: %s", jobId, accountId, errorMessage)
		ErrGetJob.Incr()
		return
	}

	mlog.Info("Job %s metadata retrieved for account %s", job.Id, accountId)
	GetJob.Incr()
	c.JSON(http.StatusOK, job)
}

func getJobQuery(accountId string, jobId string) (*model.JobData, error) {
	var job model.JobData

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, JobsTable).
		Value("accountid", &job.AccountId).
		Value("id", &job.Id).
		Value("name", &job.Name).
		Value("snippetid", &job.SnippetId).
		Value("schedule", &job.Schedule).
		Value("disabled", &job.Disabled).
		Value("updatedon", &job.UpdatedOn).
		Value("description", &job.Description).
		Where("datacenter", config.CassandraDatacenter).
		Where("accountid", accountId).
		Where("id", jobId).
		Scan(session); err != nil {
		return nil, err
	}

	return &job, nil
}

func getAllJobs(c *gin.Context) {
	results, err := getAllJobsQuery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrGetAllJobs.Incr()
		return
	}
	GetAllJobs.Incr()
	c.JSON(http.StatusOK, results)
}

func getAllJobsQuery() ([]model.JobData, error) {
	mlog.Info("Retrieving all jobs")

	results := make([]model.JobData, 0, 10)
	entry := new(model.JobData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	var datacenter string
	iter := session.Query(`SELECT `+jobsColumns+` FROM `+JobsTable+` WHERE datacenter=?`,
		config.CassandraDatacenter).Iter()
	for iter.Scan(
		&datacenter,
		&entry.AccountId,
		&entry.Id,
		&entry.Name,
		&entry.SnippetId,
		&entry.Schedule,
		&entry.Disabled,
		&entry.UpdatedOn,
		&entry.Description) {
		results = append(results, *entry)
		entry = new(model.JobData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func getJobs(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	results, err := getJobsQuery(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrGetJobs.Incr()
		return
	}
	GetJobs.Incr()
	c.JSON(http.StatusOK, results)
}

func getJobsQuery(accountId string) ([]model.JobData, error) {
	mlog.Info("Retrieving jobs for account %s", accountId)

	results := make([]model.JobData, 0, 10)
	entry := new(model.JobData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	var datacenter string
	iter := session.Query(`SELECT `+jobsColumns+` FROM `+JobsTable+` WHERE datacenter=? AND accountid=?`,
		config.CassandraDatacenter, accountId).Iter()
	for iter.Scan(
		&datacenter,
		&entry.AccountId,
		&entry.Id,
		&entry.Name,
		&entry.SnippetId,
		&entry.Schedule,
		&entry.Disabled,
		&entry.UpdatedOn,
		&entry.Description) {
		results = append(results, *entry)
		entry = new(model.JobData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func updateJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to update job due to bad request. Account Id is missing")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateJob.Incr()
		return
	}

	jobId := c.Params.ByName("jobId")

	var update = new(model.JobData)
	if err := c.Bind(update); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateJob.Incr()
		return
	}

	err := update.ValidateOnUpdate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrUpdateJob.Incr()
		return
	}

	err = updateJobQuery(accountId, jobId, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrUpdateJob.Incr()
		return
	}

	UpdateJob.Incr()
	c.String(http.StatusOK, "")
}

func updateJobQuery(accountId string, jobId string, update *model.JobData) error {
	queryBuilder := database.Update(Keyspace, JobsTable).
		Param("disabled", update.Disabled).
		Param("updatedon", update.UpdatedOn).
		Where("datacenter", config.CassandraDatacenter).
		Where("accountid", accountId).
		Where("id", jobId)

	if update.Name != "" {
		queryBuilder = queryBuilder.Param("name", update.Name)
	}

	if update.SnippetId != "" {
		queryBuilder = queryBuilder.Param("snippetid", update.SnippetId)
	}

	if update.Schedule != "" {
		queryBuilder = queryBuilder.Param("schedule", update.Schedule)
	}

	if update.Description != "" {
		queryBuilder = queryBuilder.Param("description", update.Description)
	}

	session, err := getSession()
	if err != nil {
		return err
	}

	_, err = queryBuilder.Exec(session)
	return err
}

func deleteJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete job due to bad request. Account Id is missing")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelJob.Incr()
		return
	}

	jobId := c.Params.ByName("jobId")

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelJob.Incr()
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, JobsTable).
		Where("datacenter", config.CassandraDatacenter).
		Where("accountId", accountId).
		Where("id", jobId).
		Exec(session); err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetInternalError(em))
		}
		ErrDelJob.Incr()
		return
	}

	if !success {
		ErrDelJob.Incr()
		mlog.Error("Failed to delete cron job. Cron job %s not found in account %s", jobId, accountId)
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Cron job not found."))
		return
	}

	DelJob.Incr()
	c.String(http.StatusOK, "")
}
