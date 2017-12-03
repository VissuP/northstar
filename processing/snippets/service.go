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
	"errors"
	"net/http"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	invocationsClient "github.com/verizonlabs/northstar/data/invocations/client"
	snippetClient "github.com/verizonlabs/northstar/data/snippets/client"
	"github.com/verizonlabs/northstar/pkg/rte/events"
	"github.com/verizonlabs/northstar/processing/snippets/model"
	"github.com/verizonlabs/northstar/processing/util"
)

type SnippetsService struct {
	SnippetManagerStore events.ManagerStore
	SnippetClient       snippetClient.Client
	InvocationsClient   invocationsClient.Client
}

func NewSnippetsService() (*SnippetsService, error) {
	snippetClient, err := snippetClient.NewSnippetClient()
	if err != nil {
		return nil, err
	}

	invocationsClient, err := invocationsClient.NewInvocationClient()
	if err != nil {
		return nil, err
	}

	return &SnippetsService{SnippetManagerStore: events.NewSnippetMngrStore(PROCESSING_SERVICE_NAME),
		SnippetClient:     snippetClient,
		InvocationsClient: invocationsClient}, nil
}

func (s *SnippetsService) AddRoutes() {
	grp := management.Engine().Group(util.ProcessingBasePath)
	g := grp.Group("snippets")
	g.POST(":accountId", s.startSnippet)
	g.DELETE(":accountId/:invocationId", s.stopSnippet)
}

func (s *SnippetsService) startSnippet(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("accountId is empty"))
		util.ErrStartEvent.Incr()
		return
	}

	var snippet = new(model.Snippet)
	c.Bind(snippet)

	var invocationId string
	if snippet.SnippetId != "" {
		var err error
		invocationId, err = s.StartSnippetById(accountId, snippet.SnippetId, &snippet.Options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
			util.ErrStartEvent.Incr()
			return
		}
	} else {
		err := snippet.Validate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
			util.ErrStartEvent.Incr()
			return
		}

		mlog.Debug("Validation passed")
		invocationId, err = s.invokeSnippetDirect(accountId, snippet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, management.GetInternalError(err.Error()))
			util.ErrStartEvent.Incr()
			return
		}
	}

	util.StartSnippet.Incr()
	c.String(http.StatusOK, invocationId)
}

func (s *SnippetsService) StartSnippetById(accountId string,
	snippetId string,
	options *model.Options) (string, error) {
	mlog.Debug("Starting snippet invocation by snippet id: %v", snippetId)

	snippet, mErr := s.SnippetClient.GetSnippet(accountId, snippetId)
	if mErr != nil {
		return "", errors.New(mErr.Error())
	}

	manager, err := s.SnippetManagerStore.GetManager(snippet.Runtime)
	if err != nil {
		return "", err
	}

	event := &events.SnippetStartEvent{
		SnippetId: snippet.Id,
		MainFn:    snippet.MainFn,
		Runtime:   snippet.Runtime,
		Timeout:   snippet.Timeout,
		Memory:    snippet.Memory,
		Args:      options.Args,
		URL:       snippet.URL,
		Code:      snippet.Code,
		Callback:  snippet.Callback,
	}

	invocationId, err := manager.SnippetStart(accountId, event)
	if err != nil {
		return "", err
	}

	mlog.Info("Snippet %s invoked by id with invocation id %s", snippet.Id, invocationId)
	return invocationId, nil
}

func (s *SnippetsService) invokeSnippetDirect(accountId string,
	snippet *model.Snippet) (string, error) {
	mlog.Debug("Starting direct snippet invocation for accountId: %v", accountId)
	eventsProducer, err := s.SnippetManagerStore.GetManager(snippet.Runtime)
	if err != nil {
		return "", err
	}

	start := &events.SnippetStartEvent{
		SnippetId: uuid.NewV4().String(),
		MainFn:    snippet.MainFn,
		Runtime:   snippet.Runtime,
		Timeout:   snippet.Timeout,
		Memory:    snippet.Options.Memory,
		Args:      snippet.Options.Args,
		URL:       snippet.URL,
		Code:      snippet.Code,
		Callback:  snippet.Options.Callback,
	}

	invocationId, err := eventsProducer.SnippetStart(accountId, start)
	if err != nil {
		return "", err
	}

	mlog.Info("Snippet %s invoked directly with invocation id %s", start.SnippetId, invocationId)
	return invocationId, nil
}

func (s *SnippetsService) stopSnippet(c *gin.Context) {
	accountId := c.Params.ByName("accountId")
	if accountId == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("accountId is empty"))
		util.ErrStopSnippet.Incr()
		return
	}

	invocationId := c.Params.ByName("invocationId")
	if invocationId == "" {
		c.JSON(http.StatusBadRequest, management.GetBadRequestError("invocationId is empty"))
		util.ErrStopSnippet.Incr()
		return
	}

	invocation, mErr := s.InvocationsClient.GetInvocation(accountId, invocationId)
	if mErr != nil {
		c.JSON(http.StatusBadRequest, management.GetInternalError(mErr.Error()))
		util.ErrStopSnippet.Incr()
		return
	}

	if invocation.Status != events.SNIPPET_RUNNING_EVENT {
		errM := fmt.Sprintf("Invocation is not in running state: %v", invocation.Status)
		c.JSON(http.StatusBadRequest, management.GetInternalError(errM))
		util.ErrStopSnippet.Incr()
		return
	}

	manager, err := s.SnippetManagerStore.GetManager(invocation.Runtime)
	if err != nil {
		c.JSON(http.StatusBadRequest, management.GetInternalError(err.Error()))
		util.ErrStopSnippet.Incr()
		return
	}

	stop := &events.SnippetStopEvent{
		InvocationId: invocationId,
	}

	err = manager.SnippetStop(accountId, invocation.Partition, stop)
	if err != nil {
		c.JSON(http.StatusBadRequest, management.GetInternalError(err.Error()))
		util.ErrStopSnippet.Incr()
		return
	}

	util.StopSnippet.Incr()
	c.String(http.StatusOK, "")
}
