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
)

// User defines the type used to represent a portal user.
type User struct {
	AccountID   string `json:"accountId, omitempty"`
	Id          string `json:"id,omitempty"`
	DisplayName string `json:"displayName"`
	Name        Name   `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	ImageId     string `json:"imageId,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

// Name defines the type used to represent a name.
type Name struct {
	First  string `json:"first"`
	Middle string `json:"middle"`
	Last   string `json:"last"`
}

// typeUser defines an internal event type for marshaling/unmarshaling.
type typeUser User

// MarshalJSON is a helper method used to marshal user while translating fields.
func (user *User) MarshalJSON() ([]byte, error) {
	var value typeUser

	value = typeUser(*user)

	// Make sure we do not expose the password.
	value.Password = ""

	return json.Marshal(&value)
}
