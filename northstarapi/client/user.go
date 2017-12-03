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

package client

import (
	"encoding/json"
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

const (
	UsersPath = "users"
)

// SearchUsers creates a new runtime object.
func (client *Client) SearchUsers(accessToken string, user *model.User) ([]model.User, *management.Error) {
	mlog.Debug("SearchUsers")

	// Create runtime.
	path := client.getResourcePath(UsersPath) + "/actions/search"
	headers := client.getRequestHeaders(accessToken)

	response, mErr := management.PostJSONWithHeaders(client.baseUrl, path, user, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resources.
	users := make([]model.User, 0)

	if err := json.Unmarshal(response, &users); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return users, nil
}
