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

package thingspace

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
)

// User defines the type used to represent ThingSpace User
// Resource.
type User struct {
	// Identity
	Id          string    `json:"id,omitempty"`
	Kind        string    `json:"kind,omitempty"`
	Version     string    `json:"version,omitempty"`
	VersionId   string    `json:"versionid,omitempty"`
	CreatedOn   time.Time `json:"createdon,omitempty"`
	LastUpdated time.Time `json:"lastupdated,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`

	// Foreign
	ForeignId string   `json:"foreignid,omitempty"`
	TagIds    []string `json:"tagids,omitempty"`

	// User
	CredentialsId   string    `json:"credentialsid,omitempty"`
	CredentialsType string    `json:"credentialstype,omitempty"`
	State           string    `json:"state,omitempty"`
	DisplayName     string    `json:"displayname,omitempty"`
	AckTerms        bool      `json:"ackterms,omitempty"`
	AckTermsOn      time.Time `json:"acktermson,omitempty"`
	ImageId         string    `json:"imageid,omitempty"`
	Title           string    `json:"title,omitempty"`
	FirstName       string    `json:"firstname,omitempty"`
	MiddleName      string    `json:"middlename,omitempty"`
	LastName        string    `json:"lastname,omitempty"`
	Email           string    `json:"email,omitempty"`
	Mdn             string    `json:"mdn,omitempty"`
}

// GetUser returns information of the authenticated user.
func (userClient *UserClient) GetUser(accessToken string) (*User, *management.Error) {
	mlog.Debug("GetUser")

	// Generate request headers.
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}

	// Get response body.
	body, mErr := management.GetWithHeaders(userClient.hostAndPort, "/api/v2/users/me", headers)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Failed to get information of authenticated user with error: %v", mErr))
	}

	// Unmarshal user object.
	var user User

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to unmarshal user with error: %v", err))
	}

	return &user, nil
}
