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

package events

import (
	"net/http"

	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/events/model"
	"github.com/verizonlabs/northstar/data/util"
)

var (
	eventsColumn = "id, name, createdon"
	sess         *gocql.Session
	lock         sync.Mutex
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

type EventsService struct{}

func (s *EventsService) AddRoutes() {
	grp := management.Engine().Group(util.DataBasePath)
	g := grp.Group("events")
	g.POST(":accountId", addEvent)
	g.GET(":accountId", getEvents)
	g.DELETE(":accountId/:eventId", deleteEvent)
}

func addEvent(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	var event = new(model.EventData)
	if err := c.Bind(event); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertEvent.Incr()
		return
	}

	err := event.ValidateOnAdd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrInsertEvent.Incr()
		return
	}

	if event.Id == "" {
		event.Id = uuid.NewV4().String()
	}

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertEvent.Incr()
		return
	}

	if _, err := database.Insert(Keyspace, EventsTable).
		Param("id", event.Id).
		Param("accountId", accountId).
		Param("name", event.Name).
		Param("createdon", time.Now().In(time.UTC)).
		Exec(session); err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrInsertEvent.Incr()
		return
	}

	InsertEvent.Incr()
	c.String(http.StatusCreated, event.Id)
}

func getEvents(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	results, err := getEventsQuery(accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		ErrGetEvents.Incr()
		return
	}

	GetEvents.Incr()
	c.JSON(http.StatusOK, results)
}

func getEventsQuery(accountId string) ([]model.EventData, error) {
	mlog.Info("Retrieving events for account %s", accountId)

	results := make([]model.EventData, 0, 10)
	entry := new(model.EventData)

	session, err := getSession()
	if err != nil {
		return nil, err
	}

	iter := session.
		Query(`SELECT `+eventsColumn+` FROM `+EventsTable+` WHERE accountid=?`, accountId).Iter()
	for iter.Scan(&entry.Id,
		&entry.Name,
		&entry.CreatedOn) {
		results = append(results, *entry)
		entry = new(model.EventData)
	}

	if err := iter.Close(); err != nil {
		mlog.Error("Error: ", err)
		return nil, err
	}

	return results, nil
}

func deleteEvent(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to delete event due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		ErrDelEvent.Incr()
		return
	}
	eventId := c.Params.ByName("eventId")

	session, err := getSession()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		c.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelEvent.Incr()
		return
	}

	success := true
	if success, err = database.Delete(Keyspace, EventsTable).
		Where("accountId", accountId).
		Where("id", eventId).
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
		ErrDelEvent.Incr()
		return
	}

	if !success {
		ErrDelEvent.Incr()
		mlog.Error("Event %s not found in account %s", eventId, accountId)
		c.JSON(http.StatusInternalServerError, management.GetInternalError("Event not found."))
		return
	}

	mlog.Info("Event %s deleted from account %s", eventId, accountId)
	DelEvent.Incr()
	c.String(http.StatusOK, "")
}
