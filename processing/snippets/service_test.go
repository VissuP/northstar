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
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/mlog"
	invocationsClient "github.com/verizonlabs/northstar/data/invocations/client"
	snippetsClient "github.com/verizonlabs/northstar/data/snippets/client"
	"github.com/verizonlabs/northstar/pkg/rte/events/mocks"
)

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func TestInvokeSnippetDirectly(t *testing.T) {
	service := newFakeSnippetsService(t)
	engine := gin.Default()
	engine.POST(":accountId", service.startSnippet)

	reader := bytes.NewReader([]byte(`{
    "name":"printf",
    "runtime":"lua",
    "mainfn":"main",
    "url":"base64:///",
    "code":"LS0gRmlsdGVyIHN0YXRzLg0KDQpsb2NhbCBuc1FMID0gcmVxdWlyZSgibnNRTCIpDQpsb2NhbCBuc091dHB1dCA9IHJlcXVpcmUoIm5zT3V0cHV0IikNCg0KZnVuY3Rpb24gbWFpbigpDQogICAgbG9jYWwgcXVlcnkgPSAgIFtbDQogICAgICAgIFNFTEVDVCAgSlNPTl9GRVRDSChzdGF0cy5zdGF0cywgJ2V0aW1lJykgYXMgRXRpbWUsIEpTT05fRkVUQ0goc3RhdHMuc3RhdHMsICdpbnRlcnZhbCcpIGFzIEludGVyDQogICAgICAgIEZST00gICAgZGFrb3RhX2tzLnN0YXRzOw0KICAgIF1dDQogICAgbG9jYWwgc291cmNlID0gew0KICAgICAgICBQcm90b2NvbCA9ICJjYXNzYW5kcmEiLA0KICAgICAgICBIb3N0ID0gImNhc3NhbmRyYTEtbG9nLWRha290YS5tb24tbWFyYXRob24tc2VydmljZS5tZXNvcyIsDQogICAgICAgIFBvcnQgPSAiOTA0MiIsDQogICAgICAgIEJhY2tlbmQgPSAic3BhcmsiDQogICAgfQ0KICAgIGxvY2FsIG9wdGlvbnMgPSB7UmV0dXJuUm93c0xpbWl0ID0gMTAwfQ0KICAgIGxvY2FsIHJlc3VsdCA9IHByb2Nlc3NRdWVyeShxdWVyeSwgc291cmNlLCBvcHRpb25zKQ0KICAgIHJldHVybiBnZW5lcmF0ZVRhYmxlKHJlc3VsdCkNCmVuZA0KDQpmdW5jdGlvbiBwcm9jZXNzUXVlcnkocXVlcnksIHNvdXJjZSwgb3B0aW9ucykNCiAgICBsb2NhbCByZXNwLCBlcnIgPSBuc1FMLnF1ZXJ5KHF1ZXJ5LCBzb3VyY2UsIG9wdGlvbnMpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIHJlc3ANCmVuZA0KDQpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQ0KICAgIGxvY2FsIG91dCwgZXJyID0gbnNPdXRwdXQudGFibGUodGFibGUpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIG91dA0KZW5kDQo=",
    "timeout": 5000
  }`))
	req, _ := http.NewRequest("POST", "/abcd24b8-9a30-11e6-822b-acbc12345678", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestInvokeSnippetDirectlyMissingName(t *testing.T) {
	service := newFakeSnippetsService(t)

	engine := gin.Default()
	engine.POST(":accountId", service.startSnippet)

	reader := bytes.NewReader([]byte(`{
    "runtime":"lua",
    "mainfn":"main",
    "url":"base64:///",
    "code":"LS0gRmlsdGVyIHN0YXRzLg0KDQpsb2NhbCBuc1FMID0gcmVxdWlyZSgibnNRTCIpDQpsb2NhbCBuc091dHB1dCA9IHJlcXVpcmUoIm5zT3V0cHV0IikNCg0KZnVuY3Rpb24gbWFpbigpDQogICAgbG9jYWwgcXVlcnkgPSAgIFtbDQogICAgICAgIFNFTEVDVCAgSlNPTl9GRVRDSChzdGF0cy5zdGF0cywgJ2V0aW1lJykgYXMgRXRpbWUsIEpTT05fRkVUQ0goc3RhdHMuc3RhdHMsICdpbnRlcnZhbCcpIGFzIEludGVyDQogICAgICAgIEZST00gICAgZGFrb3RhX2tzLnN0YXRzOw0KICAgIF1dDQogICAgbG9jYWwgc291cmNlID0gew0KICAgICAgICBQcm90b2NvbCA9ICJjYXNzYW5kcmEiLA0KICAgICAgICBIb3N0ID0gImNhc3NhbmRyYTEtbG9nLWRha290YS5tb24tbWFyYXRob24tc2VydmljZS5tZXNvcyIsDQogICAgICAgIFBvcnQgPSAiOTA0MiIsDQogICAgICAgIEJhY2tlbmQgPSAic3BhcmsiDQogICAgfQ0KICAgIGxvY2FsIG9wdGlvbnMgPSB7UmV0dXJuUm93c0xpbWl0ID0gMTAwfQ0KICAgIGxvY2FsIHJlc3VsdCA9IHByb2Nlc3NRdWVyeShxdWVyeSwgc291cmNlLCBvcHRpb25zKQ0KICAgIHJldHVybiBnZW5lcmF0ZVRhYmxlKHJlc3VsdCkNCmVuZA0KDQpmdW5jdGlvbiBwcm9jZXNzUXVlcnkocXVlcnksIHNvdXJjZSwgb3B0aW9ucykNCiAgICBsb2NhbCByZXNwLCBlcnIgPSBuc1FMLnF1ZXJ5KHF1ZXJ5LCBzb3VyY2UsIG9wdGlvbnMpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIHJlc3ANCmVuZA0KDQpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQ0KICAgIGxvY2FsIG91dCwgZXJyID0gbnNPdXRwdXQudGFibGUodGFibGUpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIG91dA0KZW5kDQo=",
    "timeout": 5000
  }`))
	req, _ := http.NewRequest("POST", "/abcd24b8-9a30-11e6-822b-acbc12345678", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestInvokeSnippetDirectlyMissingMainFn(t *testing.T) {
	service := newFakeSnippetsService(t)

	engine := gin.Default()
	engine.POST(":accountId", service.startSnippet)

	reader := bytes.NewReader([]byte(`{
    "name":"printf",
    "runtime":"lua",
    "url":"base64:///",
    "code":"LS0gRmlsdGVyIHN0YXRzLg0KDQpsb2NhbCBuc1FMID0gcmVxdWlyZSgibnNRTCIpDQpsb2NhbCBuc091dHB1dCA9IHJlcXVpcmUoIm5zT3V0cHV0IikNCg0KZnVuY3Rpb24gbWFpbigpDQogICAgbG9jYWwgcXVlcnkgPSAgIFtbDQogICAgICAgIFNFTEVDVCAgSlNPTl9GRVRDSChzdGF0cy5zdGF0cywgJ2V0aW1lJykgYXMgRXRpbWUsIEpTT05fRkVUQ0goc3RhdHMuc3RhdHMsICdpbnRlcnZhbCcpIGFzIEludGVyDQogICAgICAgIEZST00gICAgZGFrb3RhX2tzLnN0YXRzOw0KICAgIF1dDQogICAgbG9jYWwgc291cmNlID0gew0KICAgICAgICBQcm90b2NvbCA9ICJjYXNzYW5kcmEiLA0KICAgICAgICBIb3N0ID0gImNhc3NhbmRyYTEtbG9nLWRha290YS5tb24tbWFyYXRob24tc2VydmljZS5tZXNvcyIsDQogICAgICAgIFBvcnQgPSAiOTA0MiIsDQogICAgICAgIEJhY2tlbmQgPSAic3BhcmsiDQogICAgfQ0KICAgIGxvY2FsIG9wdGlvbnMgPSB7UmV0dXJuUm93c0xpbWl0ID0gMTAwfQ0KICAgIGxvY2FsIHJlc3VsdCA9IHByb2Nlc3NRdWVyeShxdWVyeSwgc291cmNlLCBvcHRpb25zKQ0KICAgIHJldHVybiBnZW5lcmF0ZVRhYmxlKHJlc3VsdCkNCmVuZA0KDQpmdW5jdGlvbiBwcm9jZXNzUXVlcnkocXVlcnksIHNvdXJjZSwgb3B0aW9ucykNCiAgICBsb2NhbCByZXNwLCBlcnIgPSBuc1FMLnF1ZXJ5KHF1ZXJ5LCBzb3VyY2UsIG9wdGlvbnMpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIHJlc3ANCmVuZA0KDQpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQ0KICAgIGxvY2FsIG91dCwgZXJyID0gbnNPdXRwdXQudGFibGUodGFibGUpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIG91dA0KZW5kDQo=",
    "timeout": 5000
  }`))
	req, _ := http.NewRequest("POST", "/abcd24b8-9a30-11e6-822b-acbc12345678", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusInternalServerError := w.Code == http.StatusInternalServerError
		p, err := ioutil.ReadAll(w.Body)
		pageOk := (err == nil) &&
			strings.TrimRight(string(p), "\n") == "{\"error\":\"service_error\",\"error_description\":\"MainFN is empty\"}"
		return statusInternalServerError && pageOk
	})
}

func TestInvokeSnippetDirectlyMissingUrl(t *testing.T) {
	service := newFakeSnippetsService(t)

	engine := gin.Default()
	engine.POST(":accountId", service.startSnippet)

	reader := bytes.NewReader([]byte(`{
    "name":"printf",
    "runtime":"lua",
    "mainfn":"main",
    "code":"LS0gRmlsdGVyIHN0YXRzLg0KDQpsb2NhbCBuc1FMID0gcmVxdWlyZSgibnNRTCIpDQpsb2NhbCBuc091dHB1dCA9IHJlcXVpcmUoIm5zT3V0cHV0IikNCg0KZnVuY3Rpb24gbWFpbigpDQogICAgbG9jYWwgcXVlcnkgPSAgIFtbDQogICAgICAgIFNFTEVDVCAgSlNPTl9GRVRDSChzdGF0cy5zdGF0cywgJ2V0aW1lJykgYXMgRXRpbWUsIEpTT05fRkVUQ0goc3RhdHMuc3RhdHMsICdpbnRlcnZhbCcpIGFzIEludGVyDQogICAgICAgIEZST00gICAgZGFrb3RhX2tzLnN0YXRzOw0KICAgIF1dDQogICAgbG9jYWwgc291cmNlID0gew0KICAgICAgICBQcm90b2NvbCA9ICJjYXNzYW5kcmEiLA0KICAgICAgICBIb3N0ID0gImNhc3NhbmRyYTEtbG9nLWRha290YS5tb24tbWFyYXRob24tc2VydmljZS5tZXNvcyIsDQogICAgICAgIFBvcnQgPSAiOTA0MiIsDQogICAgICAgIEJhY2tlbmQgPSAic3BhcmsiDQogICAgfQ0KICAgIGxvY2FsIG9wdGlvbnMgPSB7UmV0dXJuUm93c0xpbWl0ID0gMTAwfQ0KICAgIGxvY2FsIHJlc3VsdCA9IHByb2Nlc3NRdWVyeShxdWVyeSwgc291cmNlLCBvcHRpb25zKQ0KICAgIHJldHVybiBnZW5lcmF0ZVRhYmxlKHJlc3VsdCkNCmVuZA0KDQpmdW5jdGlvbiBwcm9jZXNzUXVlcnkocXVlcnksIHNvdXJjZSwgb3B0aW9ucykNCiAgICBsb2NhbCByZXNwLCBlcnIgPSBuc1FMLnF1ZXJ5KHF1ZXJ5LCBzb3VyY2UsIG9wdGlvbnMpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIHJlc3ANCmVuZA0KDQpmdW5jdGlvbiBnZW5lcmF0ZVRhYmxlKHRhYmxlKQ0KICAgIGxvY2FsIG91dCwgZXJyID0gbnNPdXRwdXQudGFibGUodGFibGUpDQogICAgaWYoZXJyIH49IG5pbCkgdGhlbg0KICAgICAgICBlcnJvcihlcnIpDQogICAgZW5kDQogICAgcmV0dXJuIG91dA0KZW5kDQo=",
    "timeout": 5000
  }`))
	req, _ := http.NewRequest("POST", "/abcd24b8-9a30-11e6-822b-acbc12345678", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		statusInternalServerError := w.Code == http.StatusInternalServerError
		p, err := ioutil.ReadAll(w.Body)
		pageOK := (err == nil) &&
			strings.TrimRight(string(p), "\n") == "{\"error\":\"service_error\",\"error_description\":\"URL is empty\"}"
		return statusInternalServerError && pageOK
	})
}

func TestInvokeSnippetById(t *testing.T) {
	service := newFakeSnippetsService(t)

	engine := gin.Default()
	engine.POST(":accountId", service.startSnippet)

	reader := bytes.NewReader([]byte(`{
    	"snippetId": "4536cdce-b7e3-4bca-b7ac-5cf3e902dc52"
    }`))

	req, _ := http.NewRequest("POST", "/abcd24b8-9a30-11e6-822b-acbc12345678", reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestInvokeSnippetByIdDataClientFailure(t *testing.T) {
	service := &SnippetsService{
		SnippetManagerStore: mocks.MockSnippetManagerStore{},
		SnippetClient:       snippetsClient.MockSnippetsClientFail{},
		InvocationsClient:   invocationsClient.MockInvocationClientFail{},
	}

	engine := gin.Default()
	engine.POST(":accountId", service.startSnippet)

	reader := bytes.NewReader([]byte(`{
    	"snippetId":"4536cdce-b7e3-4bca-b7ac-5cf3e902dc52"
    }`))
	req, _ := http.NewRequest("POST", "/abcd24b8-9a30-11e6-822b-acbc12345678", reader)
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

// Helper function to return a snippet service instance which is completely
// disconnected from any real data service clients for mock testing
func newFakeSnippetsService(t *testing.T) *SnippetsService {
	mlog.Info("Creating new fake service")
	return &SnippetsService{
		SnippetManagerStore: mocks.MockSnippetManagerStore{},
		SnippetClient:       snippetsClient.MockSnippetsClient{},
		InvocationsClient:   invocationsClient.MockInvocationClient{},
	}
}
