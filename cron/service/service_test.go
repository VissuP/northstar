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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/cron/model"
	"github.com/verizonlabs/northstar/cron/scheduler"
	"github.com/verizonlabs/northstar/cron/util"
	cronClient "github.com/verizonlabs/northstar/data/cron/client"
	"github.com/verizonlabs/northstar/data/snippets/client"
	snippetsClient "github.com/verizonlabs/northstar/processing/snippets/client"
)

const (
	AccountId = "a3a424b8-9a30-11e6-822b-acbc32d30e43"
	SnippetId = "c2cf4526-a90b-4d4f-bf5b-e28f7e69802f"
	JobId     = "50d74215-b89b-4c51-97b6-5cd0ef9f0ac4"
	JobName   = "Job name"
	Schedule  = "0 * * * * *"
)

var (
	errMissingName      = management.NewError(http.StatusInternalServerError, ValidationFailed, model.MissingName)
	errMissingSchedule  = management.NewError(http.StatusInternalServerError, ValidationFailed, model.MissingSchedule)
	errMissingSnippetId = management.NewError(http.StatusInternalServerError, ValidationFailed, model.MissingSnippetId)
	errMissingAccountId = management.GetBadRequestError(util.AccountIdMissing)
)

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func TestAddJob(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.POST("/:accountId", service.addJob)

	json, _ := json.Marshal(&model.Job{
		AccountId: AccountId,
		Name:      JobName,
		Disabled:  false,
		Schedule:  Schedule,
		SnippetId: SnippetId,
	})
	reader := bytes.NewReader(json)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestAddJobMissingName(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.POST("/:accountId", service.addJob)

	job, _ := json.Marshal(&model.Job{
		AccountId: AccountId,
		Name:      "",
		Disabled:  false,
		Schedule:  Schedule,
		SnippetId: SnippetId,
	})
	reader := bytes.NewReader(job)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusInternalServerError := w.Code == http.StatusInternalServerError
		var errMessage management.Error
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Error(err)
		}

		merr := json.Unmarshal(p, &errMessage)
		if merr != nil {
			t.Error(merr)
		}

		identifierOK := errMessage.Id == errMissingName.Id
		descriptionOK := errMessage.Description == errMissingName.Description
		return statusInternalServerError && identifierOK && descriptionOK
	})
}

func TestAddJobMissingSchedule(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.POST("/:accountId", service.addJob)

	job, _ := json.Marshal(&model.Job{
		AccountId: AccountId,
		Name:      JobName,
		Disabled:  false,
		SnippetId: SnippetId,
	})
	reader := bytes.NewReader(job)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusInternalServerError := w.Code == http.StatusInternalServerError
		var errMessage management.Error
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Error(err)
		}

		merr := json.Unmarshal(p, &errMessage)
		if merr != nil {
			t.Error(merr)
		}

		identifierOK := errMessage.Id == errMissingSchedule.Id
		descriptionOK := errMessage.Description == errMissingSchedule.Description
		return statusInternalServerError && identifierOK && descriptionOK
	})
}

func TestAddJobMissingSnippetId(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.POST("/:accountId", service.addJob)

	job, _ := json.Marshal(&model.Job{
		AccountId: AccountId,
		Name:      JobName,
		Disabled:  false,
		Schedule:  Schedule,
	})
	reader := bytes.NewReader(job)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusInternalServerError := w.Code == http.StatusInternalServerError
		var errMessage management.Error
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Error(err)
		}

		merr := json.Unmarshal(p, &errMessage)
		if merr != nil {
			t.Error(merr)
		}

		identifierOK := errMessage.Id == errMissingSnippetId.Id
		descriptionOK := errMessage.Description == errMissingSnippetId.Description
		return statusInternalServerError && identifierOK && descriptionOK
	})
}

func TestUpdateJob(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.PUT("/:accountId/:jobId", service.updateJob)

	json, _ := json.Marshal(&model.Job{
		AccountId: AccountId,
		Name:      JobName,
		Disabled:  false,
		Schedule:  Schedule,
		SnippetId: SnippetId,
	})
	reader := bytes.NewReader(json)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/%s/%s", AccountId, JobId), reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestUpdateJobMissingAccountId(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.PUT("/:jobId", service.updateJob)

	job, _ := json.Marshal(&model.Job{
		AccountId: AccountId,
		Name:      JobName,
		Disabled:  false,
		Schedule:  Schedule,
		SnippetId: SnippetId,
	})
	reader := bytes.NewReader(job)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/%s", JobId), reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusCodeOK := w.Code == http.StatusBadRequest
		var errMessage management.Error
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Error(err)
		}

		merr := json.Unmarshal(p, &errMessage)
		if merr != nil {
			t.Error(merr)
		}

		identifierOK := errMessage.Id == errMissingAccountId.Id
		descriptionOK := errMessage.Description == errMissingAccountId.Description
		return statusCodeOK && identifierOK && descriptionOK
	})
}

func TestDeleteJob(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.DELETE("/:accountId/:jobId", service.deleteJob)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s", AccountId, JobId), nil)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestDeleteJobMissingAccountId(t *testing.T) {
	service := newMockCronService()
	engine := gin.Default()
	engine.DELETE("/:jobId", service.deleteJob)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s", JobId), nil)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusCodeOK := w.Code == http.StatusBadRequest
		var errMessage management.Error
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Error(err)
		}

		merr := json.Unmarshal(p, &errMessage)
		if merr != nil {
			t.Error(merr)
		}

		identifierOK := errMessage.Id == errMissingAccountId.Id
		descriptionOK := errMessage.Description == errMissingAccountId.Description
		return statusCodeOK && identifierOK && descriptionOK
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

// Convenience function for creating a Cron service instance that doesn't connect
// to any other currently running services
func newMockCronService() *CronService {
	return &CronService{
		processingClient:   snippetsClient.SnippetsClientMock{},
		snippetsDataClient: client.MockSnippetsClient{},
		cronDataClient:     cronClient.CronClientMock{},
		scheduler:          scheduler.SchedulerMock{},
	}
}
