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

package nsObject

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	NsObject        = stats.New("nsObject")
	CreateBucket    = NsObject.NewCounter("CreateBucket")
	DeleteBucket    = NsObject.NewCounter("DeleteBucket")
	ListBuckets     = NsObject.NewCounter("ListBuckets")
	UploadFile      = NsObject.NewCounter("UploadFile")
	DownloadFile    = NsObject.NewCounter("DownloadFile")
	DeleteFile      = NsObject.NewCounter("DeleteFile")
	ListFiles       = NsObject.NewCounter("ListFiles")
	ErrCreateBucket = NsObject.NewCounter("ErrCreateBucket")
	ErrDeleteBucket = NsObject.NewCounter("ErrDeleteBucket")
	ErrListBuckets  = NsObject.NewCounter("ErrListBuckets")
	ErrUploadFile   = NsObject.NewCounter("ErrUploadFile")
	ErrDownloadFile = NsObject.NewCounter("ErrDownloadFile")
	ErrDeleteFile   = NsObject.NewCounter("ErrDeleteFile")
	ErrListFiles    = NsObject.NewCounter("ErrListFiles")
)
