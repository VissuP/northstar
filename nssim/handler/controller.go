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
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pmylund/go-cache"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	log "github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/service_master"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/utils"
)

import (
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/execution"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/kv"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/notebook"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/notebook/execution"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/nsql/native"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/nsql/spark"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/object"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/template"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/transformation"
	_ "github.com/verizonlabs/northstar/nssim/tests/northstar/transformation/execution"
)

var (
	MaxExecutions = 5
)

// Controller struct
type Controller struct {
	//serviceMaster *service_master.ServiceMaster
	summary       Summary
	serviceMaster *service_master.ServiceMaster
	asyncResults  *cache.Cache
	stats         *stats.Stats
	callbacks     *utils.ThreadSafeMap
}

// Returns a new controller
func NewController() (*Controller, error) {
	log.Debug("NewController")

	//Initialize test summary information based on configuration
	summary := Summary{
		StartTime: time.Now(),
	}

	testStats := stats.New(config.Configuration.ServiceName)

	for index, configTest := range config.Configuration.Tests {
		if !config.Configuration.IsTestEnabled(&configTest) {
			mlog.Info("%+v is not enabled", configTest.Name)
			continue
		}

		// TODO - Add concurrency here. For now assume just one.
		// if test.Concurrency > MaxExecutions {
		//	 MaxExecutions = test.Concurrency
		// }

		// Create test structure with up to max executions to trace.
		test := &Test{
			ID:            index,
			State:         tests.Ready,
			Name:          configTest.Name,
			Type:          configTest.Type,
			Group:         configTest.Group,
			Concurrency:   1,
			Verbose:       configTest.Verbose,
			MaxExecutions: MaxExecutions,
		}

		// Initialize the test metrics.
		test.InitializeMetrics(testStats)

		// Add to list of tests.
		summary.Tests = append(summary.Tests, test)
	}

	numWorkers := len(summary.Tests)
	cacheExpiration := 1 * time.Hour
	cachePurgeInterval := 30 * time.Minute

	// Controller
	controller := &Controller{
		serviceMaster: service_master.New(numWorkers, numWorkers),
		asyncResults:  cache.New(cacheExpiration, cachePurgeInterval),
		summary:       summary,
		stats:         testStats,
		callbacks:     utils.NewThreadSafeMap(),
	}

	return controller, nil
}

// RenderServiceError is a helper method used to render http response from given management error object.
func (controller *Controller) RenderServiceError(context *gin.Context, serviceError *management.Error) {
	// per docs, headers need to be set before calling context.JSON method
	for k, v := range serviceError.Header {
		for _, v1 := range v {
			context.Writer.Header().Add(k, v1)
		}
	}
	// now serialize rest of the response
	context.JSON(serviceError.HttpStatus, serviceError)
}

// Bind is a helper method used to bind body based on supported content types.
func (controller *Controller) Bind(context *gin.Context, resource interface{}) error {
	request := context.Request
	bind := controller.getBinding(request.Method, gin.MIMEJSON)

	if err := bind.Bind(request, resource); err != nil {
		return fmt.Errorf("Failed to bind request body with error: %v", err)
	}

	return nil
}

// getBinding is a helper method used to get content binding from content type.
func (controller *Controller) getBinding(method, contentType string) binding.Binding {
	mlog.Debug("getBinding: method:%s, contentType:%s", method, contentType)

	if method == "GET" {
		return binding.Form
	}

	// TODO - Add return by supported content types. For now, assuming JSON.
	return binding.JSON
}
