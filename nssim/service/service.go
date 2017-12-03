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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	log "github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/handler"
)

type Service struct {
	controller *handler.Controller
	engine     *gin.Engine
}

func NewService() (*Service, error) {
	log.Debug("NewService")

	err := config.Load()
	if err != nil {
		log.Error("Error, failed to load service configuration with error: %s\n", err.Error())
		return nil, err
	}

	controller, err := handler.NewController()
	if err != nil {
		log.Error("Error, failed to create example service controller with error %s.\n", err.Error())
		return nil, err
	}

	service := &Service{
		controller: controller,
		engine:     management.Engine(),
	}
	management.SetHealth(controller.GetHealth())
	engine := management.Engine()

	sim := engine.Group("sim")
	v1 := sim.Group("v1")
	v1.GET("/info", controller.GetInfo)

	tests := v1.Group("tests")
	tests.GET("/", controller.GetSummary)
	tests.GET(":test/results", controller.TestResults)
	tests.GET(":test/results/:id", controller.TestResultsById)
	tests.POST(":test/execute", controller.ExecuteTest)
	//simapi.POST("config")
	//simapi.GET("config")
	sim.StaticFS("/web", gin.Dir("./web/dist", true))

	v1.POST("/callback/:callbackID", controller.ReceiveCallback)

	return service, nil
}

func (service *Service) Start() (err error) {
	log.Debug("Start")

	if config.Configuration.Mode == config.AUTORUN_MODE {
		log.Info("Starting in auto mode.")
		// If the test is started in Auto mode, execute
		// rest service on the background and block on
		// test execution.
		go management.Listen(fmt.Sprintf("%s:%s", config.Configuration.Host, config.Configuration.Port))
		service.controller.RunTests()
	} else {
		log.Info("Starting in manual mode.")
		// If the service is started in Manual mode,
		// block on rest service execution.
		management.Listen(fmt.Sprintf("%s:%s", config.Configuration.Host, config.Configuration.Port))
	}

	return nil
}
