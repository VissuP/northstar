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
	"regexp"
	"strings"
	"sync"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/northstarapi/config"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

const (
	MaxUsers          int = 20000
	MaxSearchRoutines int = 20
)

// Basic regular expression for validating strings.
const (
	Email string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

var (
	rxEmail = regexp.MustCompile(Email)
)

// Defines the type used to support operations on ThingSpace account.
type ThingSpaceAccountProvider struct {
	client thingspace.Client
}

// Returns a new ThingSpace Account Provider.
func NewThingSpaceAccountProvider() (*ThingSpaceAccountProvider, error) {
	mlog.Info("NewThingSpaceAccountProvider")

	// Create ThingSpace client.
	client, mErr := thingspace.NewThingSpaceClient(config.Configuration.ThingSpaceSouthHostPort,
		config.Configuration.ThingSpaceClientId,
		config.Configuration.ThingSpaceClientSecret)

	if mErr != nil {
		return nil, fmt.Errorf("Failed to create thingspace client with error: %+v", mErr)
	}

	// Create the provider.
	provider := &ThingSpaceAccountProvider{
		client: client,
	}

	return provider, nil
}

// Returns the account id associated with the specified loginname.
func (provider *ThingSpaceAccountProvider) GetAccountIdForLoginname(loginname string) (string, *management.Error) {
	mlog.Debug("GetAccountIdForLoginname")

	user, mErr := provider.GetUser(loginname)

	if mErr != nil {
		return "", management.ErrorExternal
	}

	return user.AccountId, nil
}

// Returns the user with the specified loginname.
func (provider *ThingSpaceAccountProvider) GetUser(loginname string) (*model.User, *management.Error) {
	mlog.Debug("GetUser")

	// Get the user associated with the loginname
	query := api.Query{
		Filter: map[string]interface{}{
			api.FilterCredentialsID: loginname,
		},
	}

	users, mErr := provider.client.QueryUsers(query)

	if mErr != nil {
		mlog.Error("Failed to get users with error: %+v", mErr)
		return nil, management.ErrorExternal
	}

	// Validate we got a valid user.
	if len(users) == 0 {
		mlog.Error("No user found for loginname. E.g., query user returned empty list.")
		return nil, model.ErrorLoginNameNotFound
	}

	// Validate the user account has been verified.
	// Validate the account state.
	if users[0].State != api.USER_STATE_ACTIVE {
		mlog.Error("The user has not been activated.")
		return nil, management.ErrorForbidden
	}

	user := &model.User{
		Id:          users[0].Id,
		AccountId:   users[0].ForeignId,
		Email:       users[0].Email,
		DisplayName: users[0].DisplayName,
	}

	return user, nil
}

// Returns the user for the specified id.
func (provider *ThingSpaceAccountProvider) GetUserById(id string) (*model.User, *management.Error) {
	mlog.Debug("GetUserById: id:%s", id)

	// TODO(s)
	// - South seem to support only query by credentials id. Need to add support for get user by id.
	//
	tsDataHostPort := config.Configuration.ThingspaceProtocol + "://" + config.Configuration.ThingSpaceDataHostPort
	path := "/ds/v2/ts.user/" + id

	// Query thingspace user.
	response, mErr := management.Get(tsDataHostPort, path)

	if mErr != nil {
		mlog.Error("Query users returned error: %+v", mErr)
		return nil, mErr
	}

	tsUser := api.User{}

	if err := json.Unmarshal(response, &tsUser); err != nil {
		mlog.Error("Unmarshal returned error: %+v", err)
		return nil, management.GetInternalError(err.Error())
	}

	// Otherwise, we have a match.
	user := &model.User{
		Id:          tsUser.Id,
		AccountId:   tsUser.ForeignId,
		DisplayName: tsUser.DisplayName,
		Email:       tsUser.Email,
		ImageId:     tsUser.ImageId,
	}

	return user, nil
}

// Returns the collection of users that matches the specified search criteria.
func (provider *ThingSpaceAccountProvider) SearchUsers(user *model.User) ([]model.User, *management.Error) {
	mlog.Debug("SearchUsers: user:%+v", user)

	// TODO(s)
	// - This code should go away. ThingSpace should provider API that enables to search
	//   for users based on Query filter, etc. This might require ability for users
	//   to decide which information is private, etc. For now doing this by brute force.

	// If an email address is provided get user by email.
	if isEmail(user.Email) {
		user, mErr := provider.GetUser(user.Email)

		if mErr != nil {
			return nil, mErr
		}

		return []model.User{*user}, nil
	}

	// Otherwise, we need to search collection of users. Note that this method is
	// very inefficient and it will not scale to million of users, etc. Final
	// implementation should support such use cases.
	query := api.Query{
		LimitNumber: MaxUsers,
	}

	// Get all ThingSpace users (up to max users).
	tsUsers, mErr := provider.client.QueryUsers(query)

	if mErr != nil {
		mlog.Error("Query users returned error: %+v", mErr)
		return nil, management.ErrorExternal
	}

	// Search, up to 10, collection of TS users.
	return search(10, MaxSearchRoutines, user, tsUsers), nil
}

// Helper method used to verify if a string contains a substring while ignoring case.
func contains(s, substr string) bool {
	if s == "" || substr == "" {
		return false
	}

	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Helper method used to verify if string is a valid email.
func isEmail(s string) bool {
	return rxEmail.MatchString(s)
}

// Helper method used to search collection in parallel.
func search(numMatches int, maxRoutines int, criteria *model.User, users []api.User) []model.User {
	// If no users, return.
	if len(users) == 0 {
		return []model.User{}
	}

	// Define variables used to synchronize go-routines.
	lock := sync.RWMutex{}
	sem := make(chan bool, maxRoutines)

	matches := make([]model.User, 0)
	numRoutines := int(0)
	chunkSize := len(users) / maxRoutines

	// For every chunk, create a go-routine to process.
	for i := 0; i < len(users); i += chunkSize {

		// Calculate the end index based on the start index.
		endIndex := i + chunkSize

		// If the last chunck contains less than expected
		// adjust the end index.
		if endIndex > len(users) {
			endIndex = len(users)
		}

		// Process chunck in go-routine.
		go func(startIndex, endIndex int) {
			for i := startIndex; i < endIndex; i++ {

				// Verify if criteria matches current user.
				if contains(users[i].DisplayName, criteria.DisplayName) || contains(users[i].Email, criteria.Email) {
					user := model.User{
						Id:          users[i].Id,
						AccountId:   users[i].ForeignId,
						DisplayName: users[i].DisplayName,
						Email:       users[i].Email,
						ImageId:     users[i].ImageId,
					}

					lenMatches := int(0)

					// Sync update common array of matches.
					lock.Lock()
					if len(matches) < numMatches {
						matches = append(matches, user)
						lenMatches = len(matches)
					}
					lock.Unlock()

					// If we found num matches entries, stop.
					if lenMatches == numMatches {
						break
					}
				}
			}

			sem <- true

		}(i, endIndex)

		numRoutines = numRoutines + 1
	}

	// Wait for all go routines to finish
	for i := 0; i < numRoutines; i++ {
		<-sem
	}

	return matches
}
