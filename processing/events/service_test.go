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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	invocationsClient "github.com/verizonlabs/northstar/data/invocations/client"
	mappingsClient "github.com/verizonlabs/northstar/data/mappings/client"
	snippetsClient "github.com/verizonlabs/northstar/data/snippets/client"
	"github.com/verizonlabs/northstar/pkg/rte/events/mocks"
	"github.com/verizonlabs/northstar/processing/snippets"
	"github.com/verizonlabs/northstar/processing/util"
)

const (
	snippetsClientToFail = iota
	invocationClientToFail
	mappingClientToFail
)

var (
	ErrMissingAccountId = management.GetBadRequestError(util.AccountIdMissing)
)

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func TestInvokeEvent(t *testing.T) {
	service := newFakeEventsService()
	engine := gin.Default()
	engine.POST(":accountId/:eventId", service.invokeEvent)

	reader := bytes.NewReader([]byte(`{}`))
	req, _ := http.NewRequest("POST",
		"/a3a424b8-9a30-11e6-822b-acbc32d30e43/a90dc0d9-1e53-42bb-93ea-3dbc99231394", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestInvokeEventMissingAccountId(t *testing.T) {
	service := newFakeEventsService()
	engine := gin.Default()
	engine.POST(":eventId", service.invokeEvent)

	reader := bytes.NewReader([]byte(`{}`))
	req, _ := http.NewRequest("POST", "/a90dc0d9-1e53-42bb-93ea-3dbc99231394", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		isStatusBadRequest := w.Code == http.StatusBadRequest
		var errMessage management.Error
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Error(err)
		}

		merr := json.Unmarshal(p, &errMessage)
		if merr != nil {
			t.Error(merr)
		}

		identifierOK := errMessage.Id == ErrMissingAccountId.Id
		descriptionOK := errMessage.Description == ErrMissingAccountId.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestInvokeEventSnippetsClientFails(t *testing.T) {
	service := newFakeEventsServiceFails(snippetsClientToFail)
	engine := gin.Default()
	engine.POST(":accountId/:eventId", service.invokeEvent)

	reader := bytes.NewReader([]byte(`{}`))
	req, _ := http.NewRequest("POST",
		"/a3a424b8-9a30-11e6-822b-acbc32d30e43/a90dc0d9-1e53-42bb-93ea-3dbc99231394", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestInvokeEventMappingClientFails(t *testing.T) {
	service := newFakeEventsServiceFails(mappingClientToFail)
	engine := gin.Default()
	engine.POST(":accountId/:eventId", service.invokeEvent)

	reader := bytes.NewReader([]byte(`{}`))
	req, _ := http.NewRequest("POST",
		"/a3a424b8-9a30-11e6-822b-acbc32d30e43/a90dc0d9-1e53-42bb-93ea-3dbc99231394", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T,
	engine *gin.Engine,
	req *http.Request,
	f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	engine.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

// Helper function to return a Event service instance which is completely
// disconnected from any real data service clients for mock testing
func newFakeEventsService() *EventsService {
	mlog.Info("Creating new fake service")
	return &EventsService{
		Snippets: &snippets.SnippetsService{
			SnippetManagerStore: mocks.MockSnippetManagerStore{},
			SnippetClient:       snippetsClient.MockSnippetsClient{},
			InvocationsClient:   invocationsClient.MockInvocationClient{},
		},
		MappingClient: mappingsClient.MappingClientMock{},
	}
}

func newFakeEventsServiceFails(clientToFail int) *EventsService {
	mlog.Info("Creating new fake service which fails any request")

	switch clientToFail {
	case invocationClientToFail:
		return &EventsService{
			Snippets: &snippets.SnippetsService{
				SnippetManagerStore: mocks.MockSnippetManagerStore{},
				SnippetClient:       snippetsClient.MockSnippetsClient{},
				InvocationsClient:   invocationsClient.MockInvocationClientFail{},
			},
			MappingClient: mappingsClient.MappingClientMock{},
		}
	case mappingClientToFail:
		return &EventsService{
			Snippets: &snippets.SnippetsService{
				SnippetManagerStore: mocks.MockSnippetManagerStore{},
				SnippetClient:       snippetsClient.MockSnippetsClient{},
				InvocationsClient:   invocationsClient.MockInvocationClient{},
			},
			MappingClient: mappingsClient.MappingClientMockFail{},
		}
	case snippetsClientToFail:
		return &EventsService{
			Snippets: &snippets.SnippetsService{
				SnippetManagerStore: mocks.MockSnippetManagerStore{},
				SnippetClient:       snippetsClient.MockSnippetsClientFail{},
				InvocationsClient:   invocationsClient.MockInvocationClient{},
			},
			MappingClient: mappingsClient.MappingClientMock{},
		}
	}

	return nil
}
