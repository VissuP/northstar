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

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/auth/config"
	"github.com/verizonlabs/northstar/auth/handler"
)

type Service struct {
	controller *handler.Controller
	engine     *gin.Engine
}

func NewService() (*Service, error) {
	mlog.Debug("NewService")

	err := config.Load()
	if err != nil {
		mlog.Error("Error, failed to load service configuration with error: %s\n", err.Error())
		return nil, err
	}

	controller, err := handler.NewController()
	if err != nil {
		mlog.Error("Error, failed to create example service controller with error %s.\n", err.Error())
		return nil, err
	}

	service := &Service{
		controller: controller,
		engine:     management.Engine(),
	}

	engine := management.Engine()

	//When user logs in to the portal
	engine.POST("/oauth2/token", controller.Oauth2TokenHandler)
	//When user logs out of the portal
	engine.POST("/oauth2/revoke", controller.Oauth2RevokeHandler)
	//When user navigates to the transformations page
	engine.GET("/api/v2/models", controller.GetModelsHandler)
	//Request coming from portal to get information of the authenticated user.
	engine.GET("/api/v2/users/me", controller.GetUserHandler)

	// APIs coming from NS API
	engine.GET("/oauth2/token/info", controller.Oauth2TokenInfoHandler)
	//called by NorthstarAPI for sharing functionality
	engine.GET("/south/v2/users", controller.SouthUserHandler)

	return service, nil
}

func (service *Service) Start() (err error) {
	mlog.Debug("Start")
	return management.Listen(":8080")
}
