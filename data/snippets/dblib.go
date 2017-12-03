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

package snippets

import (
	"fmt"
	"net/http"
	"strings"

	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/snippets/model"
	"github.com/verizonlabs/northstar/data/util"
)

var (
	snippetColumns = "id, name, createdon, updatedon, runtime, mainfn, url, code, timeout, memory, callback, description, eventtype, eventid"
	sess           *gocql.Session
	lock           sync.Mutex
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

type SnippetService struct{}

func (s *SnippetService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("snippets")
	g.POST(":accountId", addSnippet)
	g.GET(":accountId", getSnippets)
	g.GET(":accountId/:snippetId", getSnippet)
	g.PUT(":accountId/:snippetId", updateSnippet)
	g.DELETE(":accountId/:snippetId", deleteSnippet)
}

func addSnippet(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var snippet = new(model.SnippetData)
	if err := c.Bind(snippet); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertSnippet.Incr()
		return
	}

	err := snippet.Validate()
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError(err.Error()))
		ErrInsertSnippet.Incr()
		return
	}

	if snippet.Id == "" {
		snippet.Id = uuid.NewV4().String()
	}

	err = addSnippetQuery(accountId, snippet)
	if err != nil {
		em := fmt.Sprintf("Error storing snippet data, %v", err)
		mlog.Error(em)
		if strings.Contains(fmt.Sprintf("%s", err), gocql.ErrNoConnections.Error()) {
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(em))
		}
		ErrInsertSnippet.Incr()
		return
	}

	InsertSnippet.Incr()
	c.String(http.StatusCreated, snippet.Id)
}

func addSnippetQuery(accountId string, c *model.SnippetData) error {
	session, err := getSession()
	if err != nil {
		return err
	}

	if _, err := database.Insert(Keyspace, SnippetsTable).
		Param("id", c.Id).
		Param("accountId", accountId).
		Param("name", c.Name).
		Param("createdon", time.Now().In(time.UTC)).
		Param("runtime", c.Runtime).
		Param("mainfn", c.MainFn).
		Param("url", c.URL).
		Param("code", c.Code).
		Param("timeout", c.Timeout).
		Param("memory", c.Memory).
		Param("callback", c.Callback).
		Param("description", c.Description).
		Param("eventtype", c.EventType).
		Param("eventid", c.EventId).
		Exec(session); err != nil {
		return err
	}

	return nil
}

func getSnippets(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	list, err := getSnippetsQuery(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError(err.Error()))
		ErrGetSnippets.Incr()
		return
	}
	GetSnippets.Incr()
	c.JSON(http.StatusOK, list)
}

func getSnippetsQuery(accountId string) ([]model.SnippetData, error) {
	mlog.Info("Retrieving snippets for account %s", accountId)

	results := make([]model.SnippetData, 0, 10)
	entry := new(model.SnippetData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	iter := session.
		Query(`SELECT `+snippetColumns+`  FROM `+SnippetsTable+` WHERE accountid=?`, accountId).Iter()
	for iter.Scan(&entry.Id,
		&entry.Name,
		&entry.CreatedOn,
		&entry.UpdatedOn,
		&entry.Runtime,
		&entry.MainFn,
		&entry.URL,
		&entry.Code,
		&entry.Timeout,
		&entry.Memory,
		&entry.Callback,
		&entry.Description,
		&entry.EventType,
		&entry.EventId) {
		results = append(results, *entry)
		entry = new(model.SnippetData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: %v", err)
		return nil, err
	}

	return results, nil
}

func getSnippet(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get snippet due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetSnippet.Incr()
		return
	}
	snippetId := c.Params.ByName("snippetId")
	mlog.Info("Retrieving snippet %s for account %s", snippetId, accountId)

	snippetData, err := getSnippetQuery(accountId, snippetId)
	if err != nil {
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
		ErrGetSnippet.Incr()
		return
	}

	mlog.Info("Snippet %s metadata retrieved for account %s", snippetData.Name, accountId)
	GetSnippet.Incr()
	c.JSON(http.StatusOK, snippetData)
}

func getSnippetQuery(accountId string, id string) (*model.SnippetData, error) {
	var snippet model.SnippetData

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, SnippetsTable).
		Value("id", &snippet.Id).
		Value("name", &snippet.Name).
		Value("createdon", &snippet.CreatedOn).
		Value("updatedon", &snippet.UpdatedOn).
		Value("runtime", &snippet.Runtime).
		Value("mainfn", &snippet.MainFn).
		Value("url", &snippet.URL).
		Value("code", &snippet.Code).
		Value("timeout", &snippet.Timeout).
		Value("memory", &snippet.Memory).
		Value("callback", &snippet.Callback).
		Value("description", &snippet.Description).
		Value("eventtype", &snippet.EventType).
		Value("eventid", &snippet.EventId).
		Where("accountid", accountId).
		Where("id", id).
		Scan(session); err != nil {
		return nil, err
	}

	return &snippet, nil
}

func updateSnippet(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to update snippet due to bad request. The account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrUpdateSnippet.Incr()
		return
	}
	id := c.Params.ByName("snippetId")

	var update = new(model.SnippetData)
	if err := c.Bind(update); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateSnippet.Incr()
		return
	}

	err := updateSnippetQuery(accountId, id, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError(err.Error()))
		ErrUpdateSnippet.Incr()
		return
	}

	UpdateSnippet.Incr()
	c.String(http.StatusOK, "")
}

func updateSnippetQuery(accountId string, id string, update *model.SnippetData) error {
	queryBuilder := database.Update(Keyspace, SnippetsTable).
		Param("updatedon", time.Now().In(time.UTC)).
		Where("accountid", accountId).
		Where("id", id)

	if update.Name != "" {
		queryBuilder = queryBuilder.Param("name", update.Name)
	}

	if update.Runtime != "" {
		queryBuilder = queryBuilder.Param("runtime", update.Runtime)
	}

	if update.MainFn != "" {
		queryBuilder = queryBuilder.Param("mainfn", update.MainFn)
	}

	if update.URL != "" {
		queryBuilder = queryBuilder.Param("url", update.URL)
	}

	if update.Code != "" {
		queryBuilder = queryBuilder.Param("code", update.Code)
	}

	if update.Timeout > 0 {
		queryBuilder = queryBuilder.Param("timeout", update.Timeout)
	}

	if update.Memory > 0 {
		queryBuilder = queryBuilder.Param("memory", update.Memory)
	}

	if update.Callback != "" {
		queryBuilder = queryBuilder.Param("callback", update.Callback)
	}

	if update.Description != "" {
		queryBuilder = queryBuilder.Param("description", update.Description)
	}

	// If event type provided, set value. Note that we always set the
	// event id. E.g., user might try to clear value.
	if update.EventType != "" {
		queryBuilder = queryBuilder.Param("eventtype", update.EventType)
		queryBuilder = queryBuilder.Param("eventid", update.EventId)
	}

	session, err := getSession()
	if err != nil {
		return err
	}

	_, err = queryBuilder.Exec(session)
	return err
}

func deleteSnippet(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete snippet due to bad request. The account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelSnippet.Incr()
		return
	}
	snippetId := c.Params.ByName("snippetId")

	err := deleteSnippetQuery(accountId, snippetId)
	if err != nil {
		em := ""
		if err == gocql.ErrNoConnections {
			em = fmt.Sprintf("bad_gateway: %+v", err)
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			em = fmt.Sprintf("internal_error: %+v", err)
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(em))
		}
		mlog.Error("deleteSnippet(): deleting snippet %s failed: %s", snippetId, em)
		ErrDelSnippet.Incr()
		return
	}

	mlog.Info("Snippet %s deleted from account %s", snippetId, accountId)
	DelSnippet.Incr()
	c.String(http.StatusOK, "")
}

func deleteSnippetQuery(accountId string, id string) error {
	session, err := getSession()
	if err != nil {
		return err
	}

	success := true
	if success, err = database.Delete(Keyspace, SnippetsTable).
		Where("accountId", accountId).
		Where("id", id).
		Exec(session); err != nil {
		return err
	}

	if !success {
		mlog.Info("Snippet %s not found in account %s.", id, accountId)
		return fmt.Errorf("Snippet not found.")
	}

	return nil
}
