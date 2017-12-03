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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/object/model"
)

// Defines the storage provider interface
type StorageMock struct{}

func (s StorageMock) CreateBucket(bucketName string) *management.Error {
	mlog.Info("Creating a bucket")
	return nil
}

func (s StorageMock) DeleteBucket(bucketName string) *management.Error {
	mlog.Info("Deleting bucket with name: %s", bucketName)
	return nil
}

func (s StorageMock) ListBuckets(nameFilter string) ([]model.Bucket, *management.Error) {
	mlog.Info("Listing buckets with name filter: %s", nameFilter)
	return []model.Bucket{}, nil
}

func (s StorageMock) List(bucketName, fileName string) ([]model.Object, string, *management.Error) {
	mlog.Info("Listing buckets with bucketName: %s and fileName: %s", bucketName, fileName)
	return []model.Object{}, "", nil
}

func (s StorageMock) Delete(bucketName, fileName string) *management.Error {
	mlog.Info("Deleting bucket with bucketName: %s and fileName: %s", bucketName, fileName)
	return nil
}

func (s StorageMock) Upload(bucketName string, data *model.UploadData) *management.Error {
	mlog.Info("Uploading bucket with bucketName: %s", bucketName)
	return nil
}

func (s StorageMock) Download(bucketName, fileName string) (*model.DownloadData, *management.Error) {
	mlog.Info("Downloading bucket with bucketName: %s and fileName: %s", bucketName, fileName)
	return &model.DownloadData{}, nil
}
