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

package mappings

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
	"github.com/verizonlabs/northstar/data/mappings/model"
	"github.com/verizonlabs/northstar/data/util"
)

var (
	mappingsColumns = "id, eventid, snippetid, createdon"
	sess            *gocql.Session
	lock            sync.Mutex
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

type MappingsService struct{}

func (s *MappingsService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("mappings")
	g.POST("/by-accountid/:accountId", addMapping)
	g.GET("/by-accountid/:accountId", getMappings)
	g.GET("/by-accountid/:accountId/:mappingId", getMapping)
	g.GET("/by-eventid/:accountId/:eventId", getMappingByEventId)
	g.DELETE("/by-accountid/:accountId/:mappingId", deleteMapping)
}

func addMapping(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var mapping = new(model.MappingsData)
	if err := c.Bind(mapping); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertMapping.Incr()
		return
	}

	err := mapping.ValidateOnAdd()
	if err != nil {
		c.JSON(http.StatusNotFound,
			management.GetNotFoundError(err.Error()))
		ErrInsertMapping.Incr()
		return
	}

	if mapping.Id == "" {
		mapping.Id = uuid.NewV4().String()
	}

	err = addMappingQuery(accountId, mapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			management.GetInternalError(err.Error()))
		ErrInsertMapping.Incr()
		return
	}

	mlog.Info("Snippet %s -> %s event mapping added", mapping.SnippetId, mapping.EventId)
	InsertMapping.Incr()
	c.String(http.StatusCreated, mapping.Id)
}

func addMappingQuery(accountId string, mapping *model.MappingsData) error {
	session, err := getSession()
	if err != nil {
		return err
	}

	if _, err := database.Insert(Keyspace, MappingsTable).
		Param("id", mapping.Id).
		Param("accountid", accountId).
		Param("snippetId", mapping.SnippetId).
		Param("eventId", mapping.EventId).
		Param("createdon", time.Now().In(time.UTC)).
		Exec(session); err != nil {
		return err
	}

	return nil
}

func getMappings(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	results, err := getMappingsQuery(accountId)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), gocql.ErrNoConnections.Error()) {
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(err.Error()))
		} else {
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(err.Error()))
		}
		ErrGetMappings.Incr()
		return
	}
	GetMappings.Incr()
	c.JSON(http.StatusOK, results)
}

func getMappingsQuery(accountId string) ([]model.MappingsData, error) {
	mlog.Info("Retrieving events mappings for account %s", accountId)

	results := make([]model.MappingsData, 0, 10)
	entry := new(model.MappingsData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	iter := session.
		Query(`SELECT `+mappingsColumns+` FROM `+MappingsTable+` WHERE accountid=?`, accountId).Iter()
	for iter.Scan(&entry.Id,
		&entry.EventId,
		&entry.SnippetId,
		&entry.CreatedOn) {
		results = append(results, *entry)
		entry = new(model.MappingsData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func getMapping(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get mapping due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetMapping.Incr()
		return
	}

	mappingId := c.Params.ByName("mappingId")
	results, err := getMappingQuery(accountId, mappingId)
	if err != nil {
		em := fmt.Sprintf("Error retrieving mappings %v", err)
		if strings.Contains(fmt.Sprintf("%s", err), gocql.ErrNoConnections.Error()) {
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(em))
		}
		ErrGetMapping.Incr()
		return
	}

	GetMapping.Incr()
	c.JSON(http.StatusOK, results)
}

func getMappingQuery(accountId string, mappingId string) (*model.MappingsData, error) {
	var mappingData = new(model.MappingsData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, MappingsTable).
		Value("id", &mappingData.Id).
		Value("eventid", &mappingData.EventId).
		Value("snippetid", &mappingData.SnippetId).
		Value("createdon", &mappingData.CreatedOn).
		Where("accountid", accountId).
		Where("id", mappingId).
		Scan(session); err != nil {
		return nil, err
	}

	return mappingData, nil
}

func getMappingByEventId(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to get mapping by event Id due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrGetMappingByEventId.Incr()
		return
	}

	eventId := c.Params.ByName("eventId")
	mlog.Info("Retrieving mapping for account %s and event %s", accountId, eventId)

	results, err := getMappingByEventIdQuery(accountId, eventId)
	if err != nil {
		em := fmt.Sprintf("Error retrieving mappings %v", err)
		if strings.Contains(fmt.Sprintf("%s", err), gocql.ErrNoConnections.Error()) {
			c.JSON(http.StatusBadGateway,
				management.GetExternalError(em))
		} else {
			c.JSON(http.StatusInternalServerError,
				management.GetExternalError(em))
		}
		ErrGetMappingByEventId.Incr()
		return
	}

	GetMappingByEventId.Incr()
	c.JSON(http.StatusOK, results)
}

func getMappingByEventIdQuery(accountId string, eventId string) (*model.MappingsData, error) {
	var mappingData = new(model.MappingsData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	if err := database.Select(Keyspace, MappingsTable).
		Value("id", &mappingData.Id).
		Value("eventid", &mappingData.EventId).
		Value("snippetid", &mappingData.SnippetId).
		Value("createdon", &mappingData.CreatedOn).
		Where("accountid", accountId).
		Where("eventid", eventId).
		Scan(session); err != nil {
		return nil, err
	}

	return mappingData, nil
}

func deleteMapping(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete mapping due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelMapping.Incr()
		return
	}

	mappingId := c.Params.ByName("mappingId")
	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %v", err)
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelMapping.Incr()
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, MappingsTable).
		Where("accountId", accountId).
		Where("id", mappingId).
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
		ErrDelMapping.Incr()
		return
	}

	if !success {
		ErrDelMapping.Incr()
		mlog.Error("Mapping %s not found in account %s", mappingId, accountId)
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Mapping not found. "))
		return
	}

	mlog.Info("Mapping %s deleted from account %s", mappingId, accountId)
	DelMapping.Incr()
	c.String(http.StatusOK, "")
}
