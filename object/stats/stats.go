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

package stats

import (
	"github.com/verizonlabs/northstar/pkg/stats"
)

var (
	s = stats.New("object")

	// Create bucket
	CreateBucketReqCount                  = s.NewCounter("CreateBucketReqCount")
	ErrCreateBucketCount                  = s.NewCounter("ErrCreateBucketCount")
	ErrCreateBucketMissingAccountIdCount  = s.NewCounter("ErrCreateBucketMissingAccountIdCount")
	ErrCreateBucketMissingBucketNameCount = s.NewCounter("ErrCreateBucketMissingBucketNameCount")

	// List buckets
	ListBucketsReqCount                 = s.NewCounter("ListBucketsReqCount")
	ErrListBucketsCount                 = s.NewCounter("ErrCreateBucketCount")
	ErrListBucketsMissingAccountIdCount = s.NewCounter("ErrListBucketsMissingAccountIdCount")
	ErrListFilesMissingBucketNameCount  = s.NewCounter("ErrListFilesMissingBucketNameCount")

	// Delete buckets
	DeleteBucketReqCount                  = s.NewCounter("DeleteBucketReqCount")
	ErrDeleteBucketsCount                 = s.NewCounter("ErrDeleteBucketCount")
	ErrDeleteBucketMissingBucketNameCount = s.NewCounter("ErrDeleteBucketMissingBucketNameCount")

	// Upload file
	UploadFileReqCount                  = s.NewCounter("UploadfileReqCount")
	ErrUploadFileMissingAccountIdCount  = s.NewCounter("ErrUploadFileMissingAccountIdCount")
	ErrUploadValidateCount              = s.NewCounter("ErrUploadValidateCount")
	ErrUploadFileMissingBucketNameCount = s.NewCounter("ErrUploadFileMissingBucketNameCount")
	ErrUploadWriteDataCount             = s.NewCounter("ErrUploadWriteDataCount")
	ErrUploadWriteFailCount             = s.NewCounter("ErrUploadWriteFailCount")

	// Download file
	DownloadFileReqCount                  = s.NewCounter("DownloadfileReqCount")
	ErrDownloadMissingAccountIdCount      = s.NewCounter("ErrDownloadMissingAccountIdCount")
	ErrDownloadFileMissingBucketNameCount = s.NewCounter("ErrDownloadFileMissingBucketNameCount")
	ErrDownloadMissingFileNameCount       = s.NewCounter("ErrDownloadMissingFileNameCount")
	ErrDownloadReadFailCount              = s.NewCounter("ErrDownloadReadFailCount")

	// List files
	ListFilesReqCount                 = s.NewCounter("ListFilesReqCount")
	ErrListFilesMissingAccountIdCount = s.NewCounter("ErrListFilesMissingAccountIdCount")
	ErrListFilesCount                 = s.NewCounter("ErrListFilesCount")

	// Delete files
	DeleteFileReqCount                   = s.NewCounter("DeleteFileReqCount")
	ErrDeleteFileMissingAccountIdCount   = s.NewCounter("ErrDeleteFileMissingAccountIdCount")
	ErrDeleteFilesMissingBucketNameCount = s.NewCounter("ErrDeleteFilesMissingBucketNameCount")
	ErrDeleteFileMissingFileNameCount    = s.NewCounter("ErrDeleteFileMissingFileNameCount")
	ErrDeleteFileFailCount               = s.NewCounter("ErrDeleteFileFailCount")
)
