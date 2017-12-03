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
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/auth/model"
)

// Controller struct
type Controller struct {
}

// Returns a new controller
func NewController() (*Controller, error) {
	mlog.Debug("NewController")

	controller := &Controller{}

	return controller, nil
}

//Function to send back the access token when the user logs in
func (controller *Controller) Oauth2TokenHandler(context *gin.Context) {
	mlog.Debug("Inside the Oauth2TokenHandler handler")
	context.Header("Access-Control-Allow-Credentials", "true")
	context.Header("Access-Control-Allow-Origin", " ")
	context.Header("Cache-Control", "no-store")
	context.Header("Set-Cookie", "token=bearer MTE5YWVkMmMtOWFjMi00ZDA5LWFlOWMtNTZjYTU0NmJmNGI0; Path=/; Expires=Fri, 20 Oct 2018 17:16:37 GMT; Max-Age=86400")

	userTokenBody := model.UserToken{
		AccessToken:  "MTE5YWVkMmMtOWFjMi00ZDA5LWFlOWMtNTZjYTU0NmJmNGI0",
		ExpiresIn:    2592000,
		RefreshToken: "NDNjMzVkZDYtZTI3Ni00ZjU0LWE0N2ItZGZkNzQ2ODk2N2Jl",
		Scope:        "ts.execution ts.stream ts.user ts.user.ro ts.transformation ts.transformation.ro ts.notebook ts.notebook.ro ts.model.ro ts.nsobject.ro",
		TokenType:    "bearer",
	}

	context.JSON(http.StatusOK, userTokenBody)

}

//Function to send back a response when the user logs out
func (controller *Controller) Oauth2RevokeHandler(context *gin.Context) {
	mlog.Debug("Inside the controller Oauth2RevokeHandler handler")

	context.Status(http.StatusOK)

}

//Function to validate the access token from the NS API
func (controller *Controller) Oauth2TokenInfoHandler(context *gin.Context) {
	mlog.Debug("Inside the Oauth2TokenInfoHandler")

	userTokenInfoBody := model.UserTokenInfo{
		GrantType:       "password",
		ExpiresIn:       2592000,
		CreatedAt:       "2017-10-20T21:32:49.736100969Z",
		UserName:        "ns-test-user@verizon.com",
		Scope:           []string{"ts.execution", "ts.stream", "ts.model.ro", "ts.nsobject.ro", "ts.user", "ts.user.ro", "ts.transformation", "ts.transformation.ro", "ts.notebook", "ts.notebook.ro"},
		ClientId:        "NorthStarPortalDummyAuthV1",
		CredentialsType: "ts.credential",
	}

	context.JSON(http.StatusOK, userTokenInfoBody)
}

// Function to respond to the /api/v2/users/me coming in from the portal
func (controller *Controller) GetUserHandler(context *gin.Context) {
	mlog.Debug("Inside the GetUserHandler")

	getUserBody := model.GetUserInfo{
		Id:              "50g13cb7-9efe-68cd-f696-2809185f4af9",
		Kind:            "ts.user",
		Version:         "1.0",
		VersionId:       "a8924031-ad1a-11e7-a004-02420a270d10",
		CreatedOn:       "2017-10-09T17:52:46.464Z",
		LastUpdated:     "2017-10-09T17:52:46.7Z",
		Number:          "3b60a7ea-1064-4808-bcfd-b11ec84b0c1f",
		RequestedOn:     "2017-10-09T17:52:46.463Z",
		Expiration:      2592000,
		CredentialsId:   "ns-test-user@verizon.com",
		CredentialsType: "ts.credential",
		State:           "Active",
		DisplayName:     "NS Test User",
		AckTermsOn:      "0001-01-01T00:00:00Z",
		FirstName:       "NS Test",
		LastName:        "User",
		Email:           "ns-test-user@verizon.com",
	}
	context.JSON(http.StatusOK, getUserBody)

}

//Function to return back an empty array corresponding to the /api/v2/models request when
// user navigates to the Transformations page on the portal
func (controller *Controller) GetModelsHandler(context *gin.Context) {
	mlog.Debug("Inside the GetModelsHandler")

	// Send back an empty array
	getModelBody := [0]int{}

	context.JSON(http.StatusOK, getModelBody)
}

// Function to respond to the /south/v2/users coming in from the NS API
func (controller *Controller) SouthUserHandler(context *gin.Context) {
	mlog.Debug("Inside the SouthUserHandler")

	// Send back an array with a single hard-coded user entry
	getSouthBody := model.SouthUserInfo{
		Id:              "50g13cb7-9efe-68cd-f696-2809185f4af9",
		Kind:            "ts.user",
		Version:         "1.0",
		VersionId:       "a8924031-ad1a-11e7-a004-02420a270d10",
		CreatedOn:       "2017-10-09T17:52:46.464Z",
		LastUpdated:     "2017-10-09T17:52:46.7Z",
		ForeignId:       "e121d569-c8db-600f-f443-ebb16fd54850",
		Number:          "3b60a7ea-1064-4808-bcfd-b11ec84b0c1f",
		RequestedOn:     "2017-10-09T17:52:46.463Z",
		Expiration:      2592000,
		CredentialsId:   "ns-test-user@verizon.com",
		CredentialsType: "ts.credential",
		State:           "Active",
		DisplayName:     "NS Test User",
		AckTermsOn:      "0001-01-01T00:00:00Z",
		FirstName:       "NS Test",
		LastName:        "User",
		Email:           "ns-test-user@verizon.com",
	}

	// NS API expects atleast a single user as a response to this request
	// create an array with single dummy element
	southUsersBody := make([]model.SouthUserInfo, 1)
	southUsersBody[0] = getSouthBody
	context.JSON(http.StatusOK, southUsersBody)

}
