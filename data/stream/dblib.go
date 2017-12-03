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

package stream

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/stream/model"
	"github.com/verizonlabs/northstar/data/util"
	"sync"
	"time"
)

var (
	jobColumns = "id, accountid, invocationId, memory, source, functions, createdon, updatedon, status, errordescr, description"
	sess       *gocql.Session
	lock       sync.Mutex
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

type StreamService struct{}

func (s *StreamService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("stream")
	g.POST(":accountId", addJob)
	g.GET(":accountId", getJobs)
	g.GET(":accountId/:jobId", getJob)
	g.PUT(":accountId/:jobId", updateJob)
	g.DELETE(":accountId/:jobId", deleteJob)
}

func addJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var job = new(model.JobData)

	err := c.Bind(job)
	if err != nil {
		ErrAddJob.Incr()
		c.JSON(http.StatusBadRequest, management.GetInternalError(err.Error()))
		return
	}

	err = job.Validate()
	if err != nil {
		ErrAddJob.Incr()
		c.JSON(http.StatusBadRequest, management.GetInternalError(err.Error()))
		return
	}

	session, err := getSession()
	if err != nil {
		ErrAddJob.Incr()
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		ErrAddJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		return
	}

	sourceByte, err := json.Marshal(job.Source)
	if err != nil {
		ErrAddJob.Incr()
		mlog.Error("Marshal err: %v", err)
		ErrAddJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	functionsBytes, err := json.Marshal(job.Functions)
	if err != nil {
		ErrAddJob.Incr()
		mlog.Error("Marshal err: %v", err)
		ErrAddJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	if _, err := database.Insert(Keyspace, JobsTable).
		Param("accountId", accountId).
		Param("id", job.Id).
		Param("invocationId", job.InvocationId).
		Param("memory", job.Memory).
		Param("source", sourceByte).
		Param("functions", functionsBytes).
		Param("createdon", time.Now().In(time.UTC)).
		Param("status", job.Status).
		Param("description", job.Description).
		Exec(session); err != nil {
		ErrAddJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	AddJob.Incr()
	c.String(http.StatusCreated, "")
}

func getJobs(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	results, err := getJobsQuery(accountId)
	if err != nil {
		ErrGetJobs.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}
	GetJobs.Incr()
	c.JSON(http.StatusOK, results)
}

func getJobsQuery(accountId string) ([]model.JobData, error) {
	mlog.Info("Retrieving Jobs for account %s", accountId)

	results := make([]model.JobData, 0, 10)
	entry := new(model.JobData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	var source []byte
	var functions []byte

	iter :=
		session.Query(`SELECT `+jobColumns+` FROM `+JobsTable+` WHERE accountid=?`, accountId).Iter()
	for iter.Scan(&entry.Id,
		&entry.AccountId,
		&entry.InvocationId,
		&entry.Memory,
		&source,
		&functions,
		&entry.CreatedOn,
		&entry.UpdatedOn,
		&entry.Status,
		&entry.ErrorDescr,
		&entry.Description) {
		entry.ByteArrToSource(source)
		entry.ByteArrToFunctions(functions)
		results = append(results, *entry)
		entry = new(model.JobData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func getJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	jobId := c.Params.ByName("jobId")
	mlog.Info("Retrieving job %s for account %s", jobId, accountId)

	job, err := getJobQuery(accountId, jobId)
	if err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			ErrGetJob.Incr()
			c.JSON(http.StatusBadGateway, management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			ErrGetJob.Incr()
			c.JSON(http.StatusInternalServerError, management.GetInternalError(em))
		}
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

	var source []byte
	var functions []byte

	if err := database.Select(Keyspace, JobsTable).
		Value("id", &job.Id).
		Value("accountid", &job.AccountId).
		Value("invocationId", &job.InvocationId).
		Value("memory", &job.Memory).
		Value("source", &source).
		Value("functions", &functions).
		Value("createdon", &job.CreatedOn).
		Value("updatedon", &job.UpdatedOn).
		Value("description", &job.Description).
		Where("accountid", accountId).
		Where("id", jobId).
		Scan(session); err != nil {
		return nil, err
	}

	job.ByteArrToSource(source)
	job.ByteArrToFunctions(functions)
	return &job, nil
}

func updateJob(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	jobId := c.Params.ByName("jobId")

	var update = new(model.JobData)
	if err := c.Bind(update); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateJob.Incr()
		return
	}

	err := updateJobQuery(accountId, jobId, update)
	if err != nil {
		ErrUpdateJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	UpdateJob.Incr()
	c.String(http.StatusOK, "")
}

func updateJobQuery(accountId string, jobId string, update *model.JobData) error {
	queryBuilder := database.Update(Keyspace, JobsTable).
		Param("updatedon", time.Now().In(time.UTC)).
		Where("accountid", accountId).
		Where("id", jobId)

	if update.Memory > 0 {
		queryBuilder = queryBuilder.Param("mem", update.Memory)
	}

	if update.Source.Name != "" {
		queryBuilder = queryBuilder.Param("source", update.Source)
	}

	if len(update.Functions) > 0 {
		queryBuilder = queryBuilder.Param("functions", update.Functions)
	}

	if update.Status != "" {
		queryBuilder = queryBuilder.Param("status", update.Status)
	}

	if update.ErrorDescr != "" {
		queryBuilder = queryBuilder.Param("errordescr", update.ErrorDescr)
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
	jobId := c.Params.ByName("jobId")

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		ErrDelJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, JobsTable).
		Where("accountId", accountId).
		Where("id", jobId).
		Exec(session); err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			ErrDelJob.Incr()
			c.JSON(http.StatusBadGateway, management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			ErrDelJob.Incr()
			c.JSON(http.StatusInternalServerError, management.GetInternalError(em))
		}
		return
	}

	if !success {
		mlog.Error("Failed to delete stream job. Stream job %s not found in account %s", jobId, accountId)
		ErrDelJob.Incr()
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Stream job not found."))
		return
	}

	DelJob.Incr()
	c.String(http.StatusOK, "")
}
