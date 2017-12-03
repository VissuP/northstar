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

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/northstarapi/utils"
)

// Creates a new runtime.
func (controller *Controller) SearchUsers(context *gin.Context) {
	mlog.Info("SearchUser")

	// Get the user resource used as search criteria.
	var user model.User

	// Validate request message
	if err := controller.Bind(context, &user); err != nil {
		mlog.Error("Failed to search user with error: %v.", err)
		utils.ErrSearchUsers.Incr()
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Search users.
	users, mErr := controller.accountProvider.SearchUsers(&user)

	if mErr != nil {
		mlog.Error("Failed to search users with error: %v", mErr)
		utils.ErrSearchUsers.Incr()
		controller.RenderServiceError(context, management.ErrorExternal)
		return
	}

	utils.SearchUsers.Incr()
	context.JSON(http.StatusOK, users)
}
