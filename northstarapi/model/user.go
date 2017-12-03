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

const (
	// Define the supported permissions  (e.g., owner, viewer, contributor, etc).
	ReadPermissions             string = "Read"
	ReadWritePermissions        string = "ReadWrite"
	ReadWriteExecutePermissions string = "ReadWriteExecute"
	OwnerPermissions            string = "Owner"
)

// Defines the type used to represent a notebook user.
type User struct {
	Id          string `json:"id,omitempty"`
	AccountId   string `json:"accountId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	ImageId     string `json:"imageId,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}
