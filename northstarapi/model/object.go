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

package model

import "time"

type Bucket struct {
	Name         string    `json:"name,omitempty"`
	CreationDate time.Time `json:"creationDate,omitempty"`
}

type Object struct {
	Key          string    `json:"name,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	Size         int64     `json:"size,omitempty"`
	Etag         string    `json:"etag,omitempty"`
}

type Data struct {
	Payload     []byte `json:"payload,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

//ContentTypeFolder is the content type for folders
//const ContentTypeFolder = "application/x-directory"
const ContentTypeFolder = "application/xml"
