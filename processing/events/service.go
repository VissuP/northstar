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

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/mappings/client"
	"github.com/verizonlabs/northstar/processing/events/model"
	"github.com/verizonlabs/northstar/processing/snippets"
	snippetsModel "github.com/verizonlabs/northstar/processing/snippets/model"
	"github.com/verizonlabs/northstar/processing/util"
)

type EventsService struct {
	Snippets      *snippets.SnippetsService
	MappingClient client.Client
}

func NewEventsService(snippets *snippets.SnippetsService) *EventsService {
	mappingClient, err := client.NewMappingsClient()
	if err != nil {
		return nil
	}

	return &EventsService{Snippets: snippets, MappingClient: mappingClient}
}

func (s *EventsService) AddRoutes() {
	grp := management.Engine().Group(util.ProcessingBasePath)
	g := grp.Group("events")
	g.POST(":accountId/:eventId", s.invokeEvent)
}

func (s *EventsService) invokeEvent(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		mlog.Error("Failed to invoke event due to bad request. Account Id is missing.")
		c.JSON(http.StatusBadRequest, management.GetBadRequestError(util.AccountIdMissing))
		util.ErrInvokeEvent.Incr()
		return
	}

	eventId := c.Params.ByName("eventId")

	var options = new(model.Options)
	c.Bind(options)

	mlog.Info("Invoking event %s at account %s", eventId, accountId)
	mapping, mErr := s.MappingClient.GetMappingByEventId(accountId, eventId)
	if mErr != nil {
		c.JSON(http.StatusInternalServerError, mErr)
		util.ErrGetMappingByEventId.Incr()
		return
	}

	snippetOptions := snippetsModel.Options{Args: options.Args}
	result, err := s.Snippets.StartSnippetById(accountId, mapping.SnippetId, &snippetOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
		util.ErrInvokeEvent.Incr()
		return
	}

	util.InvokeEvent.Incr()
	c.String(http.StatusOK, result)
}
