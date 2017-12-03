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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/kafkamngr/cluster"
)

var fakeCluster cluster.MockCluster

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// set up fake cluster for testing
	kafkaCluster = fakeCluster

	// Run the other tests
	os.Exit(m.Run())
}

func TestCreateTopic(t *testing.T) {
	engine := gin.Default()
	engine.POST("/", createTopic)

	reader := bytes.NewReader([]byte(`{
	  "name":"test_topic",
	  "partitions": 2,
	  "replication": 2
	}`))
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusCreated
	})
}

func TestCreateTopicFail(t *testing.T) {
	engine := gin.Default()
	engine.POST("/", createTopic)

	reader := bytes.NewReader([]byte(`{
	  "name":"",
	  "partitions": 2,
	  "replication": 2
	}`))
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestCreateTopicAlreadyExists(t *testing.T) {
	engine := gin.Default()
	engine.POST("/", createTopic)

	reader := bytes.NewReader([]byte(`{
	  "name":"existingTopic",
	  "partitions": 2,
	  "replication": 2
	}`))
	req, _ := http.NewRequest("POST", "/", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusConflict
	})
}

func TestGetTopicNames(t *testing.T) {
	engine := gin.Default()
	engine.GET("/", getTopicNames)

	req, _ := http.NewRequest("GET", "/", nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := (err == nil) &&
			strings.TrimRight(string(p), "\n") == "[\"topic1\",\"topic2\",\"topic3\"]"
		return statusOK == true && pageOK == true
	})
}

func TestUpdateTopic(t *testing.T) {
	engine := gin.Default()
	engine.POST("/:serviceName/:topicName", updateTopic)

	reader := bytes.NewReader([]byte(`{
	  "partitions": 1
	}`))
	req, _ := http.NewRequest("POST", "/kafkamngr/existingTopic", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestUpdateTopicMissingServiceName(t *testing.T) {
	engine := gin.Default()
	engine.POST("/:topicName", updateTopic)

	reader := bytes.NewReader([]byte(`{
	  "partitions": 1
	}`))
	req, _ := http.NewRequest("POST", "/existingTopic", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusBadRequest
	})
}

func TestUpdateTopicMissingTopicName(t *testing.T) {
	engine := gin.Default()
	engine.POST("/:serviceName", updateTopic)

	reader := bytes.NewReader([]byte(`{
	  "partitions": 1
	}`))
	req, _ := http.NewRequest("POST", "/serviceName", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusBadRequest
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
