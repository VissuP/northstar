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

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// Define the supported permissions  (e.g., owner, viewer, contributor, etc).
	NoPermissions               string = "None"
	ReadPermissions             string = "Read"
	ReadWritePermissions        string = "ReadWrite"
	ReadExecute                 string = "ReadExecute"
	ReadWriteExecutePermissions string = "ReadWriteExecute"
	OwnerPermissions            string = "Owner"
)

// Defines the type that represents notebook access information.
type Access struct {
	Id         string    `json:"id,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
	AccountId  string    `json:"accountId,omitempty"`
	UserId     string    `json:"userId,omitempty"`
	Permission string    `json:"permission,omitempty"`
	NotebookId string    `json:"notebookId,omitempty"`
}

// Defines internal type for marshaling/unmarshaling.
type typeAccess Access

// Helper method used to unmarshal notebook while validating required fields.
func (access *Access) UnmarshalJSON(data []byte) error {
	var value typeAccess

	// Unmarshal to the internal type
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*access = Access(value)

	// Validate the account id.
	if access.AccountId == "" {
		return fmt.Errorf("The account id is missing.")
	}

	// Validate the user id.
	if access.UserId == "" {
		return fmt.Errorf("The user id is missing.")
	}

	// Validate the permissions.
	switch access.Permission {
	case ReadPermissions:
	case ReadWritePermissions:
	case ReadWriteExecutePermissions:
	case ReadExecute:
	case OwnerPermissions:
	default:
		return fmt.Errorf("The permissions is missing or invalid.")
	}

	return nil
}

// Defines the type used to query notebooks.
type Query struct {
	AccountId  string `json:"accountId,omitempty"`
	UserId     string `json:"userId,omitempty"`
	NotebookId string `json:"notebookId,omitempty"`
}
