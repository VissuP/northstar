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
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/stats"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/utils"
	"strings"
	"sync"
	"time"
)

type Summary struct {
	StartTime   time.Time `json:"startTime"`
	Environment string    `json:"environment"`
	Tests       []*Test   `json:"tests"`
	TotalRuns   int       `json:"totalRuns"`
	Mode        string    `json:"mode"`
}

type Test struct {
	ID               int                        `json:"id,omitempty"`
	Name             string                     `json:"name,omitempty"`
	Type             string                     `json:"type,omitempty"`
	Group            string                     `json:"group,omitempty"`
	State            string                     `json:"state"`
	ConcurrencyIndex int                        `json:"concurrencyIndex"`
	Concurrency      int                        `json:"concurrency"`
	TotalExecutions  int                        `json:"totalExecutions"`
	TotalErrors      int                        `json:"totalErrors"`
	LastStatus       string                     `json:"lastStatus,omitempty"`
	LastLatency      string                     `json:"lastLatency,omitempty"`
	LastErrorMsg     string                     `json:"lastErrorMessage,omitempty"`
	ExecutionResults []ExecutionResult          `json:"executionResults, omitempty"`
	FailureStats     map[string]*utils.StepStat `json:"failureStats"`

	// For internal use only. These are to track test execution,
	// results, etc.
	Verbose           bool              `json:"-"`
	MaxExecutions     int               `json:"-"`
	LastExecution     int               `json:"-"`
	TestLogger        *utils.Logger     `json:"-"`
	FailureStatsLock  sync.RWMutex      `json:"-"`
	FailureStatsGroup *stats.Stats      `json:"-"`
	Status            chan tests.Status `json:"-"`

	// Test Metrics
	ExecutionCountMetric *stats.Counter            `json:"-"`
	ErrorCountMetric     *stats.Counter            `json:"-"`
	StatusMetric         *stats.Set                `json:"-"`
	LastErrorMetric      *stats.Set                `json:"-"`
	FailureMetrics       map[string]*stats.Counter `json:"-"`
}

// ExecutionResult defines the type used to store execution results.
type ExecutionResult struct {
	Status           string           `json:"status"`
	ConcurrencyIndex int              `json:"concurrencyIndex"`
	LogMessages      []utils.LogEntry `json:"logMessages"`
	Latency          string           `json:"latency"`
	Finished         time.Time        `json:"finished"`
}

// Initializes test metrics.
func (test *Test) InitializeMetrics(statsGroup *stats.Stats) error {
	mlog.Info("InitializeMetrics: test.Name=%s", test.Name)

	if test == nil || test.Name == "" {
		return fmt.Errorf("Test is invalid or does not have a valid name.")
	}

	test.ExecutionCountMetric = statsGroup.NewCounter(fmt.Sprintf("%sCount", test.Name))
	test.ErrorCountMetric = statsGroup.NewCounter(fmt.Sprintf("%sErrorCount", test.Name))
	test.StatusMetric = statsGroup.NewSet(fmt.Sprintf("%sStatus", test.Name))
	test.LastErrorMetric = statsGroup.NewSet(fmt.Sprintf("%sLastErrorMsg", test.Name))
	test.FailureMetrics = make(map[string]*stats.Counter)
	test.FailureStatsGroup = stats.New(test.Name)
	test.FailureStats = make(map[string]*utils.StepStat)
	return nil
}

// Updates test information from status at end of test.
func (test *Test) Finish(status *tests.Status) {
	mlog.Info("Finish: %+v", status)
	test.State = status.State
	test.LastStatus = status.Status
	test.TotalExecutions = test.TotalExecutions + 1
	test.LastLatency = fmt.Sprintf("%f", status.Latency.Seconds())
	// we count several errors in the test as onedev
	if status.NumberOfErrors > 0 {
		test.TotalErrors = test.TotalErrors + 1
	}
	test.LastErrorMsg = status.LastError

	// If we exceed the max number of results. Reset the array.
	if len(test.ExecutionResults) > MaxExecutions {
		test.ExecutionResults = make([]ExecutionResult, 0)
	}

	// If the test did not return any logs we use a default message.
	logMessages := []utils.LogEntry{
		{
			Time:    time.Now(),
			Message: "No messages provided in test results.",
		},
	}
	if test.TestLogger != nil {
		logMessages = test.TestLogger.Logs
	}

	test.updatePassFailStats(test.TestLogger.Stats)

	if status.Status != "" {
		executionResults := ExecutionResult{
			Status:           status.Status,
			ConcurrencyIndex: test.ConcurrencyIndex,
			Latency:          fmt.Sprintf("%f", status.Latency.Seconds()),
			LogMessages:      logMessages,
			Finished:         time.Now(),
		}

		test.ExecutionResults = append(test.ExecutionResults, executionResults)
	}

	// Update test metrics.
	test.ExecutionCountMetric.Incr()
	test.StatusMetric.Set(status.Status)
	test.LastErrorMetric.Set(status.LastError)

	// Note that we count several errors in the test as one.
	if status.NumberOfErrors > 0 {
		test.ErrorCountMetric.Incr()
	}
}

//UpdateStepStats updates step pass/fail stats for test steps
func (test *Test) updatePassFailStats(stats map[string]*utils.StepStat) {
	for statName, statValue := range stats {
		if statValue.Successes > 0 {
			successName := statName + " Successes"
			test.FailureStatsLock.RLock()
			statSuccess, ok := test.FailureMetrics[successName]
			test.FailureStatsLock.RUnlock()

			if !ok {

				statSuccess = test.FailureStatsGroup.NewCounter(successName)
				statSuccess.IncrBy(statValue.Successes)

				test.FailureStatsLock.Lock()
				test.FailureMetrics[successName] = statSuccess
				test.FailureStatsLock.Unlock()
			}
		}

		if statValue.Failures > 0 {
			failureName := statName + " Failures"

			test.FailureStatsLock.RLock()
			statFailure, ok := test.FailureMetrics[failureName]
			test.FailureStatsLock.RUnlock()

			if !ok {
				statFailure = test.FailureStatsGroup.NewCounter(failureName)
				statFailure.IncrBy(statValue.Successes)

				test.FailureStatsLock.Lock()
				test.FailureMetrics[failureName] = statFailure
				test.FailureStatsLock.Unlock()
			}
		}

		test.FailureStatsLock.Lock()
		stat, ok := test.FailureStats[statName]
		if !ok {
			stat = &utils.StepStat{}
			test.FailureStats[statName] = stat
		}
		stat.Successes += statValue.Successes
		stat.Failures += statValue.Failures
		test.FailureStatsLock.Unlock()
	}
}

// Returns the test summary.
func (controller *Controller) GetSummary(context *gin.Context) {
	mlog.Debug("GetSummary")

	context.JSON(http.StatusOK, controller.summary)
}

func (controller *Controller) TestResults(context *gin.Context) {
	test := strings.TrimSpace(context.Params.ByName("test"))
	mlog.Debug("TestResults Test:%s", test)
	context.String(http.StatusOK, "results")
}

func (controler *Controller) TestResultsById(context *gin.Context) {
	test := strings.TrimSpace(context.Params.ByName("test"))
	testID := strings.TrimSpace(context.Params.ByName("id"))
	mlog.Debug("TestResults %s", test)
	mlog.Debug("TestResultsById Test:%s ID:%s", test, testID)
	context.String(http.StatusOK, "results by id")
}
