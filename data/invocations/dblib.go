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

package invocations

import (
	"fmt"
	"net/http"

	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/invocations/model"
	"github.com/verizonlabs/northstar/data/util"
	"strconv"
)

var (
	invocationColumns = "id, rteid, snippetid, partition, createdon, startedon, finishedon, updatedon, elapsedtime, runtime, mainfn, url, code, timeout, memory, callback, args, stdout, result, status, errordescr"
	sess              *gocql.Session
	lock              sync.Mutex
)

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

type InvocationService struct{}

func (s *InvocationService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("invocations")
	g.POST("/invocation/:accountId", addInvocation)
	g.GET("/invocation/:accountId/:invocationId", getInvocation)
	g.POST("/invocation/:accountId/:invocationId", updateInvocation)
	g.DELETE("/invocation/:accountId/:invocationId", deleteInvocation)
	g.GET("/history/by-account/:accountId/:limit", getInvocationsByAccountId)
	g.GET("/history/by-snippet/:accountId/:snippetId/:limit", getInvocationHistory)
}

func addInvocation(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var invocation = new(model.InvocationData)
	if err := c.Bind(invocation); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertInvocation.Incr()
		return
	}

	err := invocation.ValidateOnAdd()
	if err != nil {
		ErrInsertInvocation.Incr()
		c.JSON(http.StatusBadGateway,
			management.GetExternalError(err.Error()))
		return
	}

	invocationId := gocql.TimeUUID()
	mlog.Debug("Adding invocation %s for account %s", invocationId.String(), accountId)

	err = addInvocationQuery(accountId, invocationId, invocation)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), gocql.ErrNoConnections.Error()) {
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(err.Error()))
		}
		ErrInsertInvocation.Incr()
		return
	}

	InsertInvocation.Incr()
	c.String(http.StatusCreated, invocationId.String())
}

func addInvocationQuery(accountId string,
	invocationId gocql.UUID,
	invocation *model.InvocationData) error {
	session, err := getSession()
	if err != nil {
		return err
	}

	args, err := invocation.ArgsToByteArr()
	if err != nil {
		return err
	}

	if _, err := database.Insert(Keyspace, InvocationsTable).
		Param("id", invocationId).
		Param("accountid", accountId).
		Param("snippetid", invocation.SnippetId).
		Param("createdon", time.Now().In(time.UTC)).
		Param("runtime", invocation.Runtime).
		Param("mainfn", invocation.MainFn).
		Param("url", invocation.URL).
		Param("code", invocation.Code).
		Param("timeout", invocation.Timeout).
		Param("memory", invocation.Memory).
		Param("callback", invocation.Callback).
		Param("args", args).
		Param("status", invocation.Status).
		Exec(session); err != nil {
		return err
	}

	return nil
}

func getInvocation(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get invocation due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		return
	}
	invocationId := c.Params.ByName("invocationId")
	mlog.Debug("Retrieving invocation %s data for account %s", invocationId, accountId)

	var invocationData = new(model.InvocationData)
	var args []byte

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetInvocationById.Incr()
		return
	}

	if err := database.Select(Keyspace, InvocationsTable).
		Value("id", &invocationData.Id).
		Value("rteid", &invocationData.RTEId).
		Value("snippetid", &invocationData.SnippetId).
		Value("partition", &invocationData.Partition).
		Value("createdon", &invocationData.CreatedOn).
		Value("startedon", &invocationData.StartedOn).
		Value("finishedon", &invocationData.FinishedOn).
		Value("updatedon", &invocationData.UpdatedOn).
		Value("elapsedtime", &invocationData.ElapsedTime).
		Value("runtime", &invocationData.Runtime).
		Value("mainfn", &invocationData.MainFn).
		Value("url", &invocationData.URL).
		Value("code", &invocationData.Code).
		Value("timeout", &invocationData.Timeout).
		Value("memory", &invocationData.Memory).
		Value("callback", &invocationData.Callback).
		Value("args", &args).
		Value("stdout", &invocationData.Stdout).
		Value("result", &invocationData.Result).
		Value("status", &invocationData.Status).
		Value("errordescr", &invocationData.ErrorDescr).
		Where("accountid", accountId).
		Where("id", invocationId).
		Scan(session); err != nil {
		em := fmt.Sprintf("Error getting invocation result for id %s: %v", invocationId, err)
		mlog.Error(em)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(em))
		ErrGetInvocationById.Incr()
		return
	}

	invocationData.ByteArrToArgs(args)

	mlog.Debug("Invocation %s data retrieved for account %s", invocationId, accountId)
	GetInvocationById.Incr()
	c.JSON(http.StatusOK, invocationData)
}

func updateInvocation(c *gin.Context) {
	var input = new(model.InvocationData)
	if err := c.Bind(input); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateInvocation.Incr()
		return
	}

	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to update invocation due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateInvocation.Incr()
		return
	}

	invocationId := c.Params.ByName("invocationId")
	if invocationId == "" {
		mlog.Error("Failed to update invocation due to bad request. Invocation Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.InvocationIdMissing))
		ErrUpdateInvocation.Incr()
		return
	}

	err := updateInvocationQuery(accountId, invocationId, input)
	if err != nil {
		mlog.Error("Error executing invocation update:", err)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrUpdateInvocation.Incr()
		return
	}

	mlog.Debug("Invocation %s updated for account %s", invocationId, accountId)
	UpdateInvocation.Incr()
	c.String(http.StatusCreated, "")
}

func updateInvocationQuery(accountId string, invocationId string, input *model.InvocationData) error {
	queryBuilder := database.Update(Keyspace, InvocationsTable).
		Param("updatedon", time.Now().In(time.UTC)).
		Where("accountid", accountId).
		Where("id", invocationId)

	if input.SnippetId != "" {
		queryBuilder = queryBuilder.Param("snippetid", input.SnippetId)
	}

	if input.RTEId != "" {
		queryBuilder = queryBuilder.Param("rteid", input.RTEId)
	}

	if input.Partition >= 0 {
		queryBuilder = queryBuilder.Param("partition", input.Partition)
	}

	if !input.StartedOn.IsZero() {
		queryBuilder = queryBuilder.Param("startedon", input.StartedOn)
	}

	if !input.FinishedOn.IsZero() {
		queryBuilder = queryBuilder.Param("finishedon", input.FinishedOn)
	}

	if input.ElapsedTime > 0 {
		queryBuilder = queryBuilder.Param("elapsedtime", input.ElapsedTime)
	}

	if input.Stdout != "" {
		queryBuilder = queryBuilder.Param("stdout", input.Stdout)
	}

	if input.Result != "" {
		queryBuilder = queryBuilder.Param("result", input.Result)
	}

	if input.Status != "" {
		queryBuilder = queryBuilder.Param("status", input.Status)
	}

	if input.ErrorDescr != "" {
		queryBuilder = queryBuilder.Param("errordescr", input.ErrorDescr)
	}

	session, err := getSession()
	if err != nil {
		return err
	}

	_, err = queryBuilder.Exec(session)
	if err != nil {
		return err
	}

	return nil
}

func deleteInvocation(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete invocation due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelInvocation.Incr()
		return
	}

	invocationId := c.Params.ByName("invocationId")
	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelInvocation.Incr()
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, InvocationsTable).
		Where("accountId", accountId).
		Where("id", invocationId).
		Exec(session); err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway, management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetInternalError(em))
		}
		ErrDelInvocation.Incr()
		return
	}

	if !success {
		ErrDelInvocation.Incr()
		mlog.Error("Invocation %s not found in account %s", invocationId, accountId)
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Invocation not found."))
		return
	}

	mlog.Debug("Invocation %s deleted from account %s", invocationId, accountId)
	DelInvocation.Incr()
	c.String(http.StatusOK, "")
}

func getInvocationsByAccountId(c *gin.Context) {
	accountId := c.Params.ByName("accountId")

	limit, err := strconv.Atoi(c.Params.ByName("limit"))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to convert limit to int: %v", err.Error())
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetInvocationsByAccountId.Incr()
		return
	}

	mlog.Debug("Retrieving all invocations for account %s with limit %v", accountId, limit)

	if limit < 1 {
		errorMessage := fmt.Sprintf("Limit must be greater than 0")
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetInvocationsByAccountId.Incr()
		return
	}

	results := make([]model.InvocationData, 0)
	var invocationData = new(model.InvocationData)
	var args []byte

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetInvocationsByAccountId.Incr()
		return
	}

	iter := session.Query(`SELECT `+invocationColumns+` FROM `+InvocationsTable+
		` WHERE accountid=? LIMIT ?`, accountId, limit).Iter()
	for iter.Scan(
		&invocationData.Id,
		&invocationData.RTEId,
		&invocationData.SnippetId,
		&invocationData.Partition,
		&invocationData.CreatedOn,
		&invocationData.StartedOn,
		&invocationData.FinishedOn,
		&invocationData.UpdatedOn,
		&invocationData.ElapsedTime,
		&invocationData.Runtime,
		&invocationData.MainFn,
		&invocationData.URL,
		&invocationData.Code,
		&invocationData.Timeout,
		&invocationData.Memory,
		&invocationData.Callback,
		&args,
		&invocationData.Stdout,
		&invocationData.Result,
		&invocationData.Status,
		&invocationData.ErrorDescr) {
		invocationData.ByteArrToArgs(args)
		results = append(results, *invocationData)
		invocationData = new(model.InvocationData)
	}

	mlog.Debug("Invocations retrieved for account %s", accountId)
	GetInvocationsByAccountId.Incr()
	c.JSON(http.StatusOK, results)
}

func getInvocationHistory(context *gin.Context) {
	accountId := context.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get invocation history due to bad request. Account Id is missing.")
		context.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetInvocationHistory.Incr()
		return
	}

	snippetId := context.Params.ByName("snippetId")
	if snippetId == "" {
		mlog.Error("Failed to get invocation history due to bad request. Snippet Id is missing.")
		context.JSON(http.StatusBadRequest, management.GetBadRequestError(util.SnippetIdMissing))
		ErrGetInvocationHistory.Incr()
		return
	}

	limit := context.Params.ByName("limit")
	results := make([]model.InvocationData, 0)
	var invocationData = new(model.InvocationData)

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetInvocationHistory.Incr()
		return
	}

	var args []byte

	iter := session.Query(`SELECT `+invocationColumns+` FROM `+InvocationsTable+
		` WHERE accountid=? and snippetid=? LIMIT ?`, accountId, snippetId, limit).Iter()
	for iter.Scan(
		&invocationData.Id,
		&invocationData.RTEId,
		&invocationData.SnippetId,
		&invocationData.Partition,
		&invocationData.CreatedOn,
		&invocationData.StartedOn,
		&invocationData.FinishedOn,
		&invocationData.UpdatedOn,
		&invocationData.ElapsedTime,
		&invocationData.Runtime,
		&invocationData.MainFn,
		&invocationData.URL,
		&invocationData.Code,
		&invocationData.Timeout,
		&invocationData.Memory,
		&invocationData.Callback,
		&args,
		&invocationData.Stdout,
		&invocationData.Result,
		&invocationData.Status,
		&invocationData.ErrorDescr) {
		invocationData.ByteArrToArgs(args)
		results = append(results, *invocationData)
		invocationData = new(model.InvocationData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error getting invocation results:", err)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrGetInvocationHistory.Incr()
		return
	}

	GetInvocationHistroy.Incr()
	context.JSON(http.StatusOK, results)
}
