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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/portal/model"
)

// ListPortfolios returns portfolios associated with authenticated user.
func (controller *Controller) ListPortfolios(context *gin.Context) {
	mlog.Debug("ListPortfolio")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get request token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	portfolios, mErr := controller.portalProvider.ListPortfolios(token.AccessToken)
	if mErr != nil {
		mlog.Error("Couldn't get portfolio list. Error was: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, portfolios)
}

// ListFiles returns files associated with authenticated user for a particular portfolio.
func (controller *Controller) ListFiles(context *gin.Context) {
	mlog.Debug("ListFiles")

	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get request token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	//get the portfolio name
	portfolioName := strings.TrimSpace(context.Params.ByName("portfolioName"))
	if portfolioName == "" {
		mlog.Error("Failed to get files due to bad request - Invalid portfolio name")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("portfolio name"))
		return
	}

	files, mErr := controller.portalProvider.ListFiles(token.AccessToken, portfolioName, "", 0, "") // default values ?
	if mErr != nil {
		mlog.Error("Couldn't get file list. Error was: %s", mErr.Description)
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.JSON(http.StatusOK, files)
}

// GetFile returns the file under specified portfolio.
func (controller *Controller) GetFile(context *gin.Context) {
	mlog.Debug("GetFile")
	token, tokenExists := controller.getAuthToken(context)
	if !tokenExists {
		mlog.Error("Failed to get request token value")
		controller.RenderServiceError(context, model.ErrorInvalidCookie)
		return
	}

	//get the portfolio name
	portfolioName := strings.TrimSpace(context.Params.ByName("portfolioName"))
	if portfolioName == "" {
		mlog.Error("Failed to get files due to bad request - Invalid portfolio name")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("portfolio name"))
		return
	}

	//get the file name
	fileName := strings.TrimSpace(context.Params.ByName("fileName"))
	if fileName == "" {
		mlog.Error("Failed to get file due to bad request - Invalid file name")
		controller.RenderServiceError(context, model.GetErrorMissingResourceID("file name"))
		return
	}

	file, mErr := controller.portalProvider.GetFile(token.AccessToken, portfolioName, fileName)
	if mErr != nil {
		controller.RenderServiceError(context, management.ErrorInternal)
		return
	}

	context.Data(http.StatusOK, file.ContentType, file.Payload)
}
