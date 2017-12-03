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

package storage

import (
	"net/http"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/object/model"
)

// Defines the storage provider interface
type StorageMockFailure struct{}

func (s StorageMockFailure) CreateBucket(bucketName string) *management.Error {
	mlog.Info("Returning error attempting to create a bucket")
	return management.NewError(http.StatusInternalServerError, "", "")
}

func (s StorageMockFailure) DeleteBucket(bucketName string) *management.Error {
	mlog.Info("Returning error attempting to delete bucket with name: %s", bucketName)
	return management.NewError(http.StatusInternalServerError, "", "")
}

func (s StorageMockFailure) ListBuckets(nameFilter string) ([]model.Bucket, *management.Error) {
	mlog.Info("Returning error attempting to list buckets with name filter: %s", nameFilter)
	return nil, management.NewError(http.StatusInternalServerError, "", "")
}

func (s StorageMockFailure) List(bucketName, fileName string) ([]model.Object, string, *management.Error) {
	mlog.Info("Returning error attempting to list buckets with bucketName: %s and fileName: %s", bucketName, fileName)
	return nil, "", management.NewError(http.StatusInternalServerError, "", "")
}

func (s StorageMockFailure) Delete(bucketName, fileName string) *management.Error {
	mlog.Info("Returning error attempting to delete bucket with bucketName: %s and fileName: %s", bucketName, fileName)
	return management.NewError(http.StatusInternalServerError, "", "")
}

func (s StorageMockFailure) Upload(bucketName string, data *model.UploadData) *management.Error {
	mlog.Info("Returning error attempting to upload bucket with bucketName: %s", bucketName)
	return management.NewError(http.StatusInternalServerError, "", "")
}

func (s StorageMockFailure) Download(bucketName, fileName string) (*model.DownloadData, *management.Error) {
	mlog.Info("Returning error attempting to download bucket with bucketName: %s and fileName: %s", bucketName, fileName)
	return nil, management.NewError(http.StatusInternalServerError, "", "")
}
