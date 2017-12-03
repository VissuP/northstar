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
	"github.com/verizonlabs/northstar/portal/middleware"
	"github.com/verizonlabs/northstar/portal/model"
	"github.com/verizonlabs/northstar/portal/utils"
)

const (
	// Defines the constant used as the HTTP Cookie Name.
	httpCookieName = "Ns.Http.Cookie"
)

// GetUser returns the information of the authenticated user.
func (controller *Controller) GetUser(context *gin.Context) {
	mlog.Debug("GetUser")

	// Get the token for current request.
	token, found := controller.getAuthToken(context)

	if !found {
		mlog.Error("Failed to get access token from context key %s.", middleware.AccessTokenKeyName)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	// Get information of authenticated user.
	thingSpaceUser, mErr := controller.userClient.GetUser(token.AccessToken)

	if mErr != nil {
		mlog.Error("Failed to get user information with error: %v", mErr)
		controller.RenderServiceError(context, mErr)
		return
	}

	// Create object for authenticated user.
	user := model.User{
		Id:          thingSpaceUser.Id,
		DisplayName: thingSpaceUser.DisplayName,
		Name: model.Name{
			First:  thingSpaceUser.FirstName,
			Middle: thingSpaceUser.MiddleName,
			Last:   thingSpaceUser.LastName,
		},
		Email:   thingSpaceUser.Email,
		ImageId: thingSpaceUser.ImageId,
	}

	context.JSON(http.StatusOK, user)
}

// Login creates a new session token from specified user.
func (controller *Controller) Login(context *gin.Context) {
	mlog.Debug("Login")

	var user model.User

	// Get resource from request object.
	if err := controller.Bind(context, &user); err != nil {
		mlog.Error("Failed to create session with error: %v.", err)
		controller.RenderServiceError(context, model.ErrorParseRequestBody)
		return
	}

	// Validate email is not empty.
	if user.Email == "" {
		mlog.Error("Failed to create session due to missing email")
		controller.RenderServiceError(context, model.ErrorMissingEmail)
		return
	}

	// Validate password is not empty.
	if user.Password == "" {
		mlog.Error("Failed to create session due to missing password")
		controller.RenderServiceError(context, model.ErrorMissingPassword)
		return
	}

	// Get the client token. Note that this method implemets caching.
	token, mErr := controller.authClient.GetUserToken(controller.clientId, controller.clientSecret, user.Email, user.Password, controller.userScopes)

	if mErr != nil {
		mlog.Error("Failed to get user token with error: %v", mErr)
		controller.RenderServiceError(context, mErr)
		return
	}

	// Set the HTTP Cookie.
	if err := utils.SetCookie(context, token); err != nil {
		mlog.Error("Failed to store token in cookie with error: %v", err)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Logout creates a new session token from specified user.
func (controller *Controller) Logout(context *gin.Context) {
	mlog.Debug("Logout")

	// Get the token for current request.
	token, err := utils.GetToken(context)

	if err != nil {
		mlog.Error("Failed to get request token value with error: %v", err)
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	// Call ThingSpace to revoke the access token.
	mErr := controller.authClient.RevokeAccessToken(controller.clientId, controller.clientSecret, token.AccessToken)

	if mErr != nil {
		mlog.Error("Failed to revoke token with error: %v", mErr)
		controller.RenderServiceError(context, mErr)
		return
	}

	// Delete the HTTP cookie from the request.
	utils.DeleteCookie(context)

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// Verify creates a new session token from specified user.
func (controller *Controller) Verify(context *gin.Context) {
	mlog.Debug("Verify")

	// Check if the request contains a token.
	if _, err := utils.GetToken(context); err != nil {
		mlog.Error("Failed to get request token with error: %v", err)
		controller.RenderServiceError(context, model.ErrorMissingCookie)
		return
	}

	// Note that validation only checks the request contains a cookie
	// but does not validate the actual cookie. The cookie will be
	// validated when a REST API call is made.

	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}

// QueryUsers returns a list of the users matching the query
func (controller *Controller) QueryUsers(context *gin.Context) {
	mlog.Debug("QueryUsers")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get access token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	//Get user object from request.
	userFilter := &model.User{}
	if err := controller.Bind(context, userFilter); err != nil {
		mlog.Error("Failed to get user filter with error: %+v", err)
		context.JSON(http.StatusOK, []string{})
		return
	}

	// Get our filtered list of users
	users, serviceErr := controller.portalProvider.QueryUsers(token.AccessToken, userFilter)
	if serviceErr != nil {
		mlog.Error("Failed to filter users with error: %s", serviceErr.Description)
		context.JSON(http.StatusOK, []model.User{})
		return
	}

	context.JSON(http.StatusOK, users)
}
