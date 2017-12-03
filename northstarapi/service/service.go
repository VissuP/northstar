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
	"path"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/middleware"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace"
	"github.com/verizonlabs/northstar/northstarapi/config"
	"github.com/verizonlabs/northstar/northstarapi/handler"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Defines the type that represents the service.
type Service struct {
	controller *handler.Controller
	engine     *gin.Engine
}

// Returns a new service instance.
func NewService() (*Service, error) {
	mlog.Info("NewService")

	// Load Configuration and create a controller
	if err := config.Load(); err != nil {
		return nil, fmt.Errorf("Failed to load config with error %s.", err.Error())
	}

	controller, err := handler.NewController()

	if err != nil {
		return nil, fmt.Errorf("Failed to create controller with error %s.", err.Error())
	}

	// Set CORS support access control.
	accessControl := middleware.AccessControl{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "Etag", "If-Match", "If-No-Match"},
	}

	// Create thingspace auth client.
	authClient := thingspace.NewThingSpaceAuthClient(config.Configuration.ThingSpaceAuthHostPort)

	// Create Application REST services
	engine := management.Engine()

	// Add common middlewares.
	engine.Use(middleware.Cors(accessControl))

	//Register internal endpoints
	internal := engine.Group(path.Join(model.ContextInternal, model.Version))
	{
		internal.POST("/callbacks/execution", controller.ExecutionCallback)
	}

	// Register service APIs
	v1 := engine.Group(path.Join(model.Context, model.Version))
	v1.Use(middleware.Authorization(&config.Configuration.Scopes, authClient))
	{
		// Register User endpoints.
		v1.POST("/users/actions/search", controller.SearchUsers)

		// Register Transformation endpoints.
		v1.POST("/transformations", controller.CreateTransformation)
		v1.GET("/transformations", controller.ListTransformations)
		v1.PUT("/transformations", controller.UpdateTransformation)
		v1.GET("/transformations/:transformationId", controller.GetTransformation)
		v1.GET("/transformations/:transformationId/results", controller.GetTransformationResults)
		v1.DELETE("/transformations/:transformationId", controller.DeleteTransformation)

		// Register Transformation Schedule endpoints.
		v1.GET("/transformations/:transformationId/schedule", controller.GetSchedule)
		v1.POST("/transformations/:transformationId/schedule", controller.CreateSchedule)
		v1.DELETE("/transformations/:transformationId/schedule", controller.DeleteSchedule)

		// Register Notebook endpoints.
		v1.POST("/notebooks", controller.CreateNotebook)
		v1.GET("/notebooks", controller.ListNotebooks)
		v1.GET("/notebooks/:notebookId", controller.GetNotebook)
		v1.PUT("/notebooks", controller.UpdateNotebook)
		v1.POST("/notebooks/:notebookId/cells/:cellId/actions/execute", controller.ExecuteNotebookCell)
		v1.DELETE("/notebooks/:notebookId", controller.DeleteNotebook)

		// Register Notebook Users endpoints.
		v1.GET("/notebooks/:notebookId/users", controller.GetNotebookUsers)
		v1.PUT("/notebooks/:notebookId/users", controller.UpdateNotebookUsers)

		//Register long running stream endpoints.
		v1.GET("/streams", controller.ListStreams)
		v1.GET("/streams/:id", controller.GetStream)
		v1.DELETE("/streams/:id", controller.RemoveStream)

		// Register Executions endpoints. Note that these are for one off executions.

		v1.POST("/executions/cell", controller.ExecuteCell)
		v1.POST("/executions/transformation", controller.ExecuteTransformation)

		//Common Execution Endpoints
		v1.GET("/executions", controller.ListExecutions)
		v1.POST("/executions/generic", controller.TriggerExecution)
		v1.GET("/executions/generic/:id", controller.GetExecution)
		v1.POST("/executions/generic/:id/actions/stop", controller.StopExecution)

		// Register Templates endpoints.
		v1.POST("/templates", controller.CreateTemplate)
		v1.GET("/templates", controller.ListTemplates)
		v1.GET("/templates/:templateId", controller.GetTemplate)
		v1.DELETE("/templates/:templateId", controller.DeleteTemplate)

		//Register object endpoints
		v1.GET("/objects", controller.ListBuckets)
		v1.GET("/objects/:bucket/list/*path", controller.ListObjects)
		v1.GET("/objects/:bucket/get/*path", controller.GetObject)
	}

	// Create and return the service.
	service := &Service{
		controller: controller,
		engine:     engine,
	}

	return service, nil
}

// Start starts listening for incoming requests on default service port.
func (service *Service) Start() (err error) {
	mlog.Info("Start")
	return management.Listen(":8080")
}
