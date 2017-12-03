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

package s3

import (
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/object/model"
)

type ProviderInput struct {
	Host      string
	UserId    string
	Secret    string
	DebugFlag bool
}

// Defines the storage provider interface
type StorageProvider interface {
	// Creates a bucket
	CreateBucket(bucketName string) *management.Error

	// Deletes a bucket
	DeleteBucket(bucketName string) *management.Error

	// Returns a list of buckets
	ListBuckets(nameFilter string) ([]model.Bucket, *management.Error)

	// Returns metadata for all items under the specified path.
	List(bucketName, fileName string) ([]model.Object, string, *management.Error)

	// Deletes the data for the specified file name.
	Delete(bucketName, fileName string) *management.Error

	// upload specified file name.
	Upload(bucketName string, data *model.UploadData) *management.Error

	// download specified file name.
	Download(bucketName, fileName string) (*model.DownloadData, *management.Error)
}
