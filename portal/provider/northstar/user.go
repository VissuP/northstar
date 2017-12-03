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

package northstar

import (
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/portal/model"
)

// QueryUsers returns a list of users matching the provided query
func (provider *NorthStarPortalProvider) QueryUsers(token string, user *model.User) ([]model.User, *management.Error) {
	mlog.Debug("QueryUsers")

	externalUsers, serviceErr := provider.northstarApiClient.SearchUsers(token, provider.ToExternalUser(user))
	if serviceErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("failed to get users: %s", serviceErr.Description))
	}

	users := []model.User{}
	for _, externalUser := range externalUsers {
		user := provider.FromExternalUser(&externalUser)
		users = append(users, *user)
	}

	return users, nil
}

// FromExternalUser maps from a NorthstarAPI user to a portal user
func (provider *NorthStarPortalProvider) FromExternalUser(externalUser *northstarApiModel.User) *model.User {
	return &model.User{
		AccountID:   externalUser.AccountId,
		Id:          externalUser.Id,
		DisplayName: externalUser.DisplayName,
		Email:       externalUser.Email,
		Permissions: externalUser.Permissions,
		ImageId:     externalUser.ImageId,
	}
}

// ToExternalUser maps from a portal user to a NorthstarAPI user
func (provider *NorthStarPortalProvider) ToExternalUser(user *model.User) *northstarApiModel.User {
	return &northstarApiModel.User{
		AccountId:   user.AccountID,
		Id:          user.Id,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Permissions: user.Permissions,
	}
}
