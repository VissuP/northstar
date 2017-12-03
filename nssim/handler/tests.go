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
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/utils"
)

// Starts execution of all configured tests.
func (controller *Controller) RunTests() (err error) {
	mlog.Info("Start - Running %d tests", len(config.Configuration.Tests))

	stop := false

	for !stop {
		// Execute all the tests
		for _, test := range controller.summary.Tests {
			// If the test is not already running. Launch it. Note that
			// this should avoid running load tests twice.
			if test.State == tests.Ready {
				// Initialize summary for the current test execution
				test.Status = make(chan tests.Status, 1)
				test.TestLogger = utils.NewLogger(test.Verbose)

				// Create test worker
				worker := tests.NewTestWorker(test.Type, test.Name, 1, test.Status, test.TestLogger, controller.callbacks)
				test.TestLogger.LogInfo("Calling dispatcher for test: %s. Active routines: %d, Queued: %d", test.Name, controller.serviceMaster.ActiveRoutines(), controller.serviceMaster.QueuedWork())
				controller.serviceMaster.Dispatch("SimWorker", worker)
			}
		}

		// Wait for tests results
		for _, test := range controller.summary.Tests {

			if test.Type == config.LOAD_TEST {
				// If the test is a load test (e.g., can execute for long time)
				// just check the results without waiting.
				select {
				case status := <-test.Status:
					test.Finish(&status)
					test.State = tests.Ready
				default:
					mlog.Info("Test %s still running...", test.Name)
				}

			} else {
				var status tests.Status

				// If the test is NOT a load test, wait for completion of timeout.
				select {
				case status = <-test.Status:
					test.Finish(&status)
					test.State = tests.Ready
				case <-time.After(time.Duration(config.Configuration.MaxExecutionTime) * time.Second):
					status.Status = tests.TimeoutError
					status.Latency = time.Duration(config.Configuration.MaxExecutionTime) * time.Second
					status.LastError = fmt.Sprintf("Test timeout. Execution was longer than %d seconds.", config.Configuration.MaxExecutionTime)
					status.NumberOfErrors = 1
					test.Finish(&status)
					test.State = tests.Ready
				}
			}
		}

		controller.summary.TotalRuns = controller.summary.TotalRuns + 1

		// Sleep between iterations
		delay := time.Duration(config.Configuration.IterationDelayInSec) * time.Second
		mlog.Info("Sleeping %d seconds before next iteration...", config.Configuration.IterationDelayInSec)
		time.Sleep(delay)
	}

	return nil
}

// Executes the specified tests.
func (controller *Controller) ExecuteTest(context *gin.Context) {
	name := strings.TrimSpace(context.Params.ByName("test"))
	mlog.Debug("ExecuteTest: %s", name)

	context.JSON(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
	return
}
