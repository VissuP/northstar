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

package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/object/mocks/storage"
	"github.com/verizonlabs/northstar/object/model"
	"github.com/verizonlabs/northstar/object/util"
)

const (
	AccountId  = "a3a424b8-9a30-11e6-822b-acbc32d30e43"
	BucketName = "Test bucket"
	FileName   = "testFile.txt"
)

var (
	ErrMissingBucketName  = management.GetBadRequestError(util.BucketNameMissing)
	ErrMissingFileName    = management.GetBadRequestError(util.FileNameMissing)
	ErrMissingAccountId   = management.GetBadRequestError(util.AccountIdMissing)
	ErrMissingPayload     = management.GetBadRequestError(util.PayloadMissing)
	ErrMissingContentType = management.GetBadRequestError(util.ContentTypeMissing)
)

func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Run the other tests
	os.Exit(m.Run())
}

func TestUploadFile(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST(":accountId/:bucketName", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{
		FileName:    "file.txt",
		Payload:     []byte("Random data"),
		ContentType: "test/plain",
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s/%s", AccountId, BucketName), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestUploadFileMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST(":bucketName", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{
		FileName:    "file.txt",
		Payload:     []byte("Random data"),
		ContentType: "test/plain",
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s", BucketName), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
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

func TestUploadFileMissingBucketName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST("/:accountId", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{
		FileName:    "file.txt",
		Payload:     []byte("Random data"),
		ContentType: "test/plain",
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s", AccountId), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
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

		identifierOK := errMessage.Id == ErrMissingBucketName.Id
		descriptionOK := errMessage.Description == ErrMissingBucketName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestUploadFileMissingFileName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST(":accountId/:bucketName", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{
		FileName:    "", // A required field which is missing
		Payload:     []byte("Random data"),
		ContentType: "test/plain",
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s/%s", AccountId, BucketName), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
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

		identifierOK := errMessage.Id == ErrMissingFileName.Id
		descriptionOK := errMessage.Description == ErrMissingFileName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestUploadFileMissingPayload(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST(":accountId/:bucketName", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{ // Payload is missing
		FileName:    "file.txt",
		ContentType: "test/plain",
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s/%s", AccountId, BucketName), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
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

		identifierOK := errMessage.Id == ErrMissingPayload.Id
		descriptionOK := errMessage.Description == ErrMissingPayload.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestUploadFileMissingContentType(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST(":accountId/:bucketName", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{
		FileName:    "file.txt",
		Payload:     []byte("Random data"),
		ContentType: "", // Content type is the required field that is empty
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s/%s", AccountId, BucketName), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
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

		identifierOK := errMessage.Id == ErrMissingContentType.Id
		descriptionOK := errMessage.Description == ErrMissingContentType.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestUploadFileStorageProviderFailing(t *testing.T) {
	controller := NewController(storage.StorageMockFailure{})
	engine := gin.Default()
	engine.POST(":accountId/:bucketName", controller.UploadFile)

	fileToUpload, _ := json.Marshal(&model.UploadData{
		FileName:    "file.txt",
		Payload:     []byte("Random data"),
		ContentType: "test/plain",
	})

	reader := bytes.NewReader(fileToUpload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/%s/%s", AccountId, BucketName), reader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("FileName", "file.txt")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestDownloadFile(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET("/:accountId/:bucketName/*fileName", controller.DownloadFile)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/%s", AccountId, BucketName, FileName), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestDownloadFileMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET("/:bucketName/*fileName", controller.DownloadFile)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", BucketName, FileName), nil)
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

func TestDownloadFileMissingBucketName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET("/:accountId/*fileName", controller.DownloadFile)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", AccountId, FileName), nil)
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

		identifierOK := errMessage.Id == ErrMissingBucketName.Id
		descriptionOK := errMessage.Description == ErrMissingBucketName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestDownloadFileMissingFileName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET("/:accountId/:bucketName", controller.DownloadFile)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", AccountId, BucketName), nil)
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

		identifierOK := errMessage.Id == ErrMissingFileName.Id
		descriptionOK := errMessage.Description == ErrMissingFileName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestDownloadFileFailingClient(t *testing.T) {
	controller := NewController(storage.StorageMockFailure{})
	engine := gin.Default()
	engine.GET("/:accountId/:bucketName/*fileName", controller.DownloadFile)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s/%s", AccountId, BucketName, FileName), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestDeleteFile(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE("/:accountId/:bucketName/*fileName", controller.DeleteFile)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s/%s", AccountId, BucketName, FileName), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestDeleteFileMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE("/:bucketName/*fileName", controller.DeleteFile)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s", BucketName, FileName), nil)
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

func TestDeleteFileMissingBucketName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE("/:accountId/*fileName", controller.DeleteFile)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s", AccountId, FileName), nil)
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

		identifierOK := errMessage.Id == ErrMissingBucketName.Id
		descriptionOK := errMessage.Description == ErrMissingBucketName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestDeleteFileMissingFileName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE("/:accountId/:bucketName", controller.DeleteFile)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s", AccountId, BucketName), nil)
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

		identifierOK := errMessage.Id == ErrMissingFileName.Id
		descriptionOK := errMessage.Description == ErrMissingFileName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestDeleteFileStorageProviderFailing(t *testing.T) {
	controller := NewController(storage.StorageMockFailure{})
	engine := gin.Default()
	engine.DELETE("/:accountId/:bucketName/*fileName", controller.DeleteFile)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s/%s", AccountId, BucketName, FileName), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestCreateBucket(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST("/:accountId", controller.CreateBucket)

	json, _ := json.Marshal(&model.Bucket{
		Name:         "Test bucket",
		CreationDate: time.Now(),
	})

	reader := bytes.NewReader(json)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestCreateBucketMissingName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST("/:accountId", controller.CreateBucket)

	bucket, _ := json.Marshal(&model.Bucket{
		Name:         "",
		CreationDate: time.Now(),
	})

	reader := bytes.NewReader(bucket)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
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

		identifierOK := errMessage.Id == ErrMissingBucketName.Id
		descriptionOK := errMessage.Description == ErrMissingBucketName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestCreateBucketMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.POST("/", controller.CreateBucket)

	bucket, _ := json.Marshal(&model.Bucket{
		Name:         "Test bucket",
		CreationDate: time.Now(),
	})

	reader := bytes.NewReader(bucket)
	req, _ := http.NewRequest("POST", "/", reader)
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

func TestCreateBucketStorageProviderFailing(t *testing.T) {
	controller := NewController(storage.StorageMockFailure{})
	engine := gin.Default()
	engine.POST("/:accountId", controller.CreateBucket)

	json, _ := json.Marshal(&model.Bucket{
		Name:         "Test bucket",
		CreationDate: time.Now(),
	})

	reader := bytes.NewReader(json)
	req, _ := http.NewRequest("POST", "/"+AccountId, reader)
	req.Header.Add("Content-Type", "application/json")
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestListBuckets(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET("/:accountId", controller.ListBuckets)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", AccountId), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestListBucketsMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET("/", controller.ListBuckets)

	req, _ := http.NewRequest("GET", "/", nil)
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

func TestListBucketsStorageProviderFailure(t *testing.T) {
	controller := NewController(storage.StorageMockFailure{})
	engine := gin.Default()
	engine.GET("/:accountId", controller.ListBuckets)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", AccountId), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusInternalServerError
	})
}

func TestDeleteBucket(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE(":accountId/:bucketName", controller.DeleteBucket)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s/%s", AccountId, BucketName), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestDeleteBucketMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE(":bucketName", controller.DeleteBucket)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s", BucketName), nil)
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

func TestDeleteBucketMissingBucketName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.DELETE(":accountId", controller.DeleteBucket)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s", AccountId), nil)
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

		identifierOK := errMessage.Id == ErrMissingBucketName.Id
		descriptionOK := errMessage.Description == ErrMissingBucketName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestListFiles(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET(":accountId/:bucketName", controller.ListFiles)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", AccountId, BucketName), nil)
	testHTTPResponse(t, engine, req, func(w *httptest.ResponseRecorder) bool {
		return w.Code == http.StatusOK
	})
}

func TestListFilesMissingAccountId(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET(":bucketName", controller.ListFiles)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", BucketName), nil)
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

func TestListFilesMissingBucketName(t *testing.T) {
	controller := NewController(storage.StorageMock{})
	engine := gin.Default()
	engine.GET(":accountId", controller.ListFiles)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", AccountId), nil)
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

		identifierOK := errMessage.Id == ErrMissingBucketName.Id
		descriptionOK := errMessage.Description == ErrMissingBucketName.Description
		return isStatusBadRequest && identifierOK && descriptionOK
	})
}

func TestListFilesStorageProviderFailure(t *testing.T) {
	controller := NewController(storage.StorageMockFailure{})
	engine := gin.Default()
	engine.GET(":accountId/:bucketName", controller.ListFiles)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/%s", AccountId, BucketName), nil)
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
